package order

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/service"
)

type checkAccrualJob struct {
	userID   string
	number   string
	attempts int
}

type gettingAccrualConfig struct {
	svc                  service.AccrualService
	jobs                 chan checkAccrualJob
	updateJobs           chan<- updateRepositoryJob //write only
	defaultRestartPeriod time.Duration
}

type accrualManager struct {
	svc service.AccrualService
	/*
		getting information from order service or will restart
	*/
	jobs chan checkAccrualJob
	/*
		getting pause period from child goroutine
	*/
	pauseSignal          chan time.Duration
	defaultRestartPeriod time.Duration
	/*
		write info to updater
	*/
	updateJobs chan<- updateRepositoryJob
	/*
		state manager's child
	*/
	paused *atomic.Bool
	/*
		when child got limit error it lock mutex
		and will wait context's cancel, it will be signal for unlocking mutex
	*/
	waitingPauseMx *sync.Mutex
}

func newAccrualManager(config gettingAccrualConfig) *accrualManager {
	return &accrualManager{
		svc:                  config.svc,
		jobs:                 config.jobs,
		pauseSignal:          make(chan time.Duration),
		defaultRestartPeriod: config.defaultRestartPeriod,
		updateJobs:           config.updateJobs,
		paused:               &atomic.Bool{},
		waitingPauseMx:       &sync.Mutex{},
	}
}

func (m accrualManager) process(ctx context.Context) {
	log := logger.With("component", "accrual manager", "operation", "process")
	//for pause on limit error
	pauseTick := time.NewTicker(time.Microsecond)
	pauseTick.Stop()

	accumulator := make([]checkAccrualJob, 0, 1024)

	childCtx, cancel := context.WithCancel(ctx)

	delayStep := time.Millisecond * 100
	delay := time.Duration(0)

	for {
		select {
		case <-ctx.Done():
			log.Debug("got stop signal. GettingAccrualInformation will stop")
			cancel()
			return
		case job := <-m.jobs:
			log.Debug("got accrual job",
				"order", job.number,
				"user", job.userID,
				"attempts", job.attempts,
			)

			if job.userID == "" || job.number == "" || job.attempts == 0 {
				log.Debug("got invalid job. skip it")
				continue
			}

			//if on pause then add job to accumulator
			//if on fine run checker to get accrual information and
			if m.paused.Load() {
				accumulator = append(accumulator, job)
				log.Debug("accrual manager is on pause. save it to accumulator", "len", len(accumulator))
				continue
			}
			if delay > 0 {
				time.Sleep(delay)
			}
			go m.getAccrualInfo(childCtx, job)

		case pause := <-m.pauseSignal:
			log.Debug("got pause signal. need to stop all jobs")
			m.paused.Store(true)

			delay += delayStep
			log.Debug("increase delay", "current", delay)

			//if there are processes on execution request and they can be successfully, wait little time
			//if they almost failed then where just restart they
			time.Sleep(delay)

			cancel()
			pauseTick.Reset(pause)
		case <-pauseTick.C:
			log.Debug("got pause tick. need to restart all jobs and clean accumulator")

			childCtx, cancel = context.WithCancel(ctx)

			m.paused.Store(false)
			pauseTick.Stop()

			for _, job := range accumulator {
				if delay > 0 {
					time.Sleep(delay)
				}
				go m.getAccrualInfo(childCtx, job)
			}
			accumulator = accumulator[:0]
		}
	}
}

func (m accrualManager) getAccrualInfo(ctx context.Context, job checkAccrualJob) {
	log := logger.With("component", "accrual manager", "operation", "getAccrualInfo")

	select {
	default:
	case <-ctx.Done():
		log.Debug("got stop signal. GettingAccrualInformation will stop")
		return
	}

	result, err := m.svc.Accruals(ctx, service.AccrualsFilterRequest{
		Number: job.number,
	})

	if err == nil {
		m.updateJobs <- updateRepositoryJob{
			userID:      job.userID,
			orderNumber: job.number,
			data:        result,
		}
		return
	}

	retry, raisePause := handleError(err, &job, m.defaultRestartPeriod)
	if raisePause {
		//try to lock mutex if can not it mean other child had notified manager about limit error
		log.Debug("need to lock pause mutex")
		if m.waitingPauseMx.TryLock() {
			log.Debug("pause mutex was locked")

			log.Debug("need to notify manager about limit error")
			m.pauseSignal <- retry
			log.Debug("manager was notified")

			log.Debug("pause mutex was locked, will wait for the stop signal")
			<-ctx.Done()
			log.Debug("got stop signal, will unlock pause mutex")
			m.waitingPauseMx.Unlock()
		}
		log.Debug("restart job after getting limit error", "order", job.number)
		m.jobs <- job
		return
	}
	//order does not exists on accrual service and all attempts was used
	if job.attempts == 0 {
		log.Debug("all attempts was used. save order as invalid", "order", job.number)
		m.updateJobs <- updateRepositoryJob{
			userID:      job.userID,
			orderNumber: job.number,
			data: service.AccrualsInfo{
				OrderNumber: job.number,
				Status:      StatusInvalid,
				Accrual:     0,
			},
		}
		return
	}
	time.AfterFunc(retry, func() {
		log.Debug("restart job after default delay", "order", job.number)
		m.jobs <- job
	})
}

func handleError(err error, job *checkAccrualJob, defaultRetry time.Duration) (retry time.Duration, raisePause bool) {
	log := logger.With(
		"component", "accrual manager",
		"operation", "handleError",
		"order", job.number,
		"user", job.userID,
		"attempts", job.attempts,
	)

	var limitErr service.LimitError
	retry = defaultRetry

	if errors.As(err, &limitErr) {
		log.Error("accrual service is overload", "error", err)
		retry, raisePause = limitErr.RetryAfter, true
	} else if errors.Is(err, service.ErrEntityDoesNotExists) {
		//accrual does not have information about order, may be it will be later
		//but it can mean that we loaded dummy order
		log.Error("accrual service does not have information about order, may be it will be later", "error", err)
		job.attempts--
	} else {
		log.Error("accrual service is not available, task will restart later", "error", err)
	}
	return retry, raisePause
}
