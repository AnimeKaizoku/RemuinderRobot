// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	time "time"

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

// AddReminderOnDateTime mocks base method
func (m *MockServicer) AddReminderOnDateTime(chatID int, command string, dateTime reminder.DateTime, message string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminderOnDateTime", chatID, command, dateTime, message)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReminderOnDateTime indicates an expected call of AddReminderOnDateTime
func (mr *MockServicerMockRecorder) AddReminderOnDateTime(chatID, command, dateTime, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminderOnDateTime", reflect.TypeOf((*MockServicer)(nil).AddReminderOnDateTime), chatID, command, dateTime, message)
}

// AddReminderOnWordDateTime mocks base method
func (m *MockServicer) AddReminderOnWordDateTime(chatID int, command string, dateTime reminder.WordDateTime, message string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminderOnWordDateTime", chatID, command, dateTime, message)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReminderOnWordDateTime indicates an expected call of AddReminderOnWordDateTime
func (mr *MockServicerMockRecorder) AddReminderOnWordDateTime(chatID, command, dateTime, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminderOnWordDateTime", reflect.TypeOf((*MockServicer)(nil).AddReminderOnWordDateTime), chatID, command, dateTime, message)
}

// AddRepeatableReminderOnDateTime mocks base method
func (m *MockServicer) AddRepeatableReminderOnDateTime(chatID int, command string, dateTime reminder.RepeatableDateTime, message string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddRepeatableReminderOnDateTime", chatID, command, dateTime, message)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddRepeatableReminderOnDateTime indicates an expected call of AddRepeatableReminderOnDateTime
func (mr *MockServicerMockRecorder) AddRepeatableReminderOnDateTime(chatID, command, dateTime, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddRepeatableReminderOnDateTime", reflect.TypeOf((*MockServicer)(nil).AddRepeatableReminderOnDateTime), chatID, command, dateTime, message)
}

// AddReminderIn mocks base method
func (m *MockServicer) AddReminderIn(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminderIn", chatID, command, amountDateTime, message)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReminderIn indicates an expected call of AddReminderIn
func (mr *MockServicerMockRecorder) AddReminderIn(chatID, command, amountDateTime, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminderIn", reflect.TypeOf((*MockServicer)(nil).AddReminderIn), chatID, command, amountDateTime, message)
}

// AddReminderEvery mocks base method
func (m *MockServicer) AddReminderEvery(chatID int, command string, amountDateTime reminder.AmountDateTime, message string) (time.Time, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddReminderEvery", chatID, command, amountDateTime, message)
	ret0, _ := ret[0].(time.Time)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddReminderEvery indicates an expected call of AddReminderEvery
func (mr *MockServicerMockRecorder) AddReminderEvery(chatID, command, amountDateTime, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddReminderEvery", reflect.TypeOf((*MockServicer)(nil).AddReminderEvery), chatID, command, amountDateTime, message)
}
