package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	mart "github.com/vilasle/gophermart/internal/repository/gophermart"
	"github.com/vilasle/gophermart/internal/service"
)

type AccrualRepository struct {
	addr *url.URL
}

func NewAccrualRepository(addr *url.URL) *AccrualRepository {
	return &AccrualRepository{addr: addr}
}

func (r AccrualRepository) AccrualByOrder(ctx context.Context, dto mart.AccrualRequest) (mart.AccrualInfo, error) {
	var (
		addr = r.addr.JoinPath("orders", dto.OrderNumber).String()
		cl   = &http.Client{Timeout: 10 * time.Second}
		req  *http.Request
		resp *http.Response
		err  error
	)

	if req, err = http.NewRequestWithContext(ctx, http.MethodGet, addr, nil); err != nil {
		return mart.AccrualInfo{}, err
	}

	if resp, err = cl.Do(req); err != nil {
		return mart.AccrualInfo{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		err := createLimitError(resp)
		return mart.AccrualInfo{}, err
	} else if resp.StatusCode == http.StatusNoContent {
		return mart.AccrualInfo{}, service.ErrEntityDoesNotExists
	} else if resp.StatusCode != http.StatusOK {
		return mart.AccrualInfo{}, service.ErrUnexpected
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return mart.AccrualInfo{}, err
	}

	return prepareAccrualInfo(content)
}

func prepareAccrualInfo(content []byte) (mart.AccrualInfo, error) {
	result := struct {
		Order   string  `json:"order"`
		Status  string  `json:"status"`
		Accrual float64 `json:"accrual"`
	}{}

	if err := json.Unmarshal(content, &result); err != nil {
		return mart.AccrualInfo{}, service.ErrUnexpected
	}

	return mart.AccrualInfo{
		Number:  result.Order,
		Status:  result.Status,
		Accrual: result.Accrual,
	}, nil
}

func createLimitError(resp *http.Response) error {
	retry := resp.Header.Get("Retry-After")
	if retry == "" {
		return service.ErrLimit
	}
	if pause, err := strconv.Atoi(retry); err == nil {
		return service.LimitError{
			RetryAfter: time.Duration(pause) * time.Second,
		}
	}
	return service.ErrLimit
}
