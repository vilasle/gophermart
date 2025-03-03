package calculation

import (
	"context"
	"regexp"
	"strings"

	repository "github.com/vilasle/gophermart/internal/repository/calculation"
	"github.com/vilasle/gophermart/internal/service"
)

type RuleServiceConfig struct {
	Repository   repository.CalculationRules
	EventManager *EventManager
}

type RuleService struct {
	rep     repository.CalculationRules
	manager *EventManager
}

func NewRuleService(config RuleServiceConfig) *RuleService {
	return &RuleService{
		rep:     config.Repository,
		manager: config.EventManager,
	}
}

func (s RuleService) Register(ctx context.Context, dto service.RegisterCalculationRuleRequest) error {
	//add for ignoring register of letters
	exp := strings.Join([]string{"(?i)", dto.Match}, "")

	//to be sure that exp is correct
	_, err := regexp.Compile(exp)
	if err != nil {
		return err
	}

	id, err := s.rep.AddRules(ctx, repository.AddingRule{
		Match:           dto.Match,
		Point:           dto.Point,
		CalculationType: dto.Type,
	})

	if err != nil {
		if err == repository.ErrDuplicate {
			return service.ErrDuplicate
		}
		return err
	}

	s.manager.RaiseEvent(NewRule, id)

	return nil
}

type rule struct {
	exp             *regexp.Regexp
	calculationType service.CalculationType
	value           float64
}

func (r *rule) calculate(name string, price float64) float64 {
	if r.exp.MatchString(name) {
		switch r.calculationType {
		case service.CalculationTypePercent:
			return price * (r.value / 100)
		case service.CalculationTypeFixed:
			return r.value
		}
	}
	return 0

}
