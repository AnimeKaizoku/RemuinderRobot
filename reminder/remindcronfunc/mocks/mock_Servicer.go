// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/enrico5b1b4/telegram-bot/reminder/remindcronfunc (interfaces: Servicer)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	reminder "github.com/enrico5b1b4/telegram-bot/reminder"
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

// AddReminderRepeatSchedule mocks base method
func (m *MockServicer) AddReminderRepeatSchedule(arg0 *reminder.Reminder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminderRepeatSchedule", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddReminderRepeatSchedule indicates an expected call of AddReminderRepeatSchedule
func (mr *MockServicerMockRecorder) AddReminderRepeatSchedule(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminderRepeatSchedule", reflect.TypeOf((*MockServicer)(nil).AddReminderRepeatSchedule), arg0)
}

// Complete mocks base method
func (m *MockServicer) Complete(arg0 *reminder.Reminder) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Complete", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Complete indicates an expected call of Complete
func (mr *MockServicerMockRecorder) Complete(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Complete", reflect.TypeOf((*MockServicer)(nil).Complete), arg0)
}
