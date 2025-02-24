package controller

import (
	"encoding/json"
	"github.com/vilasle/gophermart/internal/controller"
	"github.com/vilasle/gophermart/internal/service"
	"io"
	"net/http"
)

/////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

// AccrualsInf is used as a proxy struct to unmarshal response body in GET /api/orders/{number}
type AccrualsInf struct {
	OrderNumber string  `json:"order"`
	Status      string  `json:"status"`
	Accrual     float64 `json:"accrual,omitempty"` // TODO: omitempty is it ok?
}

// RegisterCalculationRequest is used to unmarshal data in POST /api/orders
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
	Match string                  `json:"match"`
	Point float64                 `json:"reward"` // TODO: тут нет же смысла выставлять omitempty?
	Type  service.CalculationType `json:"reward_type"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

type Controller struct {
	acrService      service.AccrualService
	calcService     service.CalculationService
	calcRuleService service.CalculationRuleService
}

// POST /api/orders/{number}
func (c Controller) OrderInfo(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		if r.ContentLength != 0 { //////TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.TypeText)
		}
		// get an order number
		orderNum := r.PostFormValue("number")
		// get accrual info about th order
		accrualInf, err := c.acrService.Accruals(r.Context(), service.AccrualsFilterRequest{Number: orderNum})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.TypeText)
		}
		//fill proxy-struct to mold response
		accInf := AccrualsInf{OrderNumber: accrualInf.OrderNumber, Status: accrualInf.Status, Accrual: accrualInf.Accrual}
		return controller.NewResponse(nil, accInf, "", controller.TypeJson)
	}
}

// POST /api/orders
func (c Controller) RegisterOrder(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.TypeText)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.TypeText)
		}
		// proxy struct to unmarshal
		regReq := RegisterCalculationReq{}
		// Unmarshal login and password
		err = json.Unmarshal(body, &regReq)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.TypeText)
		}
		// fill an appropriate for service struct
		regCalcReq := service.RegisterCalculationRequest{OrderNumber: regReq.OrderNumber}
		for i := range regReq.Products {
			regCalcReq.Products = append(regCalcReq.Products, service.ProductRow{Name: regReq.Products[i].Name, Price: regReq.Products[i].Price})
		}
		// call service method with a compatible structure (without struct tags)
		err = c.calcService.Register(r.Context(), regCalcReq)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.TypeText)
		}
		return controller.NewResponse(service.StatusOrderSuccessfullyAccepted, nil, "", controller.TypeText) // TODO: Заменить все такие случаи С ERROR на другое тут же нет ошибки
	}
}

// POST /api/goods
// Регистрация нового совершённого заказа. Для начисления баллов состав заказа должен быть
// проверен на совпадения с зарегистрированными записями вознаграждений за товары. Начисляется сумма совпадений.
// Принятый заказ не обязан браться в обработку непосредственно в момент получения запроса.
func (c Controller) AddCalculationRules(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		if r.Body == http.NoBody { // http.NoBody - not nil, len =0
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.TypeText)
		}
		body, err := io.ReadAll(r.Body)
		if err != nil || len(body) == 0 { // TODO: это лишняя проверка?
			return controller.NewResponse(service.ErrInvalidFormat, nil, "", controller.TypeText)
		}
		// fill proxy struct to deserialize
		prRegCalcRule := RegisterCalculationRuleReq{}
		err = json.Unmarshal(body, &prRegCalcRule)
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.TypeText)
		}
		err = c.calcRuleService.Register(r.Context(), service.RegisterCalculationRuleRequest{Match: prRegCalcRule.Match, Point: prRegCalcRule.Point, Type: prRegCalcRule.Type})
		if err != nil {
			return controller.NewResponse(err, nil, "", controller.TypeText)
		}
		return controller.NewResponse(err, nil, "", controller.TypeText)
	}
}
