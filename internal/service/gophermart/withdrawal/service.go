package withdrawal

import (
	"context"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type WithdrawalService struct {
	rep gophermart.WithdrawalRepository
}

func (s WithdrawalService) Withdraw(ctx context.Context, dto service.WithdrawalRequest) error {
	return s.rep.Expense(ctx, gophermart.WithdrawalRequest{
		UserID:      dto.UserID,
		OrderNumber: dto.OrderNumber,
		Sum:         dto.Sum,
	})
}

func (s WithdrawalService) List(ctx context.Context, dto service.WithdrawalListRequest) ([]service.WithdrawalInfo, error) {
	r, err := s.rep.History(ctx, gophermart.HistoryRequest{UserID: dto.UserID})
	if err != nil {
		return []service.WithdrawalInfo{}, err
	}
	request := make([]service.WithdrawalInfo, 0, len(r))
	for _, h := range r {
		if h.Income {
			continue
		}

		request = append(request, service.WithdrawalInfo{
			OrderNumber: h.OrderNumber,
			Sum:         h.Sum,
			CreatedAt:   h.CreatedAt,
		})
	}
	return request, nil
}

func (s WithdrawalService) Balance(ctx context.Context, dto service.UserBalanceRequest) (service.UserBalance, error) {
	r, err := s.rep.History(ctx, gophermart.HistoryRequest{UserID: dto.UserID})
	if err != nil {
		return service.UserBalance{}, err
	}

	balance := service.UserBalance{}
	for _, h := range r {
		if h.Income {
			balance.Balance += h.Sum
		} else {
			balance.Used += h.Sum
		}
	}
	balance.Balance -= balance.Used
	return balance, nil
}
