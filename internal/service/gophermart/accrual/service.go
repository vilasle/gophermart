package accrual

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/vilasle/gophermart/internal/service"
)

type AccrualServiceHTTP struct {
	addr *url.URL
}

func (s AccrualServiceHTTP) Accruals(ctx context.Context, dto service.AccrualsFilterRequest) (service.AccrualsInfo, error) {
	if dto.Number == "" {
		return service.AccrualsInfo{}, service.ErrInvalidFormat
	}

	numberAddr := s.addr.JoinPath(dto.Number)

	cl := &http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, numberAddr.String(), nil)
	if err != nil {
		return service.AccrualsInfo{}, err
	}

	resp, err := cl.Do(req)
	if err != nil {
		return service.AccrualsInfo{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		return service.AccrualsInfo{}, service.ErrLimit
	} else if resp.StatusCode == http.StatusNoContent {
		return service.AccrualsInfo{}, service.ErrEntityDoesNotExists
	} else if resp.StatusCode != http.StatusOK {
		return service.AccrualsInfo{}, service.ErrUnexpected
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return service.AccrualsInfo{}, err
	}

	result := service.AccrualsInfo{}

	if err := json.Unmarshal(content, &result); err != nil {
		return service.AccrualsInfo{}, service.ErrUnexpected
	}
	return result, nil
}
