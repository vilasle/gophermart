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

func NewOrderService(config OrderServiceConfig) OrderService {
	s := OrderService{
		rep:                    config.OrderRepository,
		repTx:                  config.WithdrawalRepository,
		accrual:                config.AccrualService,
		retryOnError:           config.RetryOnError,
		attemptsGettingAccrual: config.AttemptsGettingAccrual,
	}

	return s
}

func (s OrderService) Start(ctx context.Context) {
	log := logger.With("component", "OrderService")
	log.Info("starting service")

	//getting accrual information
	managerConfig := accrualManagerConfig{
		accrualSvc: s.accrual,
		ordersSvc:  s,
		updatingOrder: updatingOrder{
			orderRepository:       s.rep,
			transactionRepository: s.repTx,
		},
		timeoutOnError:  s.retryOnError,
		attemptsOnError: s.attemptsGettingAccrual,
		readJobsTimeout: time.Second * 5,
	}

	manager := newAccrualManager(managerConfig)

	go manager.start(ctx, 5)
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

func (s OrderService) unprocessedOrders(ctx context.Context) ([]service.OrderInfo, error) {
	result, err := s.rep.List(ctx, gophermart.OrderListRequest{
		Status: []int{gophermart.StatusNew, gophermart.StatusProcessing},
		Limit:  50,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]service.OrderInfo, len(result))
	for i, order := range result {
		orders[i] = service.OrderInfo{
			UserID:    order.UserID,
			Number:    order.Number,
			Status:    viewOfStatus(order.Status),
			Accrual:   order.Accrual,
			CreatedAt: order.CreatedAt,
		}
	}
	return orders, nil
}
