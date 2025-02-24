#!/bin/bash

MOCKBIN=$GOBIN/mockgen

# internal/service
$MOCKBIN -package=controller -destination=internal/controller/accrual/service_mock_test.go -source=internal/service/service.go