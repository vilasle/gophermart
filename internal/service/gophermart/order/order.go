package order

import (
	"context"
	"sort"
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

func (s OrderService) Register(ctx context.Context, dto service.RegisterOrderRequest) error {
	if dto.Number == "" || dto.UserID == "" {
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

	if err := s.rep.Create(ctx, rdt); err != nil {
		return err
	}
	s.jobs <- checkJob{number: dto.Number, userID: dto.UserID}
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

func (s OrderService) runChecker(ctx context.Context, job checkJob) {
	select {
	case <-ctx.Done():
		return
	default:
		result, err := s.getAccrualInformationByOrder(job.number)
		if err != nil {
			logger.Error("getting information from accrual service was failed", "error", err)

			logger.Debug("may by accrual service is not available, restart task")
			s.jobs <- job

			return
		}

		status := defineStatus(result)

		upDtp := gophermart.OrderUpdateRequest{
			UserID:  job.userID,
			Number:  job.number,
			Status:  status,
			Accrual: result.Accrual,
		}

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
				s.runChecker(ctx, job)
			})
		}
	}
}

func defineStatus(result service.AccrualsInfo) int {
	switch result.Status {
	case StatusProcessing:
		return gophermart.StatusProcessing
	case StatusInvalid:
		return gophermart.StatusInvalid
	case StatusProcessed:
		return gophermart.StatusProcessed
	}

	return gophermart.StatusNew
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
