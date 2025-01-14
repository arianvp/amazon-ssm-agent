// Code generated by mockery v2.23.2. DO NOT EDIT.

package mocks

import (
	context "context"

	request "github.com/aws/aws-sdk-go/aws/request"
	mock "github.com/stretchr/testify/mock"

	ssm "github.com/aws/aws-sdk-go/service/ssm"
)

// ISSMClient is an autogenerated mock type for the ISSMClient type
type ISSMClient struct {
	mock.Mock
}

// UpdateInstanceInformation provides a mock function with given fields: input
func (_m *ISSMClient) UpdateInstanceInformation(input *ssm.UpdateInstanceInformationInput) (*ssm.UpdateInstanceInformationOutput, error) {
	ret := _m.Called(input)

	var r0 *ssm.UpdateInstanceInformationOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(*ssm.UpdateInstanceInformationInput) (*ssm.UpdateInstanceInformationOutput, error)); ok {
		return rf(input)
	}
	if rf, ok := ret.Get(0).(func(*ssm.UpdateInstanceInformationInput) *ssm.UpdateInstanceInformationOutput); ok {
		r0 = rf(input)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssm.UpdateInstanceInformationOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(*ssm.UpdateInstanceInformationInput) error); ok {
		r1 = rf(input)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateInstanceInformationWithContext provides a mock function with given fields: ctx, input, opts
func (_m *ISSMClient) UpdateInstanceInformationWithContext(ctx context.Context, input *ssm.UpdateInstanceInformationInput, opts ...request.Option) (*ssm.UpdateInstanceInformationOutput, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, input)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *ssm.UpdateInstanceInformationOutput
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.UpdateInstanceInformationInput, ...request.Option) (*ssm.UpdateInstanceInformationOutput, error)); ok {
		return rf(ctx, input, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *ssm.UpdateInstanceInformationInput, ...request.Option) *ssm.UpdateInstanceInformationOutput); ok {
		r0 = rf(ctx, input, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*ssm.UpdateInstanceInformationOutput)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *ssm.UpdateInstanceInformationInput, ...request.Option) error); ok {
		r1 = rf(ctx, input, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewISSMClient interface {
	mock.TestingT
	Cleanup(func())
}

// NewISSMClient creates a new instance of ISSMClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISSMClient(t mockConstructorTestingTNewISSMClient) *ISSMClient {
	mock := &ISSMClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
