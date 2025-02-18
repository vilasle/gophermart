package repository

type CalculationStatus = int

const (
	Registered = iota + 1
	Invalid
	Processing
	Processed
)

func IsCorrectStatus(s CalculationStatus) bool {
	switch s {
	case Registered, Invalid, Processing, Processed:
		return true
	default:
		return false
	}
}

type AddingCalculation struct {
	OrderNumber string
	ProductName string
	Price       float64
}

type AddCalculationResult struct {
	OrderNumber string
	Status      int
	Value       float64
}

type CalculationFilter struct {
	OrderNumber string
}

type CalculationInfo struct {
	OrderNumber string
	Status      CalculationStatus
	Value       float64
}

type AddingRule struct {
	Match           string
	Point           float64
	CalculationType int
}

type RuleFilter struct {
	ID int16
}

type RuleInfo struct {
	ID              int16
	Match           string
	Point           float64
	CalculationType int
}
