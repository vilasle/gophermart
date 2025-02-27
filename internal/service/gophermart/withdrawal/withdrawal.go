package withdrawal

import (
	"context"
	"errors"
	"math"
	"sort"
	"strconv"

	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
	"github.com/vilasle/gophermart/internal/tool/order/validation"
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

	n, err := strconv.Atoi(dto.OrderNumber)
	if err != nil {
		return service.ErrWrongNumberOfOrder
	}

	if !validation.ValidNumber(n) {
		return service.ErrWrongNumberOfOrder
	}

	err = s.rep.Expense(ctx, gophermart.WithdrawalRequest{
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
			Sum:         math.Round(h.Sum*100) / 100,
			CreatedAt:   h.CreatedAt,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CreatedAt.Before(result[j].CreatedAt)
	})

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
	balance.Current += balance.Withdrawn

	balance.Current = math.Round(balance.Current*100) / 100
	balance.Withdrawn = math.Round(balance.Withdrawn*100) / 100

	return balance, nil
}
