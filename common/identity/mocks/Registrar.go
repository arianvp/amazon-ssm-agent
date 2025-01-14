// Code generated by mockery v2.23.2. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Registrar is an autogenerated mock type for the Registrar type
type Registrar struct {
	mock.Mock
}

// Register provides a mock function with given fields: _a0
func (_m *Registrar) Register(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

type mockConstructorTestingTNewRegistrar interface {
	mock.TestingT
	Cleanup(func())
}

// NewRegistrar creates a new instance of Registrar. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRegistrar(t mockConstructorTestingTNewRegistrar) *Registrar {
	mock := &Registrar{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
