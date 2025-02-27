package service

import "errors"

var ErrInvalidFormat = errors.New("invalid format")
var ErrDuplicate = errors.New("login name already exists")
var ErrEntityDoesNotExists = errors.New("entity does not exists")
var ErrLimit = errors.New("limit exceeded")
var ErrNotEnoughPoints = errors.New("not have enough points")
var ErrWrongNumberOfOrder = errors.New("wrong number of order")
var ErrOrderUploadAnotherUser = errors.New("order upload another user")
var ErrWrongNameOrPassword = errors.New("wrong name or password")
var ErrUnexpected = errors.New("unexpected error")
