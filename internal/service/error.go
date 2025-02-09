package service

import "errors"

//common
//input is wrong
var ErrInvalidFormat = errors.New("invalid format")
var ErrDuplicate = errors.New("login name already exists")
var ErrEntityDoesNotExists = errors.New("entity does not exists")
var ErrLimit = errors.New("limit exceeded")
var ErrNotHaveEnoughPoints = errors.New("not have enough points")
var ErrWrongNumberOfOrder = errors.New("wrong number of order")
var ErrOrderUploadAnotherUser = errors.New("order upload another user")
//AuthorizationService
var ErrWrongNameOrPassword = errors.New("wrong name or password")
//OrderService

//WithdrawalService

//AccrualService

//CalculationService

//CalculationRuleService
