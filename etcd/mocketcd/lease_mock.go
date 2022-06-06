// Code generated by MockGen. DO NOT EDIT.
// Source: lease.go

// Package mocketcd is a generated GoMock package.
package mocketcd

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	clientv3 "go.etcd.io/etcd/client/v3"
)

// MockLeaseProxy is a mock of LeaseProxy interface.
type MockLeaseProxy struct {
	ctrl     *gomock.Controller
	recorder *MockLeaseProxyMockRecorder
}

// MockLeaseProxyMockRecorder is the mock recorder for MockLeaseProxy.
type MockLeaseProxyMockRecorder struct {
	mock *MockLeaseProxy
}

// NewMockLeaseProxy creates a new mock instance.
func NewMockLeaseProxy(ctrl *gomock.Controller) *MockLeaseProxy {
	mock := &MockLeaseProxy{ctrl: ctrl}
	mock.recorder = &MockLeaseProxyMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLeaseProxy) EXPECT() *MockLeaseProxyMockRecorder {
	return m.recorder
}

// Grant mocks base method.
func (m *MockLeaseProxy) Grant(ctx context.Context, ttl int64) (*clientv3.LeaseGrantResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Grant", ctx, ttl)
	ret0, _ := ret[0].(*clientv3.LeaseGrantResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Grant indicates an expected call of Grant.
func (mr *MockLeaseProxyMockRecorder) Grant(ctx, ttl interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Grant", reflect.TypeOf((*MockLeaseProxy)(nil).Grant), ctx, ttl)
}

// KeepAlive mocks base method.
func (m *MockLeaseProxy) KeepAlive(ctx context.Context, id clientv3.LeaseID) (<-chan *clientv3.LeaseKeepAliveResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "KeepAlive", ctx, id)
	ret0, _ := ret[0].(<-chan *clientv3.LeaseKeepAliveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// KeepAlive indicates an expected call of KeepAlive.
func (mr *MockLeaseProxyMockRecorder) KeepAlive(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "KeepAlive", reflect.TypeOf((*MockLeaseProxy)(nil).KeepAlive), ctx, id)
}

// Leases mocks base method.
func (m *MockLeaseProxy) Leases(ctx context.Context) (*clientv3.LeaseLeasesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Leases", ctx)
	ret0, _ := ret[0].(*clientv3.LeaseLeasesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Leases indicates an expected call of Leases.
func (mr *MockLeaseProxyMockRecorder) Leases(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Leases", reflect.TypeOf((*MockLeaseProxy)(nil).Leases), ctx)
}

// Revoke mocks base method.
func (m *MockLeaseProxy) Revoke(ctx context.Context, id clientv3.LeaseID) (*clientv3.LeaseRevokeResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Revoke", ctx, id)
	ret0, _ := ret[0].(*clientv3.LeaseRevokeResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Revoke indicates an expected call of Revoke.
func (mr *MockLeaseProxyMockRecorder) Revoke(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Revoke", reflect.TypeOf((*MockLeaseProxy)(nil).Revoke), ctx, id)
}

// TimeToLive mocks base method.
func (m *MockLeaseProxy) TimeToLive(ctx context.Context, id clientv3.LeaseID, opts ...clientv3.LeaseOption) (*clientv3.LeaseTimeToLiveResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, id}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "TimeToLive", varargs...)
	ret0, _ := ret[0].(*clientv3.LeaseTimeToLiveResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TimeToLive indicates an expected call of TimeToLive.
func (mr *MockLeaseProxyMockRecorder) TimeToLive(ctx, id interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, id}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TimeToLive", reflect.TypeOf((*MockLeaseProxy)(nil).TimeToLive), varargs...)
}
