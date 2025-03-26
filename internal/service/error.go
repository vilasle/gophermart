package service

import (
	"errors"
	"fmt"
	"time"
)

var ErrInvalidFormat = errors.New("invalid format")
var ErrWasUploadEarly = errors.New("was upload early")
var ErrDuplicate = errors.New("login name already exists")
var ErrEntityDoesNotExists = errors.New("entity does not exists")
var ErrLimit = errors.New("limit exceeded")
var ErrNotEnoughPoints = errors.New("not have enough points")
var ErrWrongNumberOfOrder = errors.New("wrong number of order")
var ErrOrderUploadAnotherUser = errors.New("order upload another user")
var ErrWrongNameOrPassword = errors.New("wrong name or password")
var ErrUnexpected = errors.New("unexpected error")

type LimitError struct {
	RetryAfter time.Duration
}

func (e LimitError) Error() string {
	return fmt.Sprintf("too many requests, try it again in %d second", e.RetryAfter/time.Second)
}
