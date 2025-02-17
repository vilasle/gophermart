package repository

import "context"

//go:generate mockgen -package=calculation -destination=../service/calculation/repository_mock_test.go -source=repository.go
type CalculationRepository interface {
	AddCalculationToQueue(context.Context, ...AddingCalculation) (error)
	
	SaveCalculationResult(context.Context, AddCalculationResult) error
	
	Calculations(context.Context, CalculationFilter) ([]CalculationInfo, error)
	
	AddRules(context.Context, ...AddingRule) error
	
	Rules(context.Context, RuleFilter) ([]RuleInfo, error)
}	
