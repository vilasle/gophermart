package withdrawal

import (
	"context"
	"errors"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type WithdrawalService struct {
	rep gophermart.WithdrawalRepository
}

func NewWithdrawalService(rep gophermart.WithdrawalRepository) *WithdrawalService {
	return &WithdrawalService{rep: rep}
}

func (s WithdrawalService) Withdraw(ctx context.Context, dto service.WithdrawalRequest) error {
	if dto.UserID == "" || dto.OrderNumber == "" || dto.Sum == 0 {
		return service.ErrInvalidFormat
	}
	err := s.rep.Expense(ctx, gophermart.WithdrawalRequest{
		UserID:      dto.UserID,
		OrderNumber: dto.OrderNumber,
		Sum:         dto.Sum,
	})

	if errors.Is(err, gophermart.ErrNotEnoughPoints) {
		return service.ErrNotEnoughPoints
	}

	return err
}

func (s WithdrawalService) List(ctx context.Context, dto service.WithdrawalListRequest) ([]service.WithdrawalInfo, error) {
	if dto.UserID == "" {
		return []service.WithdrawalInfo{}, service.ErrInvalidFormat
	}

	r, err := s.rep.Transactions(ctx, gophermart.TransactionRequest{UserID: dto.UserID})
	if err != nil {
		return []service.WithdrawalInfo{}, err
	}
	result := make([]service.WithdrawalInfo, 0, len(r))
	for _, h := range r {
		if h.Income {
			continue
		}

		result = append(result, service.WithdrawalInfo{
			OrderNumber: h.OrderNumber,
			Sum:         h.Sum,
			CreatedAt:   h.CreatedAt,
		})
	}
	return result, nil
}

func (s WithdrawalService) Balance(ctx context.Context, dto service.UserBalanceRequest) (service.UserBalance, error) {
	if dto.UserID == "" {
		return service.UserBalance{}, service.ErrInvalidFormat
	}
	
	r, err := s.rep.Transactions(ctx, gophermart.TransactionRequest{UserID: dto.UserID})
	if err != nil {
		return service.UserBalance{}, err
	}

	balance := service.UserBalance{}
	for _, h := range r {
		if h.Income {
			balance.Current += h.Sum
		} else {
			balance.Withdrawn += h.Sum
		}
	}
	balance.Current -= balance.Withdrawn
	return balance, nil
}
