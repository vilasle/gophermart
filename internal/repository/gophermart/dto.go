package gophermart

import "time"

type AuthData struct {
	Login        string
	PasswordHash string
}

type UserInfo struct {
	ID           string
	Login        string
	PasswordHash string
}

type WithdrawalRequest struct {
	OrderNumber string
	UserID      string
	Sum         float64
}

type HistoryRequest struct {
	UserID string
}

type HistoryLine struct {
	Income      bool
	UserId      string
	OrderNumber string
	Sum         float64
	CreatedAt   time.Time
}
