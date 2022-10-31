// Code generated by mockery v2.14.1. DO NOT EDIT.

package discovery

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockGitHubDiscovery is an autogenerated mock type for the GitHubDiscovery type
type MockGitHubDiscovery struct {
	mock.Mock
}

type MockGitHubDiscovery_Expecter struct {
	mock *mock.Mock
}

func (_m *MockGitHubDiscovery) EXPECT() *MockGitHubDiscovery_Expecter {
	return &MockGitHubDiscovery_Expecter{mock: &_m.Mock}
}

// GetEmailFromUsername provides a mock function with given fields: _a0, _a1
func (_m *MockGitHubDiscovery) GetEmailFromUsername(_a0 context.Context, _a1 []string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubDiscovery_GetEmailFromUsername_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetEmailFromUsername'
type MockGitHubDiscovery_GetEmailFromUsername_Call struct {
	*mock.Call
}

// GetEmailFromUsername is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 []string
func (_e *MockGitHubDiscovery_Expecter) GetEmailFromUsername(_a0 interface{}, _a1 interface{}) *MockGitHubDiscovery_GetEmailFromUsername_Call {
	return &MockGitHubDiscovery_GetEmailFromUsername_Call{Call: _e.mock.On("GetEmailFromUsername", _a0, _a1)}
}

func (_c *MockGitHubDiscovery_GetEmailFromUsername_Call) Run(run func(_a0 context.Context, _a1 []string)) *MockGitHubDiscovery_GetEmailFromUsername_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockGitHubDiscovery_GetEmailFromUsername_Call) Return(_a0 []string, _a1 error) *MockGitHubDiscovery_GetEmailFromUsername_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetUsernameFromEmail provides a mock function with given fields: _a0, _a1
func (_m *MockGitHubDiscovery) GetUsernameFromEmail(_a0 context.Context, _a1 []string) ([]string, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, []string) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, []string) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockGitHubDiscovery_GetUsernameFromEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUsernameFromEmail'
type MockGitHubDiscovery_GetUsernameFromEmail_Call struct {
	*mock.Call
}

// GetUsernameFromEmail is a helper method to define mock.On call
//  - _a0 context.Context
//  - _a1 []string
func (_e *MockGitHubDiscovery_Expecter) GetUsernameFromEmail(_a0 interface{}, _a1 interface{}) *MockGitHubDiscovery_GetUsernameFromEmail_Call {
	return &MockGitHubDiscovery_GetUsernameFromEmail_Call{Call: _e.mock.On("GetUsernameFromEmail", _a0, _a1)}
}

func (_c *MockGitHubDiscovery_GetUsernameFromEmail_Call) Run(run func(_a0 context.Context, _a1 []string)) *MockGitHubDiscovery_GetUsernameFromEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]string))
	})
	return _c
}

func (_c *MockGitHubDiscovery_GetUsernameFromEmail_Call) Return(_a0 []string, _a1 error) *MockGitHubDiscovery_GetUsernameFromEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewMockGitHubDiscovery interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockGitHubDiscovery creates a new instance of MockGitHubDiscovery. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockGitHubDiscovery(t mockConstructorTestingTNewMockGitHubDiscovery) *MockGitHubDiscovery {
	mock := &MockGitHubDiscovery{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
