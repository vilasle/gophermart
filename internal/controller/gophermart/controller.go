package controller

import (
	"net/http"

	"github.com/vilasle/gophermart/internal/controller"
	"github.com/vilasle/gophermart/internal/service"
)

type Controller struct {
	authSvc service.AuthorizationService
}

// POST /api/user/register
func (c Controller) UserRegister(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// POST /api/user/login
func (c Controller) UserLogin(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// POST /api/user/orders
func (c Controller) RelateOrderWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// GET /api/user/orders
func (c Controller) ListOrdersRelatedWithUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// GET /api/user/balance
func (c Controller) BalanceStateByUser(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// POST /api/user/balance/withdraw
func (c Controller) Withdraw(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}

// GET /api/user/withdrawals
func (c Controller) ListOfWithdrawals(*http.Request) controller.ControllerHandler {
	return func(r *http.Request) controller.Response {
		//TODO implement it
		panic("not implemented")
	}
}
