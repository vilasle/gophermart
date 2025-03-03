package calculation

import(
	"math"
	"github.com/vilasle/gophermart/internal/service"
	repository "github.com/vilasle/gophermart/internal/repository/calculation"
)

// dto
func prepareCalculatedDto(orderNumber string, value float64) repository.AddCalculationResult {
	status := repository.Invalid
	if value > 0 {
		status = repository.Processed
	}

	return repository.AddCalculationResult{
		OrderNumber: orderNumber,
		Value:       math.Round(value*100) / 100,
		Status:      status,
	}
}

func prepareQueueToExpectedDto(dto []repository.CalculationQueueInfo) []service.RegisterCalculationRequest {
	m := make(map[string][]service.ProductRow)
	for _, v := range dto {
		if _, ok := m[v.OrderNumber]; !ok {
			m[v.OrderNumber] = make([]service.ProductRow, 0)
		}

		m[v.OrderNumber] = append(m[v.OrderNumber], service.ProductRow{
			Name:  v.ProductName,
			Price: v.Price,
		})
	}

	result := make([]service.RegisterCalculationRequest, 0, len(m))
	for k, v := range m {
		result = append(result, service.RegisterCalculationRequest{
			OrderNumber: k,
			Products:    v,
		})
	}
	return result
}

func prepareAddingDto(dto service.RegisterCalculationRequest) ([]repository.AddingCalculation, repository.AddCalculationResult) {
	addingDto := make([]repository.AddingCalculation, 0, len(dto.Products))

	orderNumber := dto.OrderNumber

	for _, product := range dto.Products {
		addingDto = append(addingDto, repository.AddingCalculation{
			OrderNumber: orderNumber,
			ProductName: product.Name,
			Price:       product.Price,
		})
	}

	calcDto := repository.AddCalculationResult{
		OrderNumber: orderNumber,
		Status:      repository.Registered,
		Value:       0,
	}

	return addingDto, calcDto
}

func prepareCalculatedInfo(dto repository.CalculationInfo) service.CalculationInfo {
	return service.CalculationInfo{
		OrderNumber: dto.OrderNumber,
		Status:      statusView(dto.Status),
		Accrual:     math.Round(dto.Value*100) / 100,
	}
}

func statusView(status repository.CalculationStatus) string {
	switch status {
	case repository.Invalid:
		return "INVALID"
	case repository.Processing:
		return "PROCESSING"
	case repository.Processed:
		return "PROCESSED"
	default:
		return ""
	}
}
