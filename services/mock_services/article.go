// Automatically generated by MockGen. DO NOT EDIT!
// Source: article.go

package mock_services

import (
	models "github.com/DamienFontaine/lunarc/models"
	gomock "github.com/golang/mock/gomock"
)

// Mock of IArticleService interface
type MockIArticleService struct {
	ctrl     *gomock.Controller
	recorder *_MockIArticleServiceRecorder
}

// Recorder for MockIArticleService (not exported)
type _MockIArticleServiceRecorder struct {
	mock *MockIArticleService
}

func NewMockIArticleService(ctrl *gomock.Controller) *MockIArticleService {
	mock := &MockIArticleService{ctrl: ctrl}
	mock.recorder = &_MockIArticleServiceRecorder{mock}
	return mock
}

func (_m *MockIArticleService) EXPECT() *_MockIArticleServiceRecorder {
	return _m.recorder
}

func (_m *MockIArticleService) GetByID(id string) (models.Article, error) {
	ret := _m.ctrl.Call(_m, "GetByID", id)
	ret0, _ := ret[0].(models.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIArticleServiceRecorder) GetByID(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetByID", arg0)
}

func (_m *MockIArticleService) GetByPretty(pretty string) (models.Article, error) {
	ret := _m.ctrl.Call(_m, "GetByPretty", pretty)
	ret0, _ := ret[0].(models.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIArticleServiceRecorder) GetByPretty(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "GetByPretty", arg0)
}

func (_m *MockIArticleService) FindByStatus(status string) ([]models.Article, error) {
	ret := _m.ctrl.Call(_m, "FindByStatus", status)
	ret0, _ := ret[0].([]models.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIArticleServiceRecorder) FindByStatus(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FindByStatus", arg0)
}

func (_m *MockIArticleService) Add(article models.Article) (models.Article, error) {
	ret := _m.ctrl.Call(_m, "Add", article)
	ret0, _ := ret[0].(models.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIArticleServiceRecorder) Add(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Add", arg0)
}

func (_m *MockIArticleService) FindAll() ([]models.Article, error) {
	ret := _m.ctrl.Call(_m, "FindAll")
	ret0, _ := ret[0].([]models.Article)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockIArticleServiceRecorder) FindAll() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "FindAll")
}

func (_m *MockIArticleService) Delete(article models.Article) error {
	ret := _m.ctrl.Call(_m, "Delete", article)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIArticleServiceRecorder) Delete(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Delete", arg0)
}

func (_m *MockIArticleService) Update(id string, article models.Article) error {
	ret := _m.ctrl.Call(_m, "Update", id, article)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockIArticleServiceRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Update", arg0, arg1)
}
