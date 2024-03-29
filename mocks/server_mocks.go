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

// Delete mocks base method.
func (m *MockBucketController) Delete(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", ctx)
}

// Delete indicates an expected call of Delete.
func (mr *MockBucketControllerMockRecorder) Delete(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockBucketController)(nil).Delete), ctx)
}

// List mocks base method.
func (m *MockBucketController) List(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "List", ctx)
}

// List indicates an expected call of List.
func (mr *MockBucketControllerMockRecorder) List(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockBucketController)(nil).List), ctx)
}

// MockFruitController is a mock of FruitController interface.
type MockFruitController struct {
	ctrl     *gomock.Controller
	recorder *MockFruitControllerMockRecorder
}

// MockFruitControllerMockRecorder is the mock recorder for MockFruitController.
type MockFruitControllerMockRecorder struct {
	mock *MockFruitController
}

// NewMockFruitController creates a new mock instance.
func NewMockFruitController(ctrl *gomock.Controller) *MockFruitController {
	mock := &MockFruitController{ctrl: ctrl}
	mock.recorder = &MockFruitControllerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFruitController) EXPECT() *MockFruitControllerMockRecorder {
	return m.recorder
}

// AddOnBucket mocks base method.
func (m *MockFruitController) AddOnBucket(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AddOnBucket", ctx)
}

// AddOnBucket indicates an expected call of AddOnBucket.
func (mr *MockFruitControllerMockRecorder) AddOnBucket(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOnBucket", reflect.TypeOf((*MockFruitController)(nil).AddOnBucket), ctx)
}

// Create mocks base method.
func (m *MockFruitController) Create(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Create", ctx)
}

// Create indicates an expected call of Create.
func (mr *MockFruitControllerMockRecorder) Create(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFruitController)(nil).Create), ctx)
}

// Delete mocks base method.
func (m *MockFruitController) Delete(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Delete", ctx)
}

// Delete indicates an expected call of Delete.
func (mr *MockFruitControllerMockRecorder) Delete(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFruitController)(nil).Delete), ctx)
}

// RemoveFromBucket mocks base method.
func (m *MockFruitController) RemoveFromBucket(ctx *gin.Context) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveFromBucket", ctx)
}

// RemoveFromBucket indicates an expected call of RemoveFromBucket.
func (mr *MockFruitControllerMockRecorder) RemoveFromBucket(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBucket", reflect.TypeOf((*MockFruitController)(nil).RemoveFromBucket), ctx)
}
