// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/gophermart/repository.go

// Package authorization is a generated GoMock package.
package authorization

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gophermart "github.com/vilasle/gophermart/internal/repository/gophermart"
)

// MockAuthorizationRepository is a mock of AuthorizationRepository interface.
type MockAuthorizationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAuthorizationRepositoryMockRecorder
}

// MockAuthorizationRepositoryMockRecorder is the mock recorder for MockAuthorizationRepository.
type MockAuthorizationRepositoryMockRecorder struct {
	mock *MockAuthorizationRepository
}

// NewMockAuthorizationRepository creates a new mock instance.
func NewMockAuthorizationRepository(ctrl *gomock.Controller) *MockAuthorizationRepository {
	mock := &MockAuthorizationRepository{ctrl: ctrl}
	mock.recorder = &MockAuthorizationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthorizationRepository) EXPECT() *MockAuthorizationRepositoryMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockAuthorizationRepository) AddUser(arg0 context.Context, arg1 gophermart.AuthData) (gophermart.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(gophermart.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockAuthorizationRepositoryMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockAuthorizationRepository)(nil).AddUser), arg0, arg1)
}

// CheckUser mocks base method.
func (m *MockAuthorizationRepository) CheckUser(arg0 context.Context, arg1 gophermart.AuthData) (gophermart.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUser", arg0, arg1)
	ret0, _ := ret[0].(gophermart.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUser indicates an expected call of CheckUser.
func (mr *MockAuthorizationRepositoryMockRecorder) CheckUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUser", reflect.TypeOf((*MockAuthorizationRepository)(nil).CheckUser), arg0, arg1)
}

// CheckUserByID mocks base method.
func (m *MockAuthorizationRepository) CheckUserByID(arg0 context.Context, arg1 string) (gophermart.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserByID", arg0, arg1)
	ret0, _ := ret[0].(gophermart.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserByID indicates an expected call of CheckUserByID.
func (mr *MockAuthorizationRepositoryMockRecorder) CheckUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserByID", reflect.TypeOf((*MockAuthorizationRepository)(nil).CheckUserByID), arg0, arg1)
}

// MockWithdrawalRepository is a mock of WithdrawalRepository interface.
type MockWithdrawalRepository struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalRepositoryMockRecorder
}

// MockWithdrawalRepositoryMockRecorder is the mock recorder for MockWithdrawalRepository.
type MockWithdrawalRepositoryMockRecorder struct {
	mock *MockWithdrawalRepository
}

// NewMockWithdrawalRepository creates a new mock instance.
func NewMockWithdrawalRepository(ctrl *gomock.Controller) *MockWithdrawalRepository {
	mock := &MockWithdrawalRepository{ctrl: ctrl}
	mock.recorder = &MockWithdrawalRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalRepository) EXPECT() *MockWithdrawalRepositoryMockRecorder {
	return m.recorder
}

// Expense mocks base method.
func (m *MockWithdrawalRepository) Expense(arg0 context.Context, arg1 gophermart.WithdrawalRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Expense", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Expense indicates an expected call of Expense.
func (mr *MockWithdrawalRepositoryMockRecorder) Expense(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Expense", reflect.TypeOf((*MockWithdrawalRepository)(nil).Expense), arg0, arg1)
}

// Income mocks base method.
func (m *MockWithdrawalRepository) Income(arg0 context.Context, arg1 gophermart.WithdrawalRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Income", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Income indicates an expected call of Income.
func (mr *MockWithdrawalRepositoryMockRecorder) Income(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Income", reflect.TypeOf((*MockWithdrawalRepository)(nil).Income), arg0, arg1)
}

// Transactions mocks base method.
func (m *MockWithdrawalRepository) Transactions(arg0 context.Context, arg1 gophermart.TransactionRequest) ([]gophermart.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transactions", arg0, arg1)
	ret0, _ := ret[0].([]gophermart.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transactions indicates an expected call of Transactions.
func (mr *MockWithdrawalRepositoryMockRecorder) Transactions(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transactions", reflect.TypeOf((*MockWithdrawalRepository)(nil).Transactions), arg0, arg1)
}

// MockOrderRepository is a mock of OrderRepository interface.
type MockOrderRepository struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryMockRecorder
}

// MockOrderRepositoryMockRecorder is the mock recorder for MockOrderRepository.
type MockOrderRepositoryMockRecorder struct {
	mock *MockOrderRepository
}

// NewMockOrderRepository creates a new mock instance.
func NewMockOrderRepository(ctrl *gomock.Controller) *MockOrderRepository {
	mock := &MockOrderRepository{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepository) EXPECT() *MockOrderRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockOrderRepository) Create(arg0 context.Context, arg1 gophermart.OrderCreateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockOrderRepositoryMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockOrderRepository)(nil).Create), arg0, arg1)
}

// List mocks base method.
func (m *MockOrderRepository) List(arg0 context.Context, arg1 gophermart.OrderListRequest) ([]gophermart.OrderInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]gophermart.OrderInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockOrderRepositoryMockRecorder) List(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockOrderRepository)(nil).List), arg0, arg1)
}

// Update mocks base method.
func (m *MockOrderRepository) Update(arg0 context.Context, arg1 gophermart.OrderUpdateRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockOrderRepositoryMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockOrderRepository)(nil).Update), arg0, arg1)
}

// MockAccrualRepository is a mock of AccrualRepository interface.
type MockAccrualRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualRepositoryMockRecorder
}

// MockAccrualRepositoryMockRecorder is the mock recorder for MockAccrualRepository.
type MockAccrualRepositoryMockRecorder struct {
	mock *MockAccrualRepository
}

// NewMockAccrualRepository creates a new mock instance.
func NewMockAccrualRepository(ctrl *gomock.Controller) *MockAccrualRepository {
	mock := &MockAccrualRepository{ctrl: ctrl}
	mock.recorder = &MockAccrualRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualRepository) EXPECT() *MockAccrualRepositoryMockRecorder {
	return m.recorder
}

// AccrualByOrder mocks base method.
func (m *MockAccrualRepository) AccrualByOrder(arg0 context.Context, arg1 gophermart.AccrualRequest) (gophermart.AccrualInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AccrualByOrder", arg0, arg1)
	ret0, _ := ret[0].(gophermart.AccrualInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AccrualByOrder indicates an expected call of AccrualByOrder.
func (mr *MockAccrualRepositoryMockRecorder) AccrualByOrder(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AccrualByOrder", reflect.TypeOf((*MockAccrualRepository)(nil).AccrualByOrder), arg0, arg1)
}
