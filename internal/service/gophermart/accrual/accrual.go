package accrual

import (
	"context"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type AccrualServiceHTTP struct {
	rep gophermart.AccrualRepository
}

func NewAccrualServiceHTTP(rep gophermart.AccrualRepository) *AccrualServiceHTTP {
	return &AccrualServiceHTTP{rep: rep}
}

func (s AccrualServiceHTTP) Accruals(ctx context.Context, dto service.AccrualsFilterRequest) (service.AccrualsInfo, error) {
	if dto.Number == "" {
		return service.AccrualsInfo{}, service.ErrInvalidFormat
	}

	result, err := s.rep.AccrualByOrder(ctx, gophermart.AccrualRequest{OrderNumber: dto.Number})
	if err != nil {
		return service.AccrualsInfo{}, err
	}

	return service.AccrualsInfo{
		Status:  result.Status,
		Accrual: result.Accrual,
	}, nil

}
