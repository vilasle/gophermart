#!/bin/bash

MOCKBIN=$(go env GOBIN)/mockgen

# internal/service/calculation
$MOCKBIN -package=calculation -destination=internal/service/calculation/repository_mock_test.go -source=internal/repository/calculation/repository.go
# internal/service/gophermart/
$MOCKBIN -package=accrual -destination=internal/service/gophermart/accrual/repository_mock_test.go -source=internal/repository/gophermart/repository.go
$MOCKBIN -package=authorization -destination=internal/service/gophermart/authorization/repository_mock_test.go -source=internal/repository/gophermart/repository.go
$MOCKBIN -package=withdrawal -destination=internal/service/gophermart/withdrawal/repository_mock_test.go -source=internal/repository/gophermart/repository.go

$MOCKBIN -package=order -destination=internal/service/gophermart/order/repository_mock_test.go -source=internal/repository/gophermart/repository.go
$MOCKBIN -package=order -destination=internal/service/gophermart/order/service_mock_test.go -source=internal/service/service.go
