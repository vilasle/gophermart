// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/calculation/repository.go

// Package calculation is a generated GoMock package.
package calculation

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	repository "github.com/vilasle/gophermart/internal/repository/calculation"
)

// MockCalculationRepository is a mock of CalculationRepository interface.
type MockCalculationRepository struct {
	ctrl     *gomock.Controller
	recorder *MockCalculationRepositoryMockRecorder
}

// MockCalculationRepositoryMockRecorder is the mock recorder for MockCalculationRepository.
type MockCalculationRepositoryMockRecorder struct {
	mock *MockCalculationRepository
}

// NewMockCalculationRepository creates a new mock instance.
func NewMockCalculationRepository(ctrl *gomock.Controller) *MockCalculationRepository {
	mock := &MockCalculationRepository{ctrl: ctrl}
	mock.recorder = &MockCalculationRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCalculationRepository) EXPECT() *MockCalculationRepositoryMockRecorder {
	return m.recorder
}

// AddCalculationResult mocks base method.
func (m *MockCalculationRepository) AddCalculationResult(arg0 context.Context, arg1 repository.AddCalculationResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddCalculationResult", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCalculationResult indicates an expected call of AddCalculationResult.
func (mr *MockCalculationRepositoryMockRecorder) AddCalculationResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCalculationResult", reflect.TypeOf((*MockCalculationRepository)(nil).AddCalculationResult), arg0, arg1)
}

// AddCalculationToQueue mocks base method.
func (m *MockCalculationRepository) AddCalculationToQueue(arg0 context.Context, arg1 ...repository.AddingCalculation) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddCalculationToQueue", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddCalculationToQueue indicates an expected call of AddCalculationToQueue.
func (mr *MockCalculationRepositoryMockRecorder) AddCalculationToQueue(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddCalculationToQueue", reflect.TypeOf((*MockCalculationRepository)(nil).AddCalculationToQueue), varargs...)
}

// Calculations mocks base method.
func (m *MockCalculationRepository) Calculations(arg0 context.Context, arg1 repository.CalculationFilter) ([]repository.CalculationInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Calculations", arg0, arg1)
	ret0, _ := ret[0].([]repository.CalculationInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Calculations indicates an expected call of Calculations.
func (mr *MockCalculationRepositoryMockRecorder) Calculations(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Calculations", reflect.TypeOf((*MockCalculationRepository)(nil).Calculations), arg0, arg1)
}

// ClearCalculationsQueue mocks base method.
func (m *MockCalculationRepository) ClearCalculationsQueue(arg0 context.Context, arg1 repository.ClearingCalculationQueue) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ClearCalculationsQueue", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ClearCalculationsQueue indicates an expected call of ClearCalculationsQueue.
func (mr *MockCalculationRepositoryMockRecorder) ClearCalculationsQueue(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ClearCalculationsQueue", reflect.TypeOf((*MockCalculationRepository)(nil).ClearCalculationsQueue), arg0, arg1)
}

// GetCalculationsQueue mocks base method.
func (m *MockCalculationRepository) GetCalculationsQueue(arg0 context.Context) ([]repository.CalculationQueueInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCalculationsQueue", arg0)
	ret0, _ := ret[0].([]repository.CalculationQueueInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCalculationsQueue indicates an expected call of GetCalculationsQueue.
func (mr *MockCalculationRepositoryMockRecorder) GetCalculationsQueue(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCalculationsQueue", reflect.TypeOf((*MockCalculationRepository)(nil).GetCalculationsQueue), arg0)
}

// UpdateCalculationResult mocks base method.
func (m *MockCalculationRepository) UpdateCalculationResult(arg0 context.Context, arg1 repository.AddCalculationResult) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCalculationResult", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCalculationResult indicates an expected call of UpdateCalculationResult.
func (mr *MockCalculationRepositoryMockRecorder) UpdateCalculationResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCalculationResult", reflect.TypeOf((*MockCalculationRepository)(nil).UpdateCalculationResult), arg0, arg1)
}

// MockCalculationRules is a mock of CalculationRules interface.
type MockCalculationRules struct {
	ctrl     *gomock.Controller
	recorder *MockCalculationRulesMockRecorder
}

// MockCalculationRulesMockRecorder is the mock recorder for MockCalculationRules.
type MockCalculationRulesMockRecorder struct {
	mock *MockCalculationRules
}

// NewMockCalculationRules creates a new mock instance.
func NewMockCalculationRules(ctrl *gomock.Controller) *MockCalculationRules {
	mock := &MockCalculationRules{ctrl: ctrl}
	mock.recorder = &MockCalculationRulesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCalculationRules) EXPECT() *MockCalculationRulesMockRecorder {
	return m.recorder
}

// AddRules mocks base method.
func (m *MockCalculationRules) AddRules(arg0 context.Context, arg1 ...repository.AddingRule) (int16, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AddRules", varargs...)
	ret0, _ := ret[0].(int16)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRules indicates an expected call of AddRules.
func (mr *MockCalculationRulesMockRecorder) AddRules(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRules", reflect.TypeOf((*MockCalculationRules)(nil).AddRules), varargs...)
}

// Rules mocks base method.
func (m *MockCalculationRules) Rules(arg0 context.Context, arg1 repository.RuleFilter) ([]repository.RuleInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rules", arg0, arg1)
	ret0, _ := ret[0].([]repository.RuleInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Rules indicates an expected call of Rules.
func (mr *MockCalculationRulesMockRecorder) Rules(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rules", reflect.TypeOf((*MockCalculationRules)(nil).Rules), arg0, arg1)
}
