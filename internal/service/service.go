package service

import "context"

type AuthorizationService interface {
	//can return defined errors ErrInvalidFormat, ErrDuplicate and undefined error
	Register(context.Context, RegisterRequest) (UserInfo, error)
	//can return defined errors ErrInvalidFormat, ErrWrongNameOrPassword and undefined error
	Authorize(context.Context, AuthorizeRequest) (UserInfo, error)
	//can return defined errors                              and undefined error
	CheckByUserID(context.Context, string) error
}

type OrderService interface {
	//can return defined errors ErrInvalidFormat, ErrDuplicate, ErrOrderUploadAnotherUser, ErrWrongNumberOfOrder
	//and undefined error
	Register(context.Context, RegisterOrderRequest) error
	//can return undefined error
	List(context.Context, ListOrderRequest) ([]OrderInfo, error)
	Close()
}

type WithdrawalService interface {
	//can return undefined error  // TODO: mb we should add error  402 — на счету недостаточно средств?
	Withdraw(context.Context, WithdrawalRequest) error
	//can return defined errors ErrInvalidFormat, ErrDuplicate and undefined error
	List(context.Context, WithdrawalListRequest) ([]WithdrawalInfo, error)
	//can return undefined error
	Balance(context.Context, UserBalanceRequest) (UserBalance, error)
}

type AccrualService interface {
	//can return defined errors ErrEntityDoesNotExists, ErrLimit, ErrInvalidFormat, ErrUnexpected and undefined error
	Accruals(context.Context, AccrualsFilterRequest) (AccrualsInfo, error)
}

type CalculationService interface {
	//Add new calculations to queue. Can return defined errors ErrEntityDoesNotExists and undefined error
	Register(context.Context, RegisterCalculationRequest) error
	//Return information about calculation. Can return defined errors ErrInvalidFormat, ErrDuplicate and undefined error
	Calculation(context.Context, CalculationFilterRequest) (CalculationInfo, error)
}

type CalculationRuleService interface {
	//can return defined errors ErrInvalidFormat, ErrDuplicate and undefined error
	Register(context.Context, RegisterCalculationRuleRequest) error
}
