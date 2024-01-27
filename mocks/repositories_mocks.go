// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockDB is a mock of DB interface.
type MockDB struct {
	ctrl     *gomock.Controller
	recorder *MockDBMockRecorder
}

// MockDBMockRecorder is the mock recorder for MockDB.
type MockDBMockRecorder struct {
	mock *MockDB
}

// NewMockDB creates a new mock instance.
func NewMockDB(ctrl *gomock.Controller) *MockDB {
	mock := &MockDB{ctrl: ctrl}
	mock.recorder = &MockDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDB) EXPECT() *MockDBMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockDB) Create(value interface{}) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", value)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDBMockRecorder) Create(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDB)(nil).Create), value)
}

// Model mocks base method.
func (m *MockDB) Model(value interface{}) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Model", value)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Model indicates an expected call of Model.
func (mr *MockDBMockRecorder) Model(value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Model", reflect.TypeOf((*MockDB)(nil).Model), value)
}

// Transaction mocks base method.
func (m *MockDB) Transaction(fc func(*gorm.DB) error, opts ...*sql.TxOptions) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{fc}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Transaction", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Transaction indicates an expected call of Transaction.
func (mr *MockDBMockRecorder) Transaction(fc interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{fc}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transaction", reflect.TypeOf((*MockDB)(nil).Transaction), varargs...)
}

// MockSQL is a mock of SQL interface.
type MockSQL struct {
	ctrl     *gomock.Controller
	recorder *MockSQLMockRecorder
}

// MockSQLMockRecorder is the mock recorder for MockSQL.
type MockSQLMockRecorder struct {
	mock *MockSQL
}

// NewMockSQL creates a new mock instance.
func NewMockSQL(ctrl *gomock.Controller) *MockSQL {
	mock := &MockSQL{ctrl: ctrl}
	mock.recorder = &MockSQLMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSQL) EXPECT() *MockSQLMockRecorder {
	return m.recorder
}

// PingContext mocks base method.
func (m *MockSQL) PingContext(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PingContext", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// PingContext indicates an expected call of PingContext.
func (mr *MockSQLMockRecorder) PingContext(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PingContext", reflect.TypeOf((*MockSQL)(nil).PingContext), ctx)
}
