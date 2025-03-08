package order

import (
	"context"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type updateRepositoryJob struct {
	userID      string
	orderNumber string
	data        service.AccrualsInfo
}

type updateWorkerConfig struct {
	orderRepository       gophermart.OrderRepository
	transactionRepository gophermart.WithdrawalRepository
	quantityOfWorkers     int
	jobs                  <-chan updateRepositoryJob
}

type updateWorker struct {
	orderRepository       gophermart.OrderRepository
	transactionRepository gophermart.WithdrawalRepository
	input                 <-chan updateRepositoryJob
}

// responsibility for updating orders states and commit transactions
func runUpdateWorkers(ctx context.Context, config updateWorkerConfig) {
	log := logger.With("component", "updater")
	log.Debug("starting update workers", "qty", config.quantityOfWorkers)

	for i := 0; i < config.quantityOfWorkers; i++ {
		worker := updateWorker{
			orderRepository:       config.orderRepository,
			transactionRepository: config.transactionRepository,
			input:                 config.jobs,
		}

		go worker.process(ctx)
	}

	log.Debug("update workers started")
}

func (e updateWorker) process(ctx context.Context) {
	log := logger.With("component", "updater instance")

	for {
		select {
		case <-ctx.Done():
			log.Debug("got stop signal. update worker will stop")
			return
		case job := <-e.input:
			log.Debug("got updating job",
				"order", job.orderNumber,
				"userId", job.userID)

			data, status := job.data, defineStatus(job.data.Status)

			go e.postOrderState(ctx, gophermart.OrderUpdateRequest{
				UserID:  job.userID,
				Number:  job.orderNumber,
				Status:  status,
				Accrual: data.Accrual,
			})

			if status == gophermart.StatusProcessed {
				go e.commitTransaction(ctx, gophermart.WithdrawalRequest{
					UserID:      job.userID,
					OrderNumber: job.orderNumber,
					Sum:         data.Accrual,
				})
			}
		}
	}
}

func (e updateWorker) commitTransaction(ctx context.Context, dto gophermart.WithdrawalRequest) {
	log := logger.With("component", "commitTransaction")
	select {
	default:
	case <-ctx.Done():
		log.Debug("context was canceled, committing of transaction will stop")
		return
	}
	if err := e.transactionRepository.Income(ctx, dto); err != nil {
		log.Error("adding income to user balance was failed", "dto", dto, "error", err)
		return
	}
	log.Debug("transaction was committed", "dto", dto)
}

func (e updateWorker) postOrderState(ctx context.Context, dto gophermart.OrderUpdateRequest) {
	log := logger.With("component", "postOrderState")
	select {
	default:
	case <-ctx.Done():
		log.Debug("context was canceled, updater will stop")
		return
	}

	if err := e.orderRepository.Update(ctx, dto); err != nil {
		log.Error("updating order status was failed", "dto", dto, "error", err)
	}
	log.Debug("order status was updated", "dto", dto)
}

func defineStatus(status string) int {
	switch status {
	case StatusProcessing:
		return gophermart.StatusProcessing
	case StatusInvalid:
		return gophermart.StatusInvalid
	case StatusProcessed:
		return gophermart.StatusProcessed
	}
	return gophermart.StatusNew
}
