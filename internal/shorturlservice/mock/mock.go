// Code generated by MockGen. DO NOT EDIT.
// Source: databaseInterface.go

// Package mock_shorturlservice is a generated GoMock package.
package mock_shorturlservice

import (
	"context"
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockDatabaseService is a mock of DatabaseService interface.
type MockDatabaseService struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseServiceMockRecorder
}

// MockDatabaseServiceMockRecorder is the mock recorder for MockDatabaseService.
type MockDatabaseServiceMockRecorder struct {
	mock *MockDatabaseService
}

// NewMockDatabaseService creates a new mock instance.
func NewMockDatabaseService(ctrl *gomock.Controller) *MockDatabaseService {
	mock := &MockDatabaseService{ctrl: ctrl}
	mock.recorder = &MockDatabaseServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabaseService) EXPECT() *MockDatabaseServiceMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDatabaseService) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDatabaseServiceMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDatabaseService)(nil).Close))
}

// Connect mocks base method.
func (m *MockDatabaseService) Connect(connStr string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", connStr)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockDatabaseServiceMockRecorder) Connect(connStr interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockDatabaseService)(nil).Connect), connStr)
}

// Ping mocks base method.
func (m *MockDatabaseService) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockDatabaseServiceMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockDatabaseService)(nil).Ping), ctx)
}
