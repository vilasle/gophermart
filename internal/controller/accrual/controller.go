package controller

import (
	"net/http"

	"github.com/vilasle/gophermart/internal/controller"
)

type Controller struct {
}

// POST /api/orders/{number}
func (c Controller) OrderInfo(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// POST /api/orders
func (c Controller) RegisterOrder(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// POST /api/goods
func (c Controller) AddCalculationRules(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}