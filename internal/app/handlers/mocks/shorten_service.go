// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/keyjin88/shortener/internal/app/handlers (interfaces: ShortenService)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockShortenService is a mock of ShortenService interface.
type MockShortenService struct {
	ctrl     *gomock.Controller
	recorder *MockShortenServiceMockRecorder
}

// MockShortenServiceMockRecorder is the mock recorder for MockShortenService.
type MockShortenServiceMockRecorder struct {
	mock *MockShortenService
}

// NewMockShortenService creates a new mock instance.
func NewMockShortenService(ctrl *gomock.Controller) *MockShortenService {
	mock := &MockShortenService{ctrl: ctrl}
	mock.recorder = &MockShortenServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortenService) EXPECT() *MockShortenServiceMockRecorder {
	return m.recorder
}

// GetShortenedURLByID mocks base method.
func (m *MockShortenService) GetShortenedURLByID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortenedURLByID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortenedURLByID indicates an expected call of GetShortenedURLByID.
func (mr *MockShortenServiceMockRecorder) GetShortenedURLByID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortenedURLByID", reflect.TypeOf((*MockShortenService)(nil).GetShortenedURLByID), arg0)
}

// ShortenURL mocks base method.
func (m *MockShortenService) ShortenURL(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShortenURL", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ShortenURL indicates an expected call of ShortenURL.
func (mr *MockShortenServiceMockRecorder) ShortenURL(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShortenURL", reflect.TypeOf((*MockShortenService)(nil).ShortenURL), arg0)
}