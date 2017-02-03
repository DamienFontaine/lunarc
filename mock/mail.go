// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/DamienFontaine/lunarc/smtp (interfaces: IMailService)

package mock

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of IMailService interface
type MockIMailService struct {
	ctrl     *gomock.Controller
	recorder *_MockIMailServiceRecorder
}

// Recorder for MockIMailService (not exported)
type _MockIMailServiceRecorder struct {
	mock *MockIMailService
}

func NewMockIMailService(ctrl *gomock.Controller) *MockIMailService {
	mock := &MockIMailService{ctrl: ctrl}
	mock.recorder = &_MockIMailServiceRecorder{mock}
	return mock
}

func (_m *MockIMailService) EXPECT() *_MockIMailServiceRecorder {
	return _m.recorder
}

func (_m *MockIMailService) Send(_param0 string, _param1 string, _param2 string, _param3 string) error {
	ret := _m.ctrl.Call(_m, "Send", _param0, _param1, _param2, _param3)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIMailServiceRecorder) Send(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Send", arg0, arg1, arg2, arg3)
}