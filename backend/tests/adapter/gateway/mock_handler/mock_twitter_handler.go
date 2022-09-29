// Code generated by MockGen. DO NOT EDIT.
// Source: ./adapter/gateway/handler/twitter_handler.go

// Package mock_handler is a generated GoMock package.
package mock_handler

import (
	reflect "reflect"

	handler "github.com/faciam_dev/twitter_block2mute/backend/adapter/gateway/handler"
	gomock "github.com/golang/mock/gomock"
)

// MockTwitterHandler is a mock of TwitterHandler interface.
type MockTwitterHandler struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterHandlerMockRecorder
}

// MockTwitterHandlerMockRecorder is the mock recorder for MockTwitterHandler.
type MockTwitterHandlerMockRecorder struct {
	mock *MockTwitterHandler
}

// NewMockTwitterHandler creates a new mock instance.
func NewMockTwitterHandler(ctrl *gomock.Controller) *MockTwitterHandler {
	mock := &MockTwitterHandler{ctrl: ctrl}
	mock.recorder = &MockTwitterHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterHandler) EXPECT() *MockTwitterHandlerMockRecorder {
	return m.recorder
}

// AuthorizationURL mocks base method.
func (m *MockTwitterHandler) AuthorizationURL() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AuthorizationURL")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AuthorizationURL indicates an expected call of AuthorizationURL.
func (mr *MockTwitterHandlerMockRecorder) AuthorizationURL() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AuthorizationURL", reflect.TypeOf((*MockTwitterHandler)(nil).AuthorizationURL))
}

// GetCredentials mocks base method.
func (m *MockTwitterHandler) GetCredentials(arg0, arg1 string) (handler.TwitterCredentials, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", arg0, arg1)
	ret0, _ := ret[0].(handler.TwitterCredentials)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCredentials indicates an expected call of GetCredentials.
func (mr *MockTwitterHandlerMockRecorder) GetCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockTwitterHandler)(nil).GetCredentials), arg0, arg1)
}

// GetRateLimits mocks base method.
func (m *MockTwitterHandler) GetRateLimits() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRateLimits")
	ret0, _ := ret[0].(error)
	return ret0
}

// GetRateLimits indicates an expected call of GetRateLimits.
func (mr *MockTwitterHandlerMockRecorder) GetRateLimits() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRateLimits", reflect.TypeOf((*MockTwitterHandler)(nil).GetRateLimits))
}

// SetCredentials mocks base method.
func (m *MockTwitterHandler) SetCredentials(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetCredentials", arg0, arg1)
}

// SetCredentials indicates an expected call of SetCredentials.
func (mr *MockTwitterHandlerMockRecorder) SetCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetCredentials", reflect.TypeOf((*MockTwitterHandler)(nil).SetCredentials), arg0, arg1)
}

// MockTwitterCredentials is a mock of TwitterCredentials interface.
type MockTwitterCredentials struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterCredentialsMockRecorder
}

// MockTwitterCredentialsMockRecorder is the mock recorder for MockTwitterCredentials.
type MockTwitterCredentialsMockRecorder struct {
	mock *MockTwitterCredentials
}

// NewMockTwitterCredentials creates a new mock instance.
func NewMockTwitterCredentials(ctrl *gomock.Controller) *MockTwitterCredentials {
	mock := &MockTwitterCredentials{ctrl: ctrl}
	mock.recorder = &MockTwitterCredentialsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterCredentials) EXPECT() *MockTwitterCredentialsMockRecorder {
	return m.recorder
}

// GetSecret mocks base method.
func (m *MockTwitterCredentials) GetSecret() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSecret")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetSecret indicates an expected call of GetSecret.
func (mr *MockTwitterCredentialsMockRecorder) GetSecret() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSecret", reflect.TypeOf((*MockTwitterCredentials)(nil).GetSecret))
}

// GetToken mocks base method.
func (m *MockTwitterCredentials) GetToken() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetToken")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetToken indicates an expected call of GetToken.
func (mr *MockTwitterCredentialsMockRecorder) GetToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetToken", reflect.TypeOf((*MockTwitterCredentials)(nil).GetToken))
}