// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/keyjin88/shortener/internal/app/handlers (interfaces: RequestContext)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRequestContext is a mock of RequestContext interface.
type MockRequestContext struct {
	ctrl     *gomock.Controller
	recorder *MockRequestContextMockRecorder
}

// MockRequestContextMockRecorder is the mock recorder for MockRequestContext.
type MockRequestContextMockRecorder struct {
	mock *MockRequestContext
}

// NewMockRequestContext creates a new mock instance.
func NewMockRequestContext(ctrl *gomock.Controller) *MockRequestContext {
	mock := &MockRequestContext{ctrl: ctrl}
	mock.recorder = &MockRequestContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRequestContext) EXPECT() *MockRequestContextMockRecorder {
	return m.recorder
}

// AbortWithStatus mocks base method.
func (m *MockRequestContext) AbortWithStatus(arg0 int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AbortWithStatus", arg0)
}

// AbortWithStatus indicates an expected call of AbortWithStatus.
func (mr *MockRequestContextMockRecorder) AbortWithStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AbortWithStatus", reflect.TypeOf((*MockRequestContext)(nil).AbortWithStatus), arg0)
}

// BindJSON mocks base method.
func (m *MockRequestContext) BindJSON(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BindJSON", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// BindJSON indicates an expected call of BindJSON.
func (mr *MockRequestContextMockRecorder) BindJSON(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BindJSON", reflect.TypeOf((*MockRequestContext)(nil).BindJSON), arg0)
}

// FullPath mocks base method.
func (m *MockRequestContext) FullPath() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FullPath")
	ret0, _ := ret[0].(string)
	return ret0
}

// FullPath indicates an expected call of FullPath.
func (mr *MockRequestContextMockRecorder) FullPath() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FullPath", reflect.TypeOf((*MockRequestContext)(nil).FullPath))
}

// GetRawData mocks base method.
func (m *MockRequestContext) GetRawData() ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRawData")
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRawData indicates an expected call of GetRawData.
func (mr *MockRequestContextMockRecorder) GetRawData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRawData", reflect.TypeOf((*MockRequestContext)(nil).GetRawData))
}

// GetString mocks base method.
func (m *MockRequestContext) GetString(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetString", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// GetString indicates an expected call of GetString.
func (mr *MockRequestContextMockRecorder) GetString(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetString", reflect.TypeOf((*MockRequestContext)(nil).GetString), arg0)
}

// Header mocks base method.
func (m *MockRequestContext) Header(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Header", arg0, arg1)
}

// Header indicates an expected call of Header.
func (mr *MockRequestContextMockRecorder) Header(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockRequestContext)(nil).Header), arg0, arg1)
}

// JSON mocks base method.
func (m *MockRequestContext) JSON(arg0 int, arg1 interface{}) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "JSON", arg0, arg1)
}

// JSON indicates an expected call of JSON.
func (mr *MockRequestContextMockRecorder) JSON(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JSON", reflect.TypeOf((*MockRequestContext)(nil).JSON), arg0, arg1)
}

// Param mocks base method.
func (m *MockRequestContext) Param(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Param", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Param indicates an expected call of Param.
func (mr *MockRequestContextMockRecorder) Param(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Param", reflect.TypeOf((*MockRequestContext)(nil).Param), arg0)
}

// Redirect mocks base method.
func (m *MockRequestContext) Redirect(arg0 int, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Redirect", arg0, arg1)
}

// Redirect indicates an expected call of Redirect.
func (mr *MockRequestContextMockRecorder) Redirect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Redirect", reflect.TypeOf((*MockRequestContext)(nil).Redirect), arg0, arg1)
}

// ShouldBind mocks base method.
func (m *MockRequestContext) ShouldBind(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ShouldBind", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// ShouldBind indicates an expected call of ShouldBind.
func (mr *MockRequestContextMockRecorder) ShouldBind(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ShouldBind", reflect.TypeOf((*MockRequestContext)(nil).ShouldBind), arg0)
}

// String mocks base method.
func (m *MockRequestContext) String(arg0 int, arg1 string, arg2 ...interface{}) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	m.ctrl.Call(m, "String", varargs...)
}

// String indicates an expected call of String.
func (mr *MockRequestContextMockRecorder) String(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockRequestContext)(nil).String), varargs...)
}
