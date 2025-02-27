package accrual

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/vilasle/gophermart/internal/controller"
	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/service"
)

type AccrualsInfo struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"` // TODO: omitempty is it ok?
}

type RegisterCalculationReq struct { // TODO: mb use lower case?
	OrderNumber string     `json:"order"`
	Products    []ProductR `json:"goods"`
}

// ProductRow is used to unmarshal data in POST /api/orders
type ProductR struct {
	Name  string  `json:"description"`
	Price float64 `json:"price"`
}

type RegisterCalculationRuleReq struct {
	Match string  `json:"match"`
	Point float64 `json:"reward"`
	Type  string  `json:"reward_type"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Controller struct {
	service.CalculationService
	service.CalculationRuleService
}

// GET /api/orders/{number}
func (c Controller) OrderInfo() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		var (
			calc   service.CalculationInfo
			err    error
			number = chi.URLParam(r, "number")
		)

		log := logger.GetRequestLogger(r)

		if number == "" {
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}

		log.Debug("getting order info", "number", number)

		dto := service.CalculationFilterRequest{OrderNumber: number}

		if calc, err = c.Calculation(r.Context(), dto); err != nil {
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		log.Debug("order calculation", "info", calc)

		return controller.NewResponse(nil, AccrualsInfo{
			OrderNumber: calc.OrderNumber,
			Status:      calc.Status,
			Accrual:     calc.Accrual,
		}, controller.TypeJSON, 0)
	}
}

// POST /api/orders
func (c Controller) RegisterOrder() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		log := logger.GetRequestLogger(r)

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			log.Error("uncorrected request ", "len", len(body), "error", err)

			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}

		log.Info("register order", "body", string(body))

		regReq := RegisterCalculationReq{}

		err = json.Unmarshal(body, &regReq)
		if err != nil {
			log.Error("unmarshal body failed", "error", err)
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		regCalcReq := service.RegisterCalculationRequest{OrderNumber: regReq.OrderNumber}
		for i := range regReq.Products {
			regCalcReq.Products = append(regCalcReq.Products, service.ProductRow{Name: regReq.Products[i].Name, Price: regReq.Products[i].Price})
		}

		log.Debug("register calculation", "request", regCalcReq)

		err = c.CalculationService.Register(r.Context(), regCalcReq)
		if err != nil {
			log.Error("register calculation failed", "error", err)
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		return controller.NewResponse(nil, nil, controller.TypeText, http.StatusAccepted)
	}
}

// POST /api/goods
func (c Controller) AddCalculationRules() controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		log := logger.GetRequestLogger(r)

		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 {
			log.Error("uncorrected request ", "len", len(body), "error", err)
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}
		log.Info("adding rules", "body", string(body))

		prRegCalcRule := RegisterCalculationRuleReq{}
		err = json.Unmarshal(body, &prRegCalcRule)
		if err != nil {
			log.Error("unmarshal body failed", "error", err)
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}

		rewardType := convertRewardType(prRegCalcRule.Type)
		if rewardType == service.CalculationTypeUnknown {
			return controller.NewResponse(service.ErrInvalidFormat, nil, controller.TypeText, 0)
		}

		log.Debug("register calculation rule", "request", prRegCalcRule)
		err = c.CalculationRuleService.Register(r.Context(), service.RegisterCalculationRuleRequest{Match: prRegCalcRule.Match, Point: prRegCalcRule.Point, Type: rewardType})
		if err != nil {
			log.Error("register calculation rule failed", "error", err)
			return controller.NewResponse(err, nil, controller.TypeText, 0)
		}
		return controller.NewResponse(err, nil, controller.TypeText, 0)
	}
}

func convertRewardType(t string) service.CalculationType {
	switch t {
	case "pt":
		return service.CalculationTypeFixed
	case "%":
		return service.CalculationTypePercent
	default:
		return service.CalculationTypeUnknown
	}
}
