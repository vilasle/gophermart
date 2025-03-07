package order

import (
	"context"
	"sort"
	"sync"
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
	jobNotice              chan int
	retryOnError           time.Duration
	attemptsGettingAccrual int
	stopMx                 *sync.Mutex
	needStopWorkers        chan time.Duration
	step                   time.Duration
}

type OrderServiceConfig struct {
	gophermart.OrderRepository
	service.AccrualService
	gophermart.WithdrawalRepository
	RetryOnError           time.Duration
	AttemptsGettingAccrual int
}

func NewOrderService(config OrderServiceConfig) OrderService {
	s := OrderService{
		rep:                    config.OrderRepository,
		repTx:                  config.WithdrawalRepository,
		accrual:                config.AccrualService,
		jobs:                   make(chan checkJob, 1024),
		jobNotice:              make(chan int),
		retryOnError:           config.RetryOnError,
		attemptsGettingAccrual: config.AttemptsGettingAccrual,
		needStopWorkers:        make(chan time.Duration),
		stopMx:                 &sync.Mutex{},
		step:                   time.Millisecond * 500,
	}

	return s
}

func (s OrderService) Start(ctx context.Context) error {
	go s.runCheckerDirector(ctx, 0)

	if err := s.runNotProcessedOrders(ctx); err != nil {
		logger.Error("run not processed jobs was failed", "error", err)
		return err
	}

	return nil
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

	s.jobs <- checkJob{
		number:   dto.Number,
		userID:   dto.UserID,
		attempts: s.attemptsGettingAccrual,
	}

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

func (s OrderService) Stop() {
	close(s.jobs)
	close(s.needStopWorkers)
}

func (s OrderService) runNotProcessedOrders(ctx context.Context) error {
	result, err := s.rep.List(ctx, gophermart.OrderListRequest{Status: gophermart.StatusNew})
	if err != nil {
		return err
	}

	go func() {
		for _, order := range result {
			job := checkJob{
				number:   order.Number,
				userID:   order.UserID,
				attempts: s.attemptsGettingAccrual,
			}
			select {
			case <-ctx.Done():
				return
			default:
				s.jobs <- job
			}
		}
	}()

	return nil
}
