// Code generated by MockGen. DO NOT EDIT.
// Source: ./interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dtos "github.com/viniosilva/where-are-my-fruits/internal/dtos"
	models "github.com/viniosilva/where-are-my-fruits/internal/models"
)

// MockHealthService is a mock of HealthService interface.
type MockHealthService struct {
	ctrl     *gomock.Controller
	recorder *MockHealthServiceMockRecorder
}

// MockHealthServiceMockRecorder is the mock recorder for MockHealthService.
type MockHealthServiceMockRecorder struct {
	mock *MockHealthService
}

// NewMockHealthService creates a new mock instance.
func NewMockHealthService(ctrl *gomock.Controller) *MockHealthService {
	mock := &MockHealthService{ctrl: ctrl}
	mock.recorder = &MockHealthServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHealthService) EXPECT() *MockHealthServiceMockRecorder {
	return m.recorder
}

// Check mocks base method.
func (m *MockHealthService) Check(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Check", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Check indicates an expected call of Check.
func (mr *MockHealthServiceMockRecorder) Check(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Check", reflect.TypeOf((*MockHealthService)(nil).Check), ctx)
}

// MockBucketService is a mock of BucketService interface.
type MockBucketService struct {
	ctrl     *gomock.Controller
	recorder *MockBucketServiceMockRecorder
}

// MockBucketServiceMockRecorder is the mock recorder for MockBucketService.
type MockBucketServiceMockRecorder struct {
	mock *MockBucketService
}

// NewMockBucketService creates a new mock instance.
func NewMockBucketService(ctrl *gomock.Controller) *MockBucketService {
	mock := &MockBucketService{ctrl: ctrl}
	mock.recorder = &MockBucketServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBucketService) EXPECT() *MockBucketServiceMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockBucketService) Create(ctx context.Context, data dtos.CreateBucketDto) (*models.Bucket, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].(*models.Bucket)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockBucketServiceMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockBucketService)(nil).Create), ctx, data)
}

// MockFruitService is a mock of FruitService interface.
type MockFruitService struct {
	ctrl     *gomock.Controller
	recorder *MockFruitServiceMockRecorder
}

// MockFruitServiceMockRecorder is the mock recorder for MockFruitService.
type MockFruitServiceMockRecorder struct {
	mock *MockFruitService
}

// NewMockFruitService creates a new mock instance.
func NewMockFruitService(ctrl *gomock.Controller) *MockFruitService {
	mock := &MockFruitService{ctrl: ctrl}
	mock.recorder = &MockFruitServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFruitService) EXPECT() *MockFruitServiceMockRecorder {
	return m.recorder
}

// AddOnBucket mocks base method.
func (m *MockFruitService) AddOnBucket(ctx context.Context, fruitID, bucketID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOnBucket", ctx, fruitID, bucketID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddOnBucket indicates an expected call of AddOnBucket.
func (mr *MockFruitServiceMockRecorder) AddOnBucket(ctx, fruitID, bucketID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOnBucket", reflect.TypeOf((*MockFruitService)(nil).AddOnBucket), ctx, fruitID, bucketID)
}

// Create mocks base method.
func (m *MockFruitService) Create(ctx context.Context, data dtos.CreateFruitDto) (*models.Fruit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, data)
	ret0, _ := ret[0].(*models.Fruit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockFruitServiceMockRecorder) Create(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockFruitService)(nil).Create), ctx, data)
}

// Delete mocks base method.
func (m *MockFruitService) Delete(ctx context.Context, fruitID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, fruitID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockFruitServiceMockRecorder) Delete(ctx, fruitID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockFruitService)(nil).Delete), ctx, fruitID)
}

// RemoveFromBucket mocks base method.
func (m *MockFruitService) RemoveFromBucket(ctx context.Context, fruitID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveFromBucket", ctx, fruitID)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveFromBucket indicates an expected call of RemoveFromBucket.
func (mr *MockFruitServiceMockRecorder) RemoveFromBucket(ctx, fruitID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveFromBucket", reflect.TypeOf((*MockFruitService)(nil).RemoveFromBucket), ctx, fruitID)
}
