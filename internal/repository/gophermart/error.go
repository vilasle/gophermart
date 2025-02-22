package gophermart

import "errors"

var ErrDuplicate = errors.New("duplicate")
var ErrEmptyResult = errors.New("empty result")