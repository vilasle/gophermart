package repository

import "context"

//go:generate mockgen -package=calculation -destination=../service/calculation/repository_mock_test.go -source=repository.go
type CalculationRepository interface {
	AddCalculationToQueue(context.Context, ...AddingCalculation) error
	ClearCalculationsQueue(context.Context, ClearingCalculationQueue) error
	GetCalculationsQueue(context.Context) ([]CalculationQueueInfo, error)

	AddCalculationResult(context.Context, AddCalculationResult) error
	UpdateCalculationResult(context.Context, AddCalculationResult) error

	Calculations(context.Context, CalculationFilter) ([]CalculationInfo, error)
}

type CalculationRules interface {
	AddRules(context.Context, ...AddingRule) (id int16, err error)

	Rules(context.Context, RuleFilter) ([]RuleInfo, error)
}
