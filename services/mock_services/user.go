// Automatically generated by MockGen. DO NOT EDIT!
// Source: user.go

package mock_services

import (
	models "github.com/DamienFontaine/lunarc/models"
	gomock "github.com/golang/mock/gomock"
)

// Mock of IUserService interface
type MockIUserService struct {
	ctrl     *gomock.Controller
	recorder *_MockIUserServiceRecorder
}

// Recorder for MockIUserService (not exported)
type _MockIUserServiceRecorder struct {
	mock *MockIUserService
}

func NewMockIUserService(ctrl *gomock.Controller) *MockIUserService {
	mock := &MockIUserService{ctrl: ctrl}
	mock.recorder = &_MockIUserServiceRecorder{mock}
	return mock
}

func (_m *MockIUserService) EXPECT() *_MockIUserServiceRecorder {
	return _m.recorder
}

func (_m *MockIUserService) GetByID(id string) (models.User, error) {
	ret := _m.ctrl.Call(_m, "GetByID", id)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIUserServiceRecorder) GetByID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetByID", arg0)
}

func (_m *MockIUserService) Get(username string, password string) (models.User, error) {
	ret := _m.ctrl.Call(_m, "Get", username, password)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIUserServiceRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get", arg0, arg1)
}

func (_m *MockIUserService) Add(user models.User) (models.User, error) {
	ret := _m.ctrl.Call(_m, "Add", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIUserServiceRecorder) Add(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Add", arg0)
}

func (_m *MockIUserService) FindAll() ([]models.User, error) {
	ret := _m.ctrl.Call(_m, "FindAll")
	ret0, _ := ret[0].([]models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIUserServiceRecorder) FindAll() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FindAll")
}

func (_m *MockIUserService) Delete(user models.User) error {
	ret := _m.ctrl.Call(_m, "Delete", user)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIUserServiceRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockIUserService) Update(id string, user models.User) error {
	ret := _m.ctrl.Call(_m, "Update", id, user)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIUserServiceRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1)
}
