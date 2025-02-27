package gophermart

import "time"

type AuthData struct {
	Login        string
	PasswordHash []byte
}

type UserInfo struct {
	ID           string
	Login        string
	PasswordHash []byte
}

type WithdrawalRequest struct {
	OrderNumber string
	UserID      string
	Sum         float64
}

type TransactionRequest struct {
	UserID string
}

type Transaction struct {
	Income      bool
	UserID      string
	OrderNumber string
	Sum         float64
	CreatedAt   time.Time
}

type OrderCreateRequest struct {
	UserID string
	Number string
}

type OrderUpdateRequest struct {
	UserID  string
	Number  string
	Status  int
	Accrual float64
}

type OrderListRequest struct {
	UserID      string
	OrderNumber string
}

type OrderInfo struct {
	Number    string
	Status    int
	Accrual   float64
	CreatedAt time.Time
}

type AccrualRequest struct {
	OrderNumber string
}

type AccrualInfo struct {
	Number  string
	Status  string
	Accrual float64
}
