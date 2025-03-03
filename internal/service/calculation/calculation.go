package calculation

import (
	"context"
	"errors"
	"regexp"
	"sync"
	"time"

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

	return s
}

func (c CalculationService) Start(ctx context.Context) error {
	c.readAllRules(ctx)
	if err := c.runNotProcessedOrders(ctx); err != nil {
		return errors.Join(err, errors.New("failed to run not processed orders"))
	}
	return nil
}

func (c CalculationService) Register(ctx context.Context, dto service.RegisterCalculationRequest) error {
	//save on db; line on table need for unexpected finishing service
	addingQueueDto, addingCalc := prepareAddingDto(dto)

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

	return prepareCalculatedInfo(result[0]), nil
}

// EVENTS
//
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

// event.Type = NewOrder; event.Data = service.RegisterCalculationRequest
func (c CalculationService) calculateOrder(ctx context.Context, event Event) {
	var (
		number   string
		products []service.ProductRow
	)

	if dto, ok := event.Data.(service.RegisterCalculationRequest); ok {
		number = dto.OrderNumber
		products = dto.Products
	} else {
		logger.Warn("was raise event with wrong data", "event", event.Type, "data", event.Data)
		return
	}

	updateDto := repository.AddCalculationResult{
		OrderNumber: number,
		Value:       0,
		Status:      repository.Processing,
	}

	err := c.repCalc.UpdateCalculationResult(ctx, updateDto)
	if err != nil {
		logger.Error("updating calculation result", "error", err, "data", updateDto)
	}

	bonus := c.calculateProductsBonus(products)

	logger.Debug("calculating order", "order", number, "products", products, "bonus", bonus)

	resultDto := prepareCalculatedDto(number, bonus)

	logger.Debug("update calculation dto", "dto", resultDto)

	if err := c.repCalc.UpdateCalculationResult(ctx, resultDto); err != nil {
		logger.Error("saving calculation result", "error", err, "data", resultDto)
		return
	}

	clearDto := repository.ClearingCalculationQueue{OrderNumber: number}
	if err := c.repCalc.ClearCalculationsQueue(ctx, clearDto); err != nil {
		logger.Error("clearing queue was failed", "error", err, "data", clearDto)
		return
	}

}

func (c CalculationService) calculateProductsBonus(products []service.ProductRow) float64 {
	var bonus float64
	for _, product := range products {
		bonus += c.calculateProduct(product)
	}
	return bonus
}

func (c CalculationService) runNotProcessedOrders(ctx context.Context) error {
	newCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	queue, err := c.repCalc.GetCalculationsQueue(newCtx)
	if err != nil {
		return err
	}

	orders := prepareQueueToExpectedDto(queue)
	for _, order := range orders {
		c.manager.RaiseEvent(NewOrder, order)
	}
	return nil
}

func (c CalculationService) calculateProduct(product service.ProductRow) float64 {
	c.mxRules.Lock()
	defer c.mxRules.Unlock()
	logger.Debug("calculating product", "product", product)
	for _, rule := range c.rules {
		if r := rule.calculate(product.Name, product.Price); r > 0 {
			return r
		}
	}
	return 0
}

func (c CalculationService) readAllRules(ctx context.Context) error {
	if rs, err := c.repRules.Rules(ctx, repository.RuleFilter{}); err == nil {
		return c.fillRules(rs)
	} else {
		return err
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


