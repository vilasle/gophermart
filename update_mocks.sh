#!/bin/bash

MOCKBIN=$HOME/go/bin/mockgen

# internal/service/calculation
$MOCKBIN -package=calculation -destination=internal/service/calculation/repository_mock_test.go -source=internal/repository/repository.go