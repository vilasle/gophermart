package calculation

import (
	"context"
	"errors"
	"regexp"
	"sync"

	"github.com/vilasle/gophermart/internal/repository"
	"github.com/vilasle/gophermart/internal/service"
)

type CalculationService struct {
	rep     repository.CalculationRepository
	mxRules *sync.Mutex
	rules   map[int16]rule
	manager *EventManager
}

func NewCalculationService(rep repository.CalculationRepository, manager *EventManager) *CalculationService {
	s := &CalculationService{
		rep:     rep,
		manager: manager,
		rules:   make(map[int16]rule),
	}

	s.manager.RegisterHandler(NewOrder, s.calculateOrder)
	s.manager.RegisterHandler(NewRule, s.readRule)

	s.readAllRules()
	return s
}

func (c CalculationService) Register(ctx context.Context, dto service.RegisterCalculationRequest) error {
	//save on db; line on table need for unexpected finishing service
	addingDto := c.prepareAddingDto(dto)

	if err := c.rep.AddCalculationToQueue(ctx, addingDto...); err != nil {
		return err
	}

	//raise event for running worker
	c.manager.RaiseEvent(NewOrder, dto)
	return nil
}

func (c CalculationService) prepareAddingDto(dto service.RegisterCalculationRequest) []repository.AddingCalculation {
	addingDto := make([]repository.AddingCalculation, 0, len(dto.Products))

	orderNumber := dto.OrderNumber

	for _, product := range dto.Products {
		addingDto = append(addingDto, repository.AddingCalculation{
			OrderNumber: orderNumber,
			ProductName: product.Name,
			Price:       product.Price,
		})
	}
	return addingDto
}

func (c CalculationService) Calculation(ctx context.Context, dto service.CalculationFilterRequest) (service.CalculationInfo, error) {
	result, err := c.rep.Calculations(ctx, repository.CalculationFilter{
		OrderNumber: dto.OrderNumber,
	})

	if err != nil {
		return service.CalculationInfo{}, err
	}

	if len(result) == 0 {
		//TODO define right error
		return service.CalculationInfo{}, errors.New("not found")
	}

	calc := result[0]

	return service.CalculationInfo{
		OrderNumber: calc.OrderNumber,
		//TODO wrap status
		Status:  calc.Status,
		Accrual: calc.Value,
	}, nil

}

// event.Type = NewOrder
// event.Data = service.RegisterCalculationRequest
func (c CalculationService) calculateOrder(ctx context.Context, event Event) {
	dto, ok := event.Data.(service.RegisterCalculationRequest)
	if !ok {
		return
	}

	number := dto.OrderNumber
	products := dto.Products

	//TODO what do I need if several rules are matched? Will be stop process when will find matched rule.
	//May be change it

	var bonus float64
	for _, product := range products {
		if bonus = c.calculateProduct(product); bonus > 0 {
			break
		}
	}

	resultDto := c.fillCalculatedDto(number, bonus)

	if err := c.rep.SaveCalculationResult(ctx, resultDto); err != nil {
		//TODO add log message
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
	//TODO change a choice of status
	status := "INVALID"
	if value > 0 {
		status = "PROCESSED"
	}

	return repository.AddCalculationResult{
		OrderNumber: orderNumber,
		Value:       value,
		Status:      status,
	}
}

func (c CalculationService) readAllRules() error {
	rs, err := c.rep.Rules(context.Background(), repository.RuleFilter{})
	if err != nil {
		//TODO wrap error
		return err
	}
	return c.fillRules(rs)
}

// event.Type = NewRule
// event.Data = id rule on repository
func (c CalculationService) readRule(ctx context.Context, event Event) {
	id, ok := event.Data.(int16)
	if !ok {
		return
	}
	rs, err := c.rep.Rules(ctx, repository.RuleFilter{ID: id})
	if err != nil {
		//TODO add log message
		return
	}

	if err := c.fillRules(rs); err != nil {
		//TODO add log message
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
			//TODO add more context
			errs = append(errs, err)
			continue
		}

		calcType := service.CalculationType(r.CalculationType)
		if calcType == 0 {
			//TODO add more context
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
