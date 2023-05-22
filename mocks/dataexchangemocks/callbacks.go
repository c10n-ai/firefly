// Code generated by mockery v2.20.2. DO NOT EDIT.

package dataexchangemocks

import (
	dataexchange "github.com/hyperledger/firefly/pkg/dataexchange"
	mock "github.com/stretchr/testify/mock"
)

// Callbacks is an autogenerated mock type for the Callbacks type
type Callbacks struct {
	mock.Mock
}

// DXEvent provides a mock function with given fields: plugin, event
func (_m *Callbacks) DXEvent(plugin dataexchange.Plugin, event dataexchange.DXEvent) error {
	ret := _m.Called(plugin, event)

	var r0 error
	if rf, ok := ret.Get(0).(func(dataexchange.Plugin, dataexchange.DXEvent) error); ok {
		r0 = rf(plugin, event)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewCallbacks interface {
	mock.TestingT
	Cleanup(func())
}

// NewCallbacks creates a new instance of Callbacks. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCallbacks(t mockConstructorTestingTNewCallbacks) *Callbacks {
	mock := &Callbacks{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
