// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	reminddetail "github.com/enrico5b1b4/telegram-bot/reminder/reminddetail"
	gomock "github.com/golang/mock/gomock"
)

// MockServicer is a mock of Servicer interface
type MockServicer struct {
	ctrl     *gomock.Controller
	recorder *MockServicerMockRecorder
}

// MockServicerMockRecorder is the mock recorder for MockServicer
type MockServicerMockRecorder struct {
	mock *MockServicer
}

// NewMockServicer creates a new mock instance
func NewMockServicer(ctrl *gomock.Controller) *MockServicer {
	mock := &MockServicer{ctrl: ctrl}
	mock.recorder = &MockServicerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockServicer) EXPECT() *MockServicerMockRecorder {
	return m.recorder
}

// GetReminder mocks base method
func (m *MockServicer) GetReminder(chatID, reminderID int) (*reminddetail.ReminderDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetReminder", chatID, reminderID)
	ret0, _ := ret[0].(*reminddetail.ReminderDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetReminder indicates an expected call of GetReminder
func (mr *MockServicerMockRecorder) GetReminder(chatID, reminderID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetReminder", reflect.TypeOf((*MockServicer)(nil).GetReminder), chatID, reminderID)
}

// DeleteReminder mocks base method
func (m *MockServicer) DeleteReminder(chatID, ID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteReminder", chatID, ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteReminder indicates an expected call of DeleteReminder
func (mr *MockServicerMockRecorder) DeleteReminder(chatID, ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteReminder", reflect.TypeOf((*MockServicer)(nil).DeleteReminder), chatID, ID)
}
