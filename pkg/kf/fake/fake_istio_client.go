// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GoogleCloudPlatform/kf/pkg/kf/fake (interfaces: IstioClient)

// Package fake is a generated GoMock package.
package fake

import (
	kf "github.com/GoogleCloudPlatform/kf/pkg/kf"
	gomock "github.com/golang/mock/gomock"
	v1 "k8s.io/api/core/v1"
	reflect "reflect"
)

// FakeIstioClient is a mock of IstioClient interface
type FakeIstioClient struct {
	ctrl     *gomock.Controller
	recorder *FakeIstioClientMockRecorder
}

// FakeIstioClientMockRecorder is the mock recorder for FakeIstioClient
type FakeIstioClientMockRecorder struct {
	mock *FakeIstioClient
}

// NewFakeIstioClient creates a new mock instance
func NewFakeIstioClient(ctrl *gomock.Controller) *FakeIstioClient {
	mock := &FakeIstioClient{ctrl: ctrl}
	mock.recorder = &FakeIstioClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *FakeIstioClient) EXPECT() *FakeIstioClientMockRecorder {
	return m.recorder
}

// ListIngresses mocks base method
func (m *FakeIstioClient) ListIngresses(arg0 ...kf.ListIngressesOption) ([]v1.LoadBalancerIngress, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range arg0 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ListIngresses", varargs...)
	ret0, _ := ret[0].([]v1.LoadBalancerIngress)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListIngresses indicates an expected call of ListIngresses
func (mr *FakeIstioClientMockRecorder) ListIngresses(arg0 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListIngresses", reflect.TypeOf((*FakeIstioClient)(nil).ListIngresses), arg0...)
}
