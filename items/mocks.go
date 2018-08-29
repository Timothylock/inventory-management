// Automatically generated by MockGen. DO NOT EDIT!
// Source: items.go

package items

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Persister interface
type MockPersister struct {
	ctrl     *gomock.Controller
	recorder *_MockPersisterRecorder
}

// Recorder for MockPersister (not exported)
type _MockPersisterRecorder struct {
	mock *MockPersister
}

func NewMockPersister(ctrl *gomock.Controller) *MockPersister {
	mock := &MockPersister{ctrl: ctrl}
	mock.recorder = &_MockPersisterRecorder{mock}
	return mock
}

func (_m *MockPersister) EXPECT() *_MockPersisterRecorder {
	return _m.recorder
}

func (_m *MockPersister) MoveItem(ID string, direction string) error {
	ret := _m.ctrl.Call(_m, "MoveItem", ID, direction)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPersisterRecorder) MoveItem(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "MoveItem", arg0, arg1)
}

func (_m *MockPersister) DeleteItem(ID string) error {
	ret := _m.ctrl.Call(_m, "DeleteItem", ID)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPersisterRecorder) DeleteItem(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "DeleteItem", arg0)
}
