package repository

import "context"

type CalculationRepository interface {
	AddCalculationToQueue(context.Context, ...AddingCalculation) (error)
	
	SaveCalculationResult(context.Context, AddCalculationResult) error
	
	Calculations(context.Context, CalculationFilter) ([]CalculationInfo, error)
	
	AddRules(context.Context, ...AddingRule) error
	
	Rules(context.Context, RuleFilter) ([]RuleInfo, error)
}	
