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

type updatingOrder struct {
	orderRepository       gophermart.OrderRepository
	transactionRepository gophermart.WithdrawalRepository
}

func (e updatingOrder) updateOrder(ctx context.Context, job updateRepositoryJob) error {
	log := logger.With("component", "updater instance")
	log.Debug("got updating job",
		"order", job.orderNumber,
		"userId", job.userID)

	data, status := job.data, defineStatus(job.data.Status)

	if err := e.postOrderState(ctx, gophermart.OrderUpdateRequest{
		UserID:  job.userID,
		Number:  job.orderNumber,
		Status:  status,
		Accrual: data.Accrual,
	}); err != nil {
		return err
	}

	if status == gophermart.StatusProcessed {
		if err := e.commitTransaction(ctx, gophermart.WithdrawalRequest{
			UserID:      job.userID,
			OrderNumber: job.orderNumber,
			Sum:         data.Accrual,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (e updatingOrder) commitTransaction(ctx context.Context, dto gophermart.WithdrawalRequest) error{
	log := logger.With("component", "commitTransaction")
	select {
	default:
	case <-ctx.Done():
		log.Debug("context was canceled, committing of transaction will stop")
		return ctx.Err()
	}
	if err := e.transactionRepository.Income(ctx, dto); err != nil {
		log.Error("adding income to user balance was failed", "dto", dto, "error", err)
		return err
	}
	log.Debug("transaction was committed", "dto", dto)
	return nil
}

func (e updatingOrder) postOrderState(ctx context.Context, dto gophermart.OrderUpdateRequest) error {
	log := logger.With("component", "postOrderState")
	select {
	default:
	case <-ctx.Done():
		log.Debug("context was canceled, updater will stop")
		return ctx.Err()
	}

	if err := e.orderRepository.Update(ctx, dto); err != nil {
		log.Error("updating order status was failed", "dto", dto, "error", err)
		return err
	}
	log.Debug("order status was updated", "dto", dto)
	return nil
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
