package gophermart

import "errors"

var ErrDuplicate = errors.New("duplicate")
var ErrEmptyResult = errors.New("empty result")
var ErrNotEnoughPoints = errors.New("does not enough points")