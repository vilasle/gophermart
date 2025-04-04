package service

import "time"

type CalculationType = int

func DefineCalculationType(t int) (value CalculationType, correct bool) {
	switch t {
	case CalculationTypePercent:
		value, correct = CalculationTypePercent, true
	case CalculationTypeFixed:
		value, correct = CalculationTypeFixed, true
	default:
		value, correct = 0, false
	}
	return
}

const (
	CalculationTypeUnknown CalculationType = iota
	CalculationTypePercent
	CalculationTypeFixed
)

type RegisterRequest struct {
	Login    string
	Password string
}

type AuthorizeRequest struct {
	Login    string
	Password string
}

type UserID struct {
	ID string
}

type UserInfo struct {
	ID string
}

type RegisterOrderRequest struct {
	UserID string
	Number string
}

type ListOrderRequest struct {
	UserID string
}

type OrderInfo struct {
	UserID    string
	Number    string
	Status    string
	Accrual   float64
	CreatedAt time.Time
}

type UserBalanceRequest struct {
	UserID string
}

type UserBalance struct {
	Current   float64
	Withdrawn float64
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
	CreatedAt   time.Time
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
