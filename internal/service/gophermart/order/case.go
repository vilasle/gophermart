package order

import (
	"context"
	"errors"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type checkJob struct {
	userID   string
	number   string
	attempts int
}

type checkerAction interface {
	exec(context.Context, *OrderService)
}

func createActions(result service.AccrualsInfo, job checkJob, err error) []checkerAction {
	handlers := make([]checkerAction, 0)

	if err != nil {
		handlers = append(handlers, onError{job: job, err: err})
	}

	if err == nil {
		handlers = append(handlers, updateState{job: job, result: result})
	}

	if result.Status == StatusProcessing {
		handlers = append(handlers, onProcessing{job: job})
	}

	if result.Status == StatusProcessed {
		handlers = append(handlers, onProcessed{job: job, result: result})
	}

	return handlers
}

type onError struct {
	job checkJob
	err error
}

func (h onError) exec(ctx context.Context, svc *OrderService) {
	log := logger.With("component", "onErrorAction")
	if h.err == nil {
		return
	}

	select {
	default:
	case <-ctx.Done():
		return
	}

	retry, needPause := handlerAccrualError(&h.job, h.err, svc.retryOnError)
	if needPause {
		if !svc.stopMx.TryLock() {
			return
		}
		log.Debug("accrual service is overloaded, need to stop workers")
		svc.needStopWorkers <- retry

		svc.stopMx.Unlock()
		return
	}
	if h.job.attempts == 0 {
		if err := svc.postOrderState(ctx, h.job, gophermart.StatusInvalid, 0); err != nil {
			log.Error("updating order status was failed", "error", err)
		}
	} else {
		svc.startChecker(ctx, retry, h.job)
	}

}

func handlerAccrualError(job *checkJob, err error, defaultRetry time.Duration) (retry time.Duration, needPause bool) {
	log := logger.With("component", "handlerAccrualError")
	log.Error("getting information from accrual service was failed", "error", err)

	var limitErr service.LimitError
	retry = defaultRetry

	if errors.As(err, &limitErr) {
		retry = limitErr.RetryAfter
		needPause = true
		log.Error("accrual service is overloaded, restart task",
			"order", job.number,
			"retryAfter", retry)
	} else if errors.Is(err, service.ErrEntityDoesNotExists) {
		job.attempts--
		log.Error("order does not exist on accrual service, may be it will be later, restart task",
			"order", job.number,
			"retryAfter", retry)
	} else {
		log.Error("may by accrual service is not available, restart task",
			"order", job.number,
			"retryAfter", retry)
	}
	return retry, needPause

}

type updateState struct {
	result service.AccrualsInfo
	job    checkJob
}

func (h updateState) exec(ctx context.Context, svc *OrderService) {
	log := logger.With("component", "updateState")
	if err := svc.postOrderState(ctx, h.job, defineStatus(h.result), h.result.Accrual); err != nil {
		log.Error("updating order status was failed", "error", err)
		return
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

type onProcessing struct {
	job checkJob
}

func (h onProcessing) exec(ctx context.Context, svc *OrderService) {
	svc.startChecker(ctx, svc.retryOnError, h.job)
}

type onProcessed struct {
	result service.AccrualsInfo
	job    checkJob
}

func (h onProcessed) exec(ctx context.Context, svc *OrderService) {
	log := logger.With("component", "onProcessedAction")
	//this is successful situation but 
	//if on this moment will cancel of context we could not save result 
	//because will ignore
	
	if err := svc.postTransaction(ctx, h.job, h.result.Accrual); err != nil {
		log.Error("adding income to user balance was failed", "error", err)
		return
	}
}
