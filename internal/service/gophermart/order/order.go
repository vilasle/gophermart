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
	number string
}

type OrderService struct {
	rep                  gophermart.OrderRepository
	accrual              service.AccrualService
	jobs                 chan checkJob
	checkerPeriodOnError []time.Duration
}

func NewOrderService(rep gophermart.OrderRepository, accrual service.AccrualService) OrderService {
	s := OrderService{
		rep:                  rep,
		accrual:              accrual,
		jobs:                 make(chan checkJob),
		checkerPeriodOnError: []time.Duration{10 * time.Second, 1 * time.Minute, 5 * time.Minute},
	}

	go s.runCheckerDirector()

	return s
}

func (s OrderService) runCheckerDirector() {
	for job := range s.jobs {
		go s.runChecker(job.number)
	}
}

func (s OrderService) runChecker(number string) {
	result, err := s.getAccrualInformationByOrder(number)
	if err != nil {
		logger.Error("getting information from accrual service was failed", "error", err)

		logger.Debug("may by accrual service is not available, restart task")
		s.jobs <- checkJob{number: number}

		return
	}

	if result.Status == StatusProcessing {
		time.AfterFunc(time.Second*10, func() {
			s.runChecker(number)
		})
	} else {
		s.rep.Update(context.Background(), gophermart.OrderUpdateRequest{
			Number: number,
			Status: result.Status,
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
