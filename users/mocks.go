// Code generated by MockGen. DO NOT EDIT.
// Source: users.go

// Package mock_users is a generated GoMock package.
package users

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersister is a mock of Persister interface
type MockPersister struct {
	ctrl     *gomock.Controller
	recorder *MockPersisterMockRecorder
}

// MockPersisterMockRecorder is the mock recorder for MockPersister
type MockPersisterMockRecorder struct {
	mock *MockPersister
}

// NewMockPersister creates a new mock instance
func NewMockPersister(ctrl *gomock.Controller) *MockPersister {
	mock := &MockPersister{ctrl: ctrl}
	mock.recorder = &MockPersisterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPersister) EXPECT() *MockPersisterMockRecorder {
	return m.recorder
}

// GetUser mocks base method
func (m *MockPersister) GetUser(username, password string) (User, error) {
	ret := m.ctrl.Call(m, "GetUser", username, password)
	ret0, _ := ret[0].(User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockPersisterMockRecorder) GetUser(username, password interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockPersister)(nil).GetUser), username, password)
}

// IsValidToken mocks base method
func (m *MockPersister) IsValidToken(token string) (bool, error) {
	ret := m.ctrl.Call(m, "IsValidToken", token)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsValidToken indicates an expected call of IsValidToken
func (mr *MockPersisterMockRecorder) IsValidToken(token interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsValidToken", reflect.TypeOf((*MockPersister)(nil).IsValidToken), token)
}
