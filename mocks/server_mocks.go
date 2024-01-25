// Code generated by MockGen. DO NOT EDIT.
// Source: ./server.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gin "github.com/gin-gonic/gin"
	gomock "github.com/golang/mock/gomock"
)

// MockHealthController is a mock of HealthController interface.
type MockHealthController struct {
	ctrl     *gomock.Controller
	recorder *MockHealthControllerMockRecorder
}

// MockHealthControllerMockRecorder is the mock recorder for MockHealthController.
type MockHealthControllerMockRecorder struct {
	mock *MockHealthController
}

// NewMockHealthController creates a new mock instance.
func NewMockHealthController(ctrl *gomock.Controller) *MockHealthController {
	mock := &MockHealthController{ctrl: ctrl}
	mock.recorder = &MockHealthControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthController) EXPECT() *MockHealthControllerMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockHealthController) Check(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Check", ctx)
}

// Check indicates an expected call of Check.
func (mr *MockHealthControllerMockRecorder) Check(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockHealthController)(nil).Check), ctx)
}

// MockBucketController is a mock of BucketController interface.
type MockBucketController struct {
	ctrl     *gomock.Controller
	recorder *MockBucketControllerMockRecorder
}

// MockBucketControllerMockRecorder is the mock recorder for MockBucketController.
type MockBucketControllerMockRecorder struct {
	mock *MockBucketController
}

// NewMockBucketController creates a new mock instance.
func NewMockBucketController(ctrl *gomock.Controller) *MockBucketController {
	mock := &MockBucketController{ctrl: ctrl}
	mock.recorder = &MockBucketControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucketController) EXPECT() *MockBucketControllerMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockBucketController) Create(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", ctx)
}

// Create indicates an expected call of Create.
func (mr *MockBucketControllerMockRecorder) Create(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBucketController)(nil).Create), ctx)
}
