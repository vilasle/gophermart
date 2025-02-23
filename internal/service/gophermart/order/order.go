package order

import (
	"context"
	"errors"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

const (
	StatusNew        = "REGISTERED"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type checkJob struct {
	userID string
	number string
}

type OrderService struct {
	rep                  gophermart.OrderRepository
	repTx                gophermart.WithdrawalRepository
	accrual              service.AccrualService
	jobs                 chan checkJob
	checkerPeriodOnError []time.Duration
}

func NewOrderService(rep gophermart.OrderRepository, accrual service.AccrualService, repTx gophermart.WithdrawalRepository) OrderService {
	s := OrderService{
		rep:                  rep,
		repTx:                repTx,
		accrual:              accrual,
		jobs:                 make(chan checkJob),
		checkerPeriodOnError: []time.Duration{10 * time.Second, 1 * time.Minute, 5 * time.Minute},
	}

	go s.runCheckerDirector()

	return s
}

func (s OrderService) runCheckerDirector() {
	for job := range s.jobs {
		go s.runChecker(job)
	}
}

func (s OrderService) runChecker(job checkJob) {
	result, err := s.getAccrualInformationByOrder(job.number)
	if err != nil {
		logger.Error("getting information from accrual service was failed", "error", err)

		logger.Debug("may by accrual service is not available, restart task")
		s.jobs <- job

		return
	}

	status := gophermart.StatusNew
	switch result.Status {
	case StatusProcessing:
		status = gophermart.StatusProcessing
	case StatusInvalid:
		status = gophermart.StatusInvalid
	case StatusProcessed:
		status = gophermart.StatusProcessed
	}

	upDtp := gophermart.OrderUpdateRequest{
		Number: job.number,
		Status: status,
	}

	ctx := context.Background()
	if err := s.rep.Update(ctx, upDtp); err != nil {
		logger.Error("updating order status was failed", "error", err)
		return
	}

	if result.Status == StatusProcessed {
		inDto := gophermart.WithdrawalRequest{
			UserID:      job.userID,
			OrderNumber: job.number,
			Sum:         result.Accrual,
		}
		if err := s.repTx.Income(ctx, inDto); err != nil {
			logger.Error("adding income to user balance was failed", "error", err)
			return
		}
		return
	}

	if result.Status == StatusProcessing {
		time.AfterFunc(time.Second*10, func() {
			s.runChecker(job)
		})
	}
}

func (s OrderService) getAccrualInformationByOrder(orderNumber string) (result service.AccrualsInfo, err error) {
	for _, period := range s.checkerPeriodOnError {
		result, err = s.accrual.Accruals(context.Background(), service.AccrualsFilterRequest{
			Number: orderNumber,
		})
		if err != nil {
			logger.Error("getting information from accrual service was failed", "error", err)
			logger.Debug("may by accrual service is not available, go to sleep", "sec", period/time.Second)
			time.Sleep(period)
			continue
		} else {
			return result, nil
		}
	}
	return result, err
}

func (s OrderService) Register(ctx context.Context, dto service.RegisterOrderRequest) error {
	if dto.Number == "" {
		return service.ErrInvalidFormat
	}

	rld := gophermart.OrderListRequest{
		UserID:      dto.UserID,
		OrderNumber: dto.Number,
	}
	if result, err := s.rep.List(ctx, rld); err == nil && len(result) > 0 {
		return service.ErrDuplicate
	} else if err != nil {
		return err
	}

	rdt := gophermart.OrderCreateRequest{
		Number: dto.Number,
		UserID: dto.UserID,
	}

	err := s.rep.Create(ctx, rdt)

	if err != nil {
		if errors.Is(err, gophermart.ErrDuplicate) {
			return service.ErrOrderUploadAnotherUser
		}
		return err
	}

	s.jobs <- checkJob{number: dto.Number}
	return nil
}

func (s OrderService) List(ctx context.Context, dto service.ListOrderRequest) ([]service.OrderInfo, error) {
	rld := gophermart.OrderListRequest{
		UserID: dto.UserID,
	}
	result, err := s.rep.List(ctx, rld)
	if err != nil {
		return nil, err
	}

	orders := make([]service.OrderInfo, len(result))
	for i, order := range result {
		orders[i] = service.OrderInfo{
			Number:    order.Number,
			Status:    order.Status,
			Accrual:   order.Accrual,
			CreatedAt: order.CreatedAt,
		}
	}
	return orders, nil
}
