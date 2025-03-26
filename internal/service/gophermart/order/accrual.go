package order

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/vilasle/gophermart/internal/logger"
	"github.com/vilasle/gophermart/internal/service"
)

type checkAccrualJob struct {
	userID   string
	number   string
	attempts int
}

type accrualManager struct {
	/*
		getting accrual information
	*/
	accrualSvc service.AccrualService
	/*
		how often need send request and read new jobs
	*/
	readJobsTimeout time.Duration
	/*
		getting unprocessed orders
	*/
	ordersSvc OrderService
	/*
		updating order state and posting transactions
	*/
	updatingOrder updatingOrder
	/*
		interaction job's reader and workers
	*/
	jobs chan checkAccrualJob
	/*
		sleeping timeout when accrual service does not contain order or return another error, except ErrLimit
	*/
	timeoutOnError time.Duration
	/*
		quantity of attempts to get accrual information before set invalid state
		if accrual service does not contain order
	*/
	attemptsOnError int
	/*
		if got ErrLimit workers will wait when can continue work
	*/
	limit   time.Time
	limitMx *sync.Mutex
	wg      *sync.WaitGroup
	//protection from run the same jobs
	mxProcessingMx *sync.Mutex
	orderOnProcess map[string]struct{}
}

type accrualManagerConfig struct {
	accrualSvc      service.AccrualService
	ordersSvc       OrderService
	updatingOrder   updatingOrder
	timeoutOnError  time.Duration
	attemptsOnError int
	readJobsTimeout time.Duration
}

func newAccrualManager(config accrualManagerConfig) *accrualManager {
	return &accrualManager{
		jobs:            make(chan checkAccrualJob),
		limit:           time.Now(),
		limitMx:         &sync.Mutex{},
		accrualSvc:      config.accrualSvc,
		ordersSvc:       config.ordersSvc,
		updatingOrder:   config.updatingOrder,
		timeoutOnError:  config.timeoutOnError,
		attemptsOnError: config.attemptsOnError,
		readJobsTimeout: config.readJobsTimeout,
		wg:              &sync.WaitGroup{},
		mxProcessingMx:  &sync.Mutex{},
		orderOnProcess:  make(map[string]struct{}),
	}
}

func (m *accrualManager) start(ctx context.Context, qtyWorkers int) {
	m.wg.Add(qtyWorkers + 1)
	go m.runJobReader(ctx)

	for i := 0; i < qtyWorkers; i++ {
		go m.runWorker(ctx)
	}

	<-ctx.Done()

	m.stop()
}

/*
the worker listen to channel of jobs, get information about order from accrual service
and update state in order service
*/
func (m *accrualManager) runWorker(ctx context.Context) {
	defer m.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-m.jobs:
			if !ok {
				return
			}
			m.mxProcessingMx.Lock()
			if _, ok := m.orderOnProcess[job.number]; ok {
				m.mxProcessingMx.Unlock()
				continue
			}
			m.orderOnProcess[job.number] = struct{}{}
			m.mxProcessingMx.Unlock()
			m.processJob(ctx, job)
		}
	}
}

func (m *accrualManager) processJob(ctx context.Context, job checkAccrualJob) {
	log := logger.With("component", "accrual manager", "operation", "processJob", "order", job.number)

	defer func() {
		m.mxProcessingMx.Lock()
		delete(m.orderOnProcess, job.number)
		m.mxProcessingMx.Unlock()
	}()

	//set state anyway
	if err := m.updatingOrder.updateOrder(ctx, updateRepositoryJob{
		userID:      job.userID,
		orderNumber: job.number,
		data: service.AccrualsInfo{
			OrderNumber: job.number,
			Status:      StatusProcessing,
			Accrual:     0,
		},
	}); err != nil {
		log.Error("failed to update order", "order", job.number, "error", err)
		return
	}

	for job.attempts > 0 {
		//was there limit error and should it wait?
		now := time.Now()
		if now.Before(m.limit) {
			log.Debug("was limit error. need to wait", "start", now.Format(time.RFC3339), "finish", m.limit.Format(time.RFC3339))
			time.Sleep(m.limit.Sub(now))
		}

		result, err := m.accrualSvc.Accruals(ctx, service.AccrualsFilterRequest{
			Number: job.number,
		})

		retry, raisePause := handleError(err, &job, m.timeoutOnError)

		//does it need to set new limit?
		if raisePause {
			log.Debug("need to lock limit mutex")
			if m.limitMx.TryLock() {
				log.Debug("limit mutex was locked")
				m.limit = time.Now().Add(retry)
				m.limitMx.Unlock()
				log.Debug("new workers limit", "date", m.limit.Format(time.RFC3339))
			}
			continue
		} else if err != nil && retry == 0 {
			//service is not available, job will  restart when reader get one on next time
			break
		} else if retry > 0 {
			now := time.Now()
			log.Debug("order does not found on accrual service, need to wait",
				"start", now.Format(time.RFC3339),
				"finish", now.Add(retry).Format(time.RFC3339),
			)
			time.Sleep(retry)
			continue
		}

		//handler normal situation
		if err := m.updatingOrder.updateOrder(ctx, updateRepositoryJob{
			userID:      job.userID,
			orderNumber: job.number,
			data:        result,
		}); err != nil {
			log.Error("failed to update order", "order", job.number, "error", err)
		}
		break
	}

	if job.attempts == 0 {
		log.Debug("order does not found on accrual service, all attempts was used, mark order as invalid")
		err := m.updatingOrder.updateOrder(ctx, updateRepositoryJob{
			userID:      job.userID,
			orderNumber: job.number,
			data: service.AccrualsInfo{
				OrderNumber: job.number,
				Status:      StatusInvalid,
				Accrual:     0,
			},
		})
		if err != nil {
			log.Error("failed to update order", "order", job.number, "error", err)
		}
	}

}

func handleError(err error, job *checkAccrualJob, defaultRetry time.Duration) (retry time.Duration, raisePause bool) {
	log := logger.With(
		"component", "accrual manager",
		"operation", "handleError",
		"order", job.number,
		"user", job.userID,
		"attempts", job.attempts,
	)

	if err == nil {
		return 0, false
	}

	var limitErr service.LimitError

	if errors.As(err, &limitErr) {
		log.Error("accrual service is overload", "error", err)
		retry, raisePause = limitErr.RetryAfter, true
	} else if errors.Is(err, service.ErrEntityDoesNotExists) {
		//accrual does not have information about order, may be it will be later
		//but it can mean that we loaded dummy order
		log.Error("accrual service does not have information about order, may be it will be later", "error", err)
		job.attempts--
		retry = defaultRetry
	} else {
		log.Error("accrual service is not available, task will restart later", "error", err)
	}
	return retry, raisePause
}

/*
reader by ticker get information from order service about unprocessed orders
and push they to channel of jobs
*/
func (m *accrualManager) runJobReader(ctx context.Context) {
	defer m.wg.Done()

	ticker := time.NewTicker(m.readJobsTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			m.readJobs(ctx)
		}
	}
}

func (m *accrualManager) readJobs(ctx context.Context) {
	orders, err := m.ordersSvc.unprocessedOrders(ctx)
	if err != nil {
		logger.Error("reading unprocessed order failed", "error", err)
		return
	}

	for _, order := range orders {
		select {
		case <-ctx.Done():
			return
		default:
		}

		m.jobs <- checkAccrualJob{
			userID:   order.UserID,
			number:   order.Number,
			attempts: m.attemptsOnError,
		}
	}
}

// waiting when all workers finished and close channels
func (m *accrualManager) stop() {
	m.wg.Wait()
	close(m.jobs)
}
