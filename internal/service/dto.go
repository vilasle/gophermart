package service

import "time"

type CalculationType = int

const (
	CalculationTypePercent CalculationType = iota + 1
	CalculationTypeSpecificValue
)

type RegisterRequest struct {
	Login    string
	Password string
}

type AuthorizeRequest struct {
	Login    string
	Password string
}

type UserInfo struct {
	ID string
}

type RegisterOrderRequest struct {
	Number string
}

type ListOrderRequest struct {
	UserID string
}

type OrderInfo struct {
	Number    string
	Status    string
	Accrual   float64
	CreatedAt time.Time
}

type UserBalanceRequest struct {
	UserID string
}

type UserBalance struct {
	Balance float64
	Used    float64
}

type WithdrawalRequest struct {
	UserID      string
	OrderNumber string
	Sum         float64
}

type WithdrawalListRequest struct {
	UserID string
}

type WithdrawalInfo struct {
	OrderNumber string
	Sum         float64
	Status      string
}

type AccrualsFilterRequest struct {
	Number string
}

type AccrualsInfo struct {
	OrderNumber string
	Status      string
	Accrual     float64
}

type RegisterCalculationRequest struct {
	OrderNumber string
	Products    []ProductRow
}

type ProductRow struct {
	Name  string
	Price float64
}

type CalculationFilterRequest struct {
	OrderNumber string
}

type CalculationInfo struct {
	OrderNumber string
	Status      string
	Accrual     float64
}

type RegisterCalculationRuleRequest struct {
	Match string
	Point float64
	Type  CalculationType
}
