package controller

import (
	"encoding/json"
	"github.com/vilasle/gophermart/internal/service"
	"io"
	"net/http"

	"github.com/vilasle/gophermart/internal/controller"
)

type Controller struct {
}

// POST /api/orders/{number}
func (c Controller) OrderInfo(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		/////TODO: НУЖНО ЛИ ПРОВЕРЯТЬ CONTENT-LENGTH == 0?
		// orderNum := r.URL.Query().Get("orderNum")
		orderNum := r.PathValue("number")
		// TODO: send orderNum to the validator??
	}
}

// POST /api/orders
func (c Controller) RegisterOrder(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// TODO: create joint getStatusCode?
		}
		// Unmarshal a request
		var calc service.RegisterCalculationRequest
		err = json.Unmarshal(body, &calc)
		if err != nil {
			// TODO: create joint getStatusCode?
		}

	}
}

// POST /api/goods
// Регистрация нового совершённого заказа. Для начисления баллов состав заказа должен быть
// проверен на совпадения с зарегистрированными записями вознаграждений за товары. Начисляется сумма совпадений.
// Принятый заказ не обязан браться в обработку непосредственно в момент получения запроса.
func (c Controller) AddCalculationRules(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//check the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			// TODO: create joint getStatusCode?
		}

		// Unmarshal a request
		var calc service.RegisterCalculationRuleRequest
		err = json.Unmarshal(body, &calc)
		if err != nil {
			// TODO: create joint getStatusCode?
		}
		panic("not implemented")
	}
}
