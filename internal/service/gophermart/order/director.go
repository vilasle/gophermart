package order

import (
	"context"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type jobState struct {
	id             int
	gotAccrualInfo bool
	cancel         context.CancelFunc
}

func (s OrderService) runCheckerDirector(ctx context.Context, step time.Duration) {
	log := logger.With("component", "checker_director")
	//keep states of checkers and don't cancel all checker if one of them is failed
	//if got signal from main context, when cancel all of them
	//if one of failed on interaction with accrual service,
	//some checkers could be success and we don't need to interrupt them
	states := make(map[int]*jobState)
	var counter int
	for {
		select {
		case <-ctx.Done():
			//stop all of jobs
			log.Debug("got stop signal")
			log.Debug("checker director finished")
			for _, state := range states {
				state.cancel()
			}

			return
		case job := <-s.jobs:
			log.Debug("got job for handling", "job", job)
			if job.userID == "" {
				log.Error("user id is empty")
				continue
			}
			counter++
			newCtx, cancel := context.WithCancel(ctx)
			states[counter] = &jobState{
				id:     counter,
				cancel: cancel,
			}
			log.Debug("run checker", "id", counter)
			go s.runChecker(newCtx, states[counter], job)
			if step > 0 {
				time.Sleep(step)
			}
		case jobId := <-s.jobNotice:
			log.Debug("finished job", "jobId", jobId)
			if states[jobId] != nil {
				states[jobId].cancel()
			}
			delete(states, jobId)

			if step > 0 {
				step -= step
			}

		case pause := <-s.needStopWorkers:
			//do not stop succeeded jobs
			for _, state := range states {
				if !state.gotAccrualInfo {
					state.cancel()
				}
				delete(states, state.id)
			}

			log.Debug("got pause signal", "pause(sec)", pause/time.Second)
			time.AfterFunc(pause, func() {
				log.Debug("pause finished. restart director and jobs")
				newStep := step + s.step
				log.Debug("increase step", "step(ms)", newStep/time.Millisecond)
				go s.runCheckerDirector(ctx, newStep)
				go s.runNotProcessedOrders(ctx)
			})
			return
		}
	}

}

func (s OrderService) runChecker(ctx context.Context, state *jobState, job checkJob) {
	log := logger.With("component", "checker")

	select {
	default:
	case <-ctx.Done():
		log.Debug("got cancel signal, stop job", "job", job)
		return
	}

	result, err := s.getAccrualInformationByOrder(ctx, job.number)
	//mark for director, don't cancel this context
	if err == nil {
		log.Debug("got accrual information, mark state as succeed", "state", state)
		state.gotAccrualInfo = true
	}

	for _, action := range createActions(result, job, err) {
		action.exec(ctx, &s)
	}
	s.jobNotice <- state.id
}

func (s OrderService) getAccrualInformationByOrder(ctx context.Context, orderNumber string) (result service.AccrualsInfo, err error) {
	timeCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()

	return s.accrual.Accruals(timeCtx, service.AccrualsFilterRequest{
		Number: orderNumber,
	})
}

func (s OrderService) startChecker(ctx context.Context, retry time.Duration, job checkJob) {
	log := logger.With("component", "restart_checker")

	time.AfterFunc(retry, func() {
		select {
		default:
		case <-ctx.Done():
			log.Debug("got cancel signal, do not restart job", "job", job)
			return
		}
		log.Debug("restart job", "job", job)
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
