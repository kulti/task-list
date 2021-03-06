// Code generated by MockGen. DO NOT EDIT.
// Source: sprintstore.go

// Package sprintstore_test is a generated GoMock package.
package sprintstore_test

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	storages "github.com/kulti/task-list/server/internal/storages"
	reflect "reflect"
)

// MockDBStore is a mock of dbStore interface
type MockDBStore struct {
	ctrl     *gomock.Controller
	recorder *MockDBStoreMockRecorder
}

// MockDBStoreMockRecorder is the mock recorder for MockDBStore
type MockDBStoreMockRecorder struct {
	mock *MockDBStore
}

// NewMockDBStore creates a new mock instance
func NewMockDBStore(ctrl *gomock.Controller) *MockDBStore {
	mock := &MockDBStore{ctrl: ctrl}
	mock.recorder = &MockDBStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDBStore) EXPECT() *MockDBStoreMockRecorder {
	return m.recorder
}

// NewSprint mocks base method
func (m *MockDBStore) NewSprint(ctx context.Context, opts storages.SprintOpts) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewSprint", ctx, opts)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewSprint indicates an expected call of NewSprint
func (mr *MockDBStoreMockRecorder) NewSprint(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewSprint", reflect.TypeOf((*MockDBStore)(nil).NewSprint), ctx, opts)
}

// CreateTask mocks base method
func (m *MockDBStore) CreateTask(ctx context.Context, task storages.Task, sprintID string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTask", ctx, task, sprintID)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateTask indicates an expected call of CreateTask
func (mr *MockDBStoreMockRecorder) CreateTask(ctx, task, sprintID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTask", reflect.TypeOf((*MockDBStore)(nil).CreateTask), ctx, task, sprintID)
}

// ListTasks mocks base method
func (m *MockDBStore) ListTasks(ctx context.Context, sprintID string) (storages.TaskList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListTasks", ctx, sprintID)
	ret0, _ := ret[0].(storages.TaskList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListTasks indicates an expected call of ListTasks
func (mr *MockDBStoreMockRecorder) ListTasks(ctx, sprintID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListTasks", reflect.TypeOf((*MockDBStore)(nil).ListTasks), ctx, sprintID)
}
