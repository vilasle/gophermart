package gophermart

import "context"

type AuthorizationRepository interface {
	AddUser(context.Context, AuthData) (UserInfo, error)
	CheckUser(context.Context, AuthData) (UserInfo, error)
}

type WithdrawalRepository interface {
	Expense(context.Context, WithdrawalRequest) error
	History(context.Context, HistoryRequest) ([]HistoryLine, error)
} 