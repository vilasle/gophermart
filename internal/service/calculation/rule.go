package calculation

import (
	"regexp"

	"github.com/vilasle/gophermart/internal/service"
)

type rule struct {
	exp *regexp.Regexp
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
