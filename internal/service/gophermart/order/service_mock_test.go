// Code generated by MockGen. DO NOT EDIT.
// Source: internal/service/service.go

// Package order is a generated GoMock package.
package order

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	service "github.com/vilasle/gophermart/internal/service"
)

// MockAuthorizationService is a mock of AuthorizationService interface.
type MockAuthorizationService struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationServiceMockRecorder
}

// MockAuthorizationServiceMockRecorder is the mock recorder for MockAuthorizationService.
type MockAuthorizationServiceMockRecorder struct {
	mock *MockAuthorizationService
}

// NewMockAuthorizationService creates a new mock instance.
func NewMockAuthorizationService(ctrl *gomock.Controller) *MockAuthorizationService {
	mock := &MockAuthorizationService{ctrl: ctrl}
	mock.recorder = &MockAuthorizationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationService) EXPECT() *MockAuthorizationServiceMockRecorder {
	return m.recorder
}

// Authorize mocks base method.
func (m *MockAuthorizationService) Authorize(arg0 context.Context, arg1 service.AuthorizeRequest) (service.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authorize", arg0, arg1)
	ret0, _ := ret[0].(service.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authorize indicates an expected call of Authorize.
func (mr *MockAuthorizationServiceMockRecorder) Authorize(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authorize", reflect.TypeOf((*MockAuthorizationService)(nil).Authorize), arg0, arg1)
}

// CheckByUserID mocks base method.
func (m *MockAuthorizationService) CheckByUserID(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckByUserID", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckByUserID indicates an expected call of CheckByUserID.
func (mr *MockAuthorizationServiceMockRecorder) CheckByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckByUserID", reflect.TypeOf((*MockAuthorizationService)(nil).CheckByUserID), arg0, arg1)
}

// Register mocks base method.
func (m *MockAuthorizationService) Register(arg0 context.Context, arg1 service.RegisterRequest) (service.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(service.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockAuthorizationServiceMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAuthorizationService)(nil).Register), arg0, arg1)
}

// MockOrderService is a mock of OrderService interface.
type MockOrderService struct {
	ctrl     *gomock.Controller
	recorder *MockOrderServiceMockRecorder
}

// MockOrderServiceMockRecorder is the mock recorder for MockOrderService.
type MockOrderServiceMockRecorder struct {
	mock *MockOrderService
}

// NewMockOrderService creates a new mock instance.
func NewMockOrderService(ctrl *gomock.Controller) *MockOrderService {
	mock := &MockOrderService{ctrl: ctrl}
	mock.recorder = &MockOrderServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderService) EXPECT() *MockOrderServiceMockRecorder {
	return m.recorder
}

// List mocks base method.
func (m *MockOrderService) List(arg0 context.Context, arg1 service.ListOrderRequest) ([]service.OrderInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]service.OrderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockOrderServiceMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockOrderService)(nil).List), arg0, arg1)
}

// Register mocks base method.
func (m *MockOrderService) Register(arg0 context.Context, arg1 service.RegisterOrderRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockOrderServiceMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockOrderService)(nil).Register), arg0, arg1)
}

// MockWithdrawalService is a mock of WithdrawalService interface.
type MockWithdrawalService struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalServiceMockRecorder
}

// MockWithdrawalServiceMockRecorder is the mock recorder for MockWithdrawalService.
type MockWithdrawalServiceMockRecorder struct {
	mock *MockWithdrawalService
}

// NewMockWithdrawalService creates a new mock instance.
func NewMockWithdrawalService(ctrl *gomock.Controller) *MockWithdrawalService {
	mock := &MockWithdrawalService{ctrl: ctrl}
	mock.recorder = &MockWithdrawalServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalService) EXPECT() *MockWithdrawalServiceMockRecorder {
	return m.recorder
}

// Balance mocks base method.
func (m *MockWithdrawalService) Balance(arg0 context.Context, arg1 service.UserBalanceRequest) (service.UserBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Balance", arg0, arg1)
	ret0, _ := ret[0].(service.UserBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Balance indicates an expected call of Balance.
func (mr *MockWithdrawalServiceMockRecorder) Balance(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Balance", reflect.TypeOf((*MockWithdrawalService)(nil).Balance), arg0, arg1)
}

// List mocks base method.
func (m *MockWithdrawalService) List(arg0 context.Context, arg1 service.WithdrawalListRequest) ([]service.WithdrawalInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]service.WithdrawalInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockWithdrawalServiceMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockWithdrawalService)(nil).List), arg0, arg1)
}

// Withdraw mocks base method.
func (m *MockWithdrawalService) Withdraw(arg0 context.Context, arg1 service.WithdrawalRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Withdraw", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Withdraw indicates an expected call of Withdraw.
func (mr *MockWithdrawalServiceMockRecorder) Withdraw(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Withdraw", reflect.TypeOf((*MockWithdrawalService)(nil).Withdraw), arg0, arg1)
}

// MockAccrualService is a mock of AccrualService interface.
type MockAccrualService struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualServiceMockRecorder
}

// MockAccrualServiceMockRecorder is the mock recorder for MockAccrualService.
type MockAccrualServiceMockRecorder struct {
	mock *MockAccrualService
}

// NewMockAccrualService creates a new mock instance.
func NewMockAccrualService(ctrl *gomock.Controller) *MockAccrualService {
	mock := &MockAccrualService{ctrl: ctrl}
	mock.recorder = &MockAccrualServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualService) EXPECT() *MockAccrualServiceMockRecorder {
	return m.recorder
}

// Accruals mocks base method.
func (m *MockAccrualService) Accruals(arg0 context.Context, arg1 service.AccrualsFilterRequest) (service.AccrualsInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Accruals", arg0, arg1)
	ret0, _ := ret[0].(service.AccrualsInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Accruals indicates an expected call of Accruals.
func (mr *MockAccrualServiceMockRecorder) Accruals(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Accruals", reflect.TypeOf((*MockAccrualService)(nil).Accruals), arg0, arg1)
}

// MockCalculationService is a mock of CalculationService interface.
type MockCalculationService struct {
	ctrl     *gomock.Controller
	recorder *MockCalculationServiceMockRecorder
}

// MockCalculationServiceMockRecorder is the mock recorder for MockCalculationService.
type MockCalculationServiceMockRecorder struct {
	mock *MockCalculationService
}

// NewMockCalculationService creates a new mock instance.
func NewMockCalculationService(ctrl *gomock.Controller) *MockCalculationService {
	mock := &MockCalculationService{ctrl: ctrl}
	mock.recorder = &MockCalculationServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCalculationService) EXPECT() *MockCalculationServiceMockRecorder {
	return m.recorder
}

// Calculation mocks base method.
func (m *MockCalculationService) Calculation(arg0 context.Context, arg1 service.CalculationFilterRequest) (service.CalculationInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calculation", arg0, arg1)
	ret0, _ := ret[0].(service.CalculationInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calculation indicates an expected call of Calculation.
func (mr *MockCalculationServiceMockRecorder) Calculation(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calculation", reflect.TypeOf((*MockCalculationService)(nil).Calculation), arg0, arg1)
}

// Register mocks base method.
func (m *MockCalculationService) Register(arg0 context.Context, arg1 service.RegisterCalculationRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockCalculationServiceMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCalculationService)(nil).Register), arg0, arg1)
}

// MockCalculationRuleService is a mock of CalculationRuleService interface.
type MockCalculationRuleService struct {
	ctrl     *gomock.Controller
	recorder *MockCalculationRuleServiceMockRecorder
}

// MockCalculationRuleServiceMockRecorder is the mock recorder for MockCalculationRuleService.
type MockCalculationRuleServiceMockRecorder struct {
	mock *MockCalculationRuleService
}

// NewMockCalculationRuleService creates a new mock instance.
func NewMockCalculationRuleService(ctrl *gomock.Controller) *MockCalculationRuleService {
	mock := &MockCalculationRuleService{ctrl: ctrl}
	mock.recorder = &MockCalculationRuleServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCalculationRuleService) EXPECT() *MockCalculationRuleServiceMockRecorder {
	return m.recorder
}

// Register mocks base method.
func (m *MockCalculationRuleService) Register(arg0 context.Context, arg1 service.RegisterCalculationRuleRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register.
func (mr *MockCalculationRuleServiceMockRecorder) Register(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockCalculationRuleService)(nil).Register), arg0, arg1)
}
