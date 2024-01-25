// Code generated by MockGen. DO NOT EDIT.
// Source: ./bucket.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dtos "github.com/viniosilva/where-are-my-fruits/internal/dtos"
	models "github.com/viniosilva/where-are-my-fruits/internal/models"
)

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
