package controller

import (
	"errors"
)

var (
	ErrInvalidFormat          = errors.New("invalid format")
	ErrDuplicate              = errors.New("duplicate")
	ErrOrderUploadAnotherUser = errors.New("order upload another user")
	ErrWrongNumberOfOrder     = errors.New("wrong number of order")
	ErrEntityDoesNotExists    = errors.New("entity does not exists")
	ErrLimit                  = errors.New("limit exceeded")
	InternalError             = errors.New("internal error")
	ErrWrongNameOrPassword    = errors.New("wrong name or password")
)
