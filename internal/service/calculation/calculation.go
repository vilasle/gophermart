package calculation

import (
	"context"
	"errors"
	"regexp"
	"sync"

	wrap "github.com/pkg/errors"

	"github.com/vilasle/gophermart/internal/logger"
	repository "github.com/vilasle/gophermart/internal/repository/calculation"
	"github.com/vilasle/gophermart/internal/service"
)

type CalculationService struct {
	repCalc  repository.CalculationRepository
	repRules repository.CalculationRules
	mxRules  *sync.Mutex
	rules    map[int16]rule
	manager  *EventManager
}

type CalculationServiceConfig struct {
	repository.CalculationRepository
	repository.CalculationRules
	*EventManager
}

func NewCalculationService(config CalculationServiceConfig) *CalculationService {
	s := &CalculationService{
		repCalc:  config.CalculationRepository,
		repRules: config.CalculationRules,
		manager:  config.EventManager,
		mxRules:  &sync.Mutex{},
		rules:    make(map[int16]rule),
	}

	s.manager.RegisterHandler(NewOrder, s.calculateOrder)
	s.manager.RegisterHandler(NewRule, s.readRule)

	s.readAllRules()
	//TODO if in repository there are rows in queue need to raise events
	return s
}

func (c CalculationService) Register(ctx context.Context, dto service.RegisterCalculationRequest) error {
	//save on db; line on table need for unexpected finishing service
	addingQueueDto, addingCalc := c.prepareAddingDto(dto)

	if err := c.repCalc.AddCalculationResult(ctx, addingCalc); err != nil {
		return err
	}

	if err := c.repCalc.AddCalculationToQueue(ctx, addingQueueDto...); err != nil {
		return err
	}

	//raise event for running worker
	c.manager.RaiseEvent(NewOrder, dto)
	return nil
}

func (c CalculationService) prepareAddingDto(dto service.RegisterCalculationRequest) ([]repository.AddingCalculation, repository.AddCalculationResult) {
	addingDto := make([]repository.AddingCalculation, 0, len(dto.Products))

	orderNumber := dto.OrderNumber

	for _, product := range dto.Products {
		addingDto = append(addingDto, repository.AddingCalculation{
			OrderNumber: orderNumber,
			ProductName: product.Name,
			Price:       product.Price,
		})
	}

	calcDto := repository.AddCalculationResult{
		OrderNumber: orderNumber,
		Status:      repository.Registered,
		Value:       0,
	}

	return addingDto, calcDto
}

func (c CalculationService) Calculation(ctx context.Context, dto service.CalculationFilterRequest) (service.CalculationInfo, error) {
	result, err := c.repCalc.Calculations(ctx, repository.CalculationFilter{
		OrderNumber: dto.OrderNumber,
	})

	if err != nil {
		return service.CalculationInfo{}, err
	}

	if len(result) == 0 {
		return service.CalculationInfo{}, service.ErrEntityDoesNotExists
	}

	return c.fillCalculatedInfo(result[0]), nil
}

func (c CalculationService) fillCalculatedInfo(dto repository.CalculationInfo) service.CalculationInfo {
	return service.CalculationInfo{
		OrderNumber: dto.OrderNumber,
		Status:      c.statusView(dto.Status),
		Accrual:     dto.Value,
	}
}

func (c CalculationService) statusView(status repository.CalculationStatus) string {
	switch status {
	case repository.Invalid:
		return "INVALID"
	case repository.Processing:
		return "PROCESSING"
	case repository.Processed:
		return "PROCESSED"
	default:
		return ""
	}
}

// TODO may I need to change calculation with help regexp and check matching on sql?
// event.Type = NewOrder; event.Data = service.RegisterCalculationRequest
func (c CalculationService) calculateOrder(ctx context.Context, event Event) {
	dto, ok := event.Data.(service.RegisterCalculationRequest)
	if !ok {
		logger.Warn("was raise event with wrong data", "event", event.Type, "data", event.Data)
		return
	}

	number := dto.OrderNumber
	products := dto.Products

	updateDto := repository.AddCalculationResult{
		OrderNumber: number,
		Value:       0,
		Status:      repository.Processing,
	}

	err := c.repCalc.UpdateCalculationResult(ctx, updateDto)
	if err != nil {
		logger.Error("updating calculation result", "error", err, "data", dto)
	}

	var bonus float64
	for _, product := range products {
		bonus += c.calculateProduct(product)
	}

	resultDto := c.fillCalculatedDto(number, bonus)

	if err := c.repCalc.UpdateCalculationResult(ctx, resultDto); err != nil {
		logger.Error("saving calculation result", "error", err, "data", dto)
		return
	}

	clearDto := repository.ClearingCalculationQueue{OrderNumber: number}
	if err := c.repCalc.ClearCalculationsQueue(ctx, clearDto); err != nil {
		logger.Error("clearing queue was failed", "error", err, "data", clearDto)
		return
	}

}

func (c CalculationService) calculateProduct(product service.ProductRow) float64 {
	c.mxRules.Lock()
	defer c.mxRules.Unlock()

	for _, rule := range c.rules {
		if r := rule.calculate(product.Name, product.Price); r > 0 {
			return r
		}
	}
	return 0
}

func (c CalculationService) fillCalculatedDto(orderNumber string, value float64) repository.AddCalculationResult {
	status := repository.Invalid
	if value > 0 {
		status = repository.Processed
	}

	return repository.AddCalculationResult{
		OrderNumber: orderNumber,
		Value:       value,
		Status:      status,
	}
}

func (c CalculationService) readAllRules() error {
	rs, err := c.repRules.Rules(context.Background(), repository.RuleFilter{})
	if err != nil {

		return err
	}
	return c.fillRules(rs)
}

// event.Type = NewRule; event.Data = id rule on repository
func (c CalculationService) readRule(ctx context.Context, event Event) {
	id, ok := event.Data.(int16)
	if !ok {
		logger.Warn("was raise event with wrong data", "event", event.Type, "data", event.Data)
		return
	}
	rs, err := c.repRules.Rules(ctx, repository.RuleFilter{ID: id})
	if err != nil {
		logger.Error("getting specific rule", "id", id, "error", err)
		return
	}

	if err := c.fillRules(rs); err != nil {
		logger.Error("preparing rules for using", "error", err)
		return
	}
}

func (c *CalculationService) fillRules(rs []repository.RuleInfo) error {
	c.mxRules.Lock()
	defer c.mxRules.Unlock()

	errs := make([]error, 0, len(rs))

	for _, r := range rs {
		exp, err := regexp.Compile(r.Match)
		if err != nil {
			errs = append(errs, wrap.Wrapf(err, "invalid regexp %s", r.Match))
			continue
		}

		calcType, correct := service.DefineCalculationType(r.CalculationType)
		if !correct {
			errs = append(errs, errors.New("invalid calculation type"))
			continue
		}

		c.rules[r.ID] = rule{
			exp:             exp,
			calculationType: calcType,
			value:           r.Point,
		}
	}
	return errors.Join(errs...)
}
