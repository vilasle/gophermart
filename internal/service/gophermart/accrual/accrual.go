package accrual

import (
	"context"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type AccrualService struct {
	rep gophermart.AccrualRepository
}

func NewAccrualService(rep gophermart.AccrualRepository) *AccrualService {
	return &AccrualService{rep: rep}
}

func (s AccrualService) Accruals(ctx context.Context, dto service.AccrualsFilterRequest) (service.AccrualsInfo, error) {
	if dto.Number == "" {
		return service.AccrualsInfo{}, service.ErrInvalidFormat
	}

	result, err := s.rep.AccrualByOrder(ctx, gophermart.AccrualRequest{OrderNumber: dto.Number})
	if err != nil {
		return service.AccrualsInfo{}, err
	}

	return service.AccrualsInfo{
		OrderNumber: result.Number,
		Status:  result.Status,
		Accrual: result.Accrual,
	}, nil

}
