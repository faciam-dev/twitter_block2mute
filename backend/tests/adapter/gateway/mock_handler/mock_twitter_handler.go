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

// GetBlockedUser mocks base method.
func (m *MockTwitterHandler) GetBlockedUser(arg0 string) (handler.TwitterUserIds, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlockedUser", arg0)
	ret0, _ := ret[0].(handler.TwitterUserIds)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBlockedUser indicates an expected call of GetBlockedUser.
func (mr *MockTwitterHandlerMockRecorder) GetBlockedUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlockedUser", reflect.TypeOf((*MockTwitterHandler)(nil).GetBlockedUser), arg0)
}

// GetCredentials mocks base method.
func (m *MockTwitterHandler) GetCredentials(arg0, arg1 string) (handler.TwitterCredentials, handler.TwitterValues, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCredentials", arg0, arg1)
	ret0, _ := ret[0].(handler.TwitterCredentials)
	ret1, _ := ret[1].(handler.TwitterValues)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// GetCredentials indicates an expected call of GetCredentials.
func (mr *MockTwitterHandlerMockRecorder) GetCredentials(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCredentials", reflect.TypeOf((*MockTwitterHandler)(nil).GetCredentials), arg0, arg1)
}

// GetUser mocks base method.
func (m *MockTwitterHandler) GetUser(arg0 string) (handler.TwitterUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0)
	ret0, _ := ret[0].(handler.TwitterUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockTwitterHandlerMockRecorder) GetUser(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockTwitterHandler)(nil).GetUser), arg0)
}

// UpdateTwitterApi mocks base method.
func (m *MockTwitterHandler) UpdateTwitterApi(arg0, arg1 string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "UpdateTwitterApi", arg0, arg1)
}

// UpdateTwitterApi indicates an expected call of UpdateTwitterApi.
func (mr *MockTwitterHandlerMockRecorder) UpdateTwitterApi(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTwitterApi", reflect.TypeOf((*MockTwitterHandler)(nil).UpdateTwitterApi), arg0, arg1)
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

// MockTwitterValues is a mock of TwitterValues interface.
type MockTwitterValues struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterValuesMockRecorder
}

// MockTwitterValuesMockRecorder is the mock recorder for MockTwitterValues.
type MockTwitterValuesMockRecorder struct {
	mock *MockTwitterValues
}

// NewMockTwitterValues creates a new mock instance.
func NewMockTwitterValues(ctrl *gomock.Controller) *MockTwitterValues {
	mock := &MockTwitterValues{ctrl: ctrl}
	mock.recorder = &MockTwitterValuesMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterValues) EXPECT() *MockTwitterValuesMockRecorder {
	return m.recorder
}

// GetTwitterID mocks base method.
func (m *MockTwitterValues) GetTwitterID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTwitterID indicates an expected call of GetTwitterID.
func (mr *MockTwitterValuesMockRecorder) GetTwitterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterID", reflect.TypeOf((*MockTwitterValues)(nil).GetTwitterID))
}

// GetTwitterScreenName mocks base method.
func (m *MockTwitterValues) GetTwitterScreenName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterScreenName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTwitterScreenName indicates an expected call of GetTwitterScreenName.
func (mr *MockTwitterValuesMockRecorder) GetTwitterScreenName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterScreenName", reflect.TypeOf((*MockTwitterValues)(nil).GetTwitterScreenName))
}

// MockTwitterUser is a mock of TwitterUser interface.
type MockTwitterUser struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterUserMockRecorder
}

// MockTwitterUserMockRecorder is the mock recorder for MockTwitterUser.
type MockTwitterUserMockRecorder struct {
	mock *MockTwitterUser
}

// NewMockTwitterUser creates a new mock instance.
func NewMockTwitterUser(ctrl *gomock.Controller) *MockTwitterUser {
	mock := &MockTwitterUser{ctrl: ctrl}
	mock.recorder = &MockTwitterUserMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterUser) EXPECT() *MockTwitterUserMockRecorder {
	return m.recorder
}

// GetTwitterID mocks base method.
func (m *MockTwitterUser) GetTwitterID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTwitterID indicates an expected call of GetTwitterID.
func (mr *MockTwitterUserMockRecorder) GetTwitterID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterID", reflect.TypeOf((*MockTwitterUser)(nil).GetTwitterID))
}

// GetTwitterName mocks base method.
func (m *MockTwitterUser) GetTwitterName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTwitterName indicates an expected call of GetTwitterName.
func (mr *MockTwitterUserMockRecorder) GetTwitterName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterName", reflect.TypeOf((*MockTwitterUser)(nil).GetTwitterName))
}

// GetTwitterScreenName mocks base method.
func (m *MockTwitterUser) GetTwitterScreenName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterScreenName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetTwitterScreenName indicates an expected call of GetTwitterScreenName.
func (mr *MockTwitterUserMockRecorder) GetTwitterScreenName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterScreenName", reflect.TypeOf((*MockTwitterUser)(nil).GetTwitterScreenName))
}

// MockTwitterUserIds is a mock of TwitterUserIds interface.
type MockTwitterUserIds struct {
	ctrl     *gomock.Controller
	recorder *MockTwitterUserIdsMockRecorder
}

// MockTwitterUserIdsMockRecorder is the mock recorder for MockTwitterUserIds.
type MockTwitterUserIdsMockRecorder struct {
	mock *MockTwitterUserIds
}

// NewMockTwitterUserIds creates a new mock instance.
func NewMockTwitterUserIds(ctrl *gomock.Controller) *MockTwitterUserIds {
	mock := &MockTwitterUserIds{ctrl: ctrl}
	mock.recorder = &MockTwitterUserIdsMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTwitterUserIds) EXPECT() *MockTwitterUserIdsMockRecorder {
	return m.recorder
}

// GetTotal mocks base method.
func (m *MockTwitterUserIds) GetTotal() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTotal")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetTotal indicates an expected call of GetTotal.
func (mr *MockTwitterUserIdsMockRecorder) GetTotal() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTotal", reflect.TypeOf((*MockTwitterUserIds)(nil).GetTotal))
}

// GetTwitterIDs mocks base method.
func (m *MockTwitterUserIds) GetTwitterIDs() []string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTwitterIDs")
	ret0, _ := ret[0].([]string)
	return ret0
}

// GetTwitterIDs indicates an expected call of GetTwitterIDs.
func (mr *MockTwitterUserIdsMockRecorder) GetTwitterIDs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTwitterIDs", reflect.TypeOf((*MockTwitterUserIds)(nil).GetTwitterIDs))
}
