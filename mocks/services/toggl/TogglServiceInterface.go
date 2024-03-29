// Code generated by mockery v2.26.1. DO NOT EDIT.

package mocks

import (
	dto "github.com/sharovik/toggl-jira/src/dto"
	arguments "github.com/sharovik/toggl-jira/src/services/arguments"

	mock "github.com/stretchr/testify/mock"
)

// TogglServiceInterface is an autogenerated mock type for the TogglServiceInterface type
type TogglServiceInterface struct {
	mock.Mock
}

// GetReport provides a mock function with given fields: args
func (_m *TogglServiceInterface) GetReport(args arguments.OutputArgs) (dto.TogglDetailsResponse, error) {
	ret := _m.Called(args)

	var r0 dto.TogglDetailsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(arguments.OutputArgs) (dto.TogglDetailsResponse, error)); ok {
		return rf(args)
	}
	if rf, ok := ret.Get(0).(func(arguments.OutputArgs) dto.TogglDetailsResponse); ok {
		r0 = rf(args)
	} else {
		r0 = ret.Get(0).(dto.TogglDetailsResponse)
	}

	if rf, ok := ret.Get(1).(func(arguments.OutputArgs) error); ok {
		r1 = rf(args)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewTogglServiceInterface interface {
	mock.TestingT
	Cleanup(func())
}

// NewTogglServiceInterface creates a new instance of TogglServiceInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewTogglServiceInterface(t mockConstructorTestingTNewTogglServiceInterface) *TogglServiceInterface {
	mock := &TogglServiceInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
