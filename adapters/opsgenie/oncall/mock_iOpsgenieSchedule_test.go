// Code generated by mockery. DO NOT EDIT.

package oncall

import (
	context "context"

	schedule "github.com/opsgenie/opsgenie-go-sdk-v2/schedule"
	mock "github.com/stretchr/testify/mock"
)

// mockIOpsgenieSchedule is an autogenerated mock type for the iOpsgenieSchedule type
type mockIOpsgenieSchedule struct {
	mock.Mock
}

type mockIOpsgenieSchedule_Expecter struct {
	mock *mock.Mock
}

func (_m *mockIOpsgenieSchedule) EXPECT() *mockIOpsgenieSchedule_Expecter {
	return &mockIOpsgenieSchedule_Expecter{mock: &_m.Mock}
}

// GetOnCalls provides a mock function with given fields: _a0, request
func (_m *mockIOpsgenieSchedule) GetOnCalls(_a0 context.Context, request *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error) {
	ret := _m.Called(_a0, request)

	var r0 *schedule.GetOnCallsResult
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)); ok {
		return rf(_a0, request)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *schedule.GetOnCallsRequest) *schedule.GetOnCallsResult); ok {
		r0 = rf(_a0, request)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*schedule.GetOnCallsResult)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *schedule.GetOnCallsRequest) error); ok {
		r1 = rf(_a0, request)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockIOpsgenieSchedule_GetOnCalls_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetOnCalls'
type mockIOpsgenieSchedule_GetOnCalls_Call struct {
	*mock.Call
}

// GetOnCalls is a helper method to define mock.On call
//   - _a0 context.Context
//   - request *schedule.GetOnCallsRequest
func (_e *mockIOpsgenieSchedule_Expecter) GetOnCalls(_a0 interface{}, request interface{}) *mockIOpsgenieSchedule_GetOnCalls_Call {
	return &mockIOpsgenieSchedule_GetOnCalls_Call{Call: _e.mock.On("GetOnCalls", _a0, request)}
}

func (_c *mockIOpsgenieSchedule_GetOnCalls_Call) Run(run func(_a0 context.Context, request *schedule.GetOnCallsRequest)) *mockIOpsgenieSchedule_GetOnCalls_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*schedule.GetOnCallsRequest))
	})
	return _c
}

func (_c *mockIOpsgenieSchedule_GetOnCalls_Call) Return(_a0 *schedule.GetOnCallsResult, _a1 error) *mockIOpsgenieSchedule_GetOnCalls_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockIOpsgenieSchedule_GetOnCalls_Call) RunAndReturn(run func(context.Context, *schedule.GetOnCallsRequest) (*schedule.GetOnCallsResult, error)) *mockIOpsgenieSchedule_GetOnCalls_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTnewMockIOpsgenieSchedule interface {
	mock.TestingT
	Cleanup(func())
}

// newMockIOpsgenieSchedule creates a new instance of mockIOpsgenieSchedule. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockIOpsgenieSchedule(t mockConstructorTestingTnewMockIOpsgenieSchedule) *mockIOpsgenieSchedule {
	mock := &mockIOpsgenieSchedule{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
