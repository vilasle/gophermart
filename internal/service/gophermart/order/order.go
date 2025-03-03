package order

import (
	"context"
	"sort"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
	"github.com/vilasle/gophermart/internal/tool/order/validation"
)

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
)

type OrderService struct {
	rep                    gophermart.OrderRepository
	repTx                  gophermart.WithdrawalRepository
	accrual                service.AccrualService
	jobs                   chan checkJob
	retryOnError           time.Duration
	attemptsGettingAccrual int
}

type OrderServiceConfig struct {
	gophermart.OrderRepository
	service.AccrualService
	gophermart.WithdrawalRepository
	RetryOnError           time.Duration
	AttemptsGettingAccrual int
}

func NewOrderService(ctx context.Context, config OrderServiceConfig) OrderService {
	s := OrderService{
		rep:                    config.OrderRepository,
		repTx:                  config.WithdrawalRepository,
		accrual:                config.AccrualService,
		jobs:                   make(chan checkJob),
		retryOnError:           config.RetryOnError,
		attemptsGettingAccrual: config.AttemptsGettingAccrual,
	}

	go s.runCheckerDirector()

	err := s.runNotProcessedOrders(ctx)
	if err != nil {
		logger.Error("run not processed jobs was failed", "error", err)
	}

	return s
}

func (s OrderService) Register(ctx context.Context, dto service.RegisterOrderRequest) error {
	if dto.Number == "" || dto.UserID == "" {
		return service.ErrInvalidFormat
	}

	if !validation.IsValidNumber(dto.Number) {
		return service.ErrWrongNumberOfOrder
	}

	if err := s.checkDuplicate(ctx, dto); err != nil {
		return err
	}

	rdt := gophermart.OrderCreateRequest{
		Number: dto.Number,
		UserID: dto.UserID,
	}

	if err := s.rep.Create(ctx, rdt); err != nil {
		return err
	}

	s.startChecker(0, checkJob{
		number:   dto.Number,
		userID:   dto.UserID,
		attempts: s.attemptsGettingAccrual,
	})

	return nil
}

func (s OrderService) checkDuplicate(ctx context.Context, dto service.RegisterOrderRequest) error {
	rld := gophermart.OrderListRequest{OrderNumber: dto.Number}

	if result, err := s.rep.List(ctx, rld); err == nil && len(result) > 0 {
		if result[0].UserID == dto.UserID {
			return service.ErrWasUploadEarly
		}
		return service.ErrDuplicate
	} else if err != nil {
		return err
	}
	return nil
}

func (s OrderService) List(ctx context.Context, dto service.ListOrderRequest) ([]service.OrderInfo, error) {
	if dto.UserID == "" {
		return nil, service.ErrInvalidFormat
	}

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
			Status:    viewOfStatus(order.Status),
			Accrual:   order.Accrual,
			CreatedAt: order.CreatedAt,
		}
	}

	sort.Slice(orders, func(i, j int) bool {
		return orders[i].CreatedAt.Before(orders[j].CreatedAt)
	})

	return orders, nil
}

func viewOfStatus(status int) string {
	switch status {
	case gophermart.StatusNew:
		return StatusNew
	case gophermart.StatusProcessing:
		return StatusProcessing
	case gophermart.StatusInvalid:
		return StatusInvalid
	case gophermart.StatusProcessed:
		return StatusProcessed
	default:
		return ""
	}
}

func (s OrderService) Close() {
	close(s.jobs)
}

func (s OrderService) runCheckerDirector() {
	baseCtx := context.Background()
	ctx, cancel := context.WithCancel(baseCtx)
	defer cancel()

	for job := range s.jobs {
		if job.userID == "" {
			logger.Error("user id is empty")
			continue
		}
		go s.runChecker(ctx, job)
	}
}

func (s OrderService) runNotProcessedOrders(ctx context.Context) error {
	result, err := s.rep.List(ctx, gophermart.OrderListRequest{Status: gophermart.StatusNew})
	if err != nil {
		return err
	}

	go func() {
		for _, order := range result {
			s.jobs <- checkJob{
				number:   order.Number,
				userID:   order.UserID,
				attempts: s.attemptsGettingAccrual,
			}
		}
	}()

	return nil
}

func (s OrderService) runChecker(ctx context.Context, job checkJob) {
	select {
	case <-ctx.Done():
		return
	default:
		result, err := s.getAccrualInformationByOrder(ctx, job.number)

		for _, action := range createActions(result, job, err) {
			action.exec(ctx, &s)
		}
	}
}

func (s OrderService) getAccrualInformationByOrder(ctx context.Context, orderNumber string) (result service.AccrualsInfo, err error) {
	timeCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return s.accrual.Accruals(timeCtx, service.AccrualsFilterRequest{
		Number: orderNumber,
	})
}

func (s OrderService) startChecker(retry time.Duration, job checkJob) {
	time.AfterFunc(retry, func() {
		s.jobs <- job
	})
}

func (s OrderService) postOrderState(ctx context.Context, job checkJob, status int, accrual float64) error {
	dto := gophermart.OrderUpdateRequest{
		UserID:  job.userID,
		Number:  job.number,
		Status:  status,
		Accrual: accrual,
	}
	if err := s.rep.Update(ctx, dto); err != nil {
		logger.Error("updating order status was failed", "error", err)
		return err
	}
	return nil
}

func (s OrderService) postTransaction(ctx context.Context, job checkJob, accrual float64) error {
	inDto := gophermart.WithdrawalRequest{
		UserID:      job.userID,
		OrderNumber: job.number,
		Sum:         accrual,
	}
	return s.repTx.Income(ctx, inDto)
}
