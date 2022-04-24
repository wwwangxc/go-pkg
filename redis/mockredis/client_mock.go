// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package mockredis is a generated GoMock package.
package mockredis

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	redis "github.com/gomodule/redigo/redis"
	redis0 "github.com/wwwangxc/go-pkg/redis"
)

// MockClientProxy is a mock of ClientProxy interface.
type MockClientProxy struct {
	ctrl     *gomock.Controller
	recorder *MockClientProxyMockRecorder
}

// MockClientProxyMockRecorder is the mock recorder for MockClientProxy.
type MockClientProxyMockRecorder struct {
	mock *MockClientProxy
}

// NewMockClientProxy creates a new mock instance.
func NewMockClientProxy(ctrl *gomock.Controller) *MockClientProxy {
	mock := &MockClientProxy{ctrl: ctrl}
	mock.recorder = &MockClientProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClientProxy) EXPECT() *MockClientProxyMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockClientProxy) Do(ctx context.Context, cmd string, args ...interface{}) (interface{}, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, cmd}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Do", varargs...)
	ret0, _ := ret[0].(interface{})
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do.
func (mr *MockClientProxyMockRecorder) Do(ctx, cmd interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, cmd}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockClientProxy)(nil).Do), varargs...)
}

// GetConn mocks base method.
func (m *MockClientProxy) GetConn() redis.Conn {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConn")
	ret0, _ := ret[0].(redis.Conn)
	return ret0
}

// GetConn indicates an expected call of GetConn.
func (mr *MockClientProxyMockRecorder) GetConn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConn", reflect.TypeOf((*MockClientProxy)(nil).GetConn))
}

// GetFetcher mocks base method.
func (m *MockClientProxy) GetFetcher() redis0.Fetcher {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFetcher")
	ret0, _ := ret[0].(redis0.Fetcher)
	return ret0
}

// GetFetcher indicates an expected call of GetFetcher.
func (mr *MockClientProxyMockRecorder) GetFetcher() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFetcher", reflect.TypeOf((*MockClientProxy)(nil).GetFetcher))
}

// GetLocker mocks base method.
func (m *MockClientProxy) GetLocker() redis0.Locker {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLocker")
	ret0, _ := ret[0].(redis0.Locker)
	return ret0
}

// GetLocker indicates an expected call of GetLocker.
func (mr *MockClientProxyMockRecorder) GetLocker() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLocker", reflect.TypeOf((*MockClientProxy)(nil).GetLocker))
}
