package gophermart

import "context"

type AuthorizationRepository interface {
	AddUser(context.Context, AuthData) (UserInfo, error)
	CheckUser(context.Context, AuthData) (UserInfo, error)
}

type WithdrawalRepository interface {
	Expense(context.Context, WithdrawalRequest) error
	Transactions(context.Context, TransactionRequest) ([]Transaction, error)
} 

type OrderRepository interface {
	Create(context.Context, OrderCreateRequest) error
	Update(context.Context, OrderUpdateRequest) error
	List(context.Context, OrderListRequest) ([]OrderInfo, error)
}