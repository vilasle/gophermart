package calculation

import (
	"context"
	"regexp"
	"strings"

	"github.com/vilasle/gophermart/internal/repository"
	"github.com/vilasle/gophermart/internal/service"
)

type RuleServiceConfig struct {
	CalculationRules repository.CalculationRules
	EventManager     *EventManager
}

type RuleService struct {
	rep     repository.CalculationRules
	manager *EventManager
}

func NewRuleService(config RuleServiceConfig) *RuleService {
	return &RuleService{
		rep:     config.CalculationRules,
		manager: config.EventManager,
	}
}

func (s RuleService) Register(ctx context.Context, dto service.RegisterCalculationRuleRequest) error {
	//add for ignoring register of letters
	exp := strings.Join([]string{"(?i)", dto.Match}, "")

	//to be sure that exp is correct
	_, err := regexp.Compile(exp)
	if err != nil {
		//TODO wrap error
		return err
	}

	id, err := s.rep.AddRules(ctx, repository.AddingRule{
		Match:           dto.Match,
		Point:           dto.Point,
		CalculationType: dto.Type,
	})

	if err != nil {
		//TODO wrap error
		return err
	}

	s.manager.RaiseEvent(NewRule, id)

	return nil

}
