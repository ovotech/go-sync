// Code generated by mockery v2.14.1. DO NOT EDIT.

package user

import (
	context "context"

	tfe "github.com/hashicorp/go-tfe"
	mock "github.com/stretchr/testify/mock"
)

// mockITeamMembers is an autogenerated mock type for the iTeamMembers type
type mockITeamMembers struct {
	mock.Mock
}

type mockITeamMembers_Expecter struct {
	mock *mock.Mock
}

func (_m *mockITeamMembers) EXPECT() *mockITeamMembers_Expecter {
	return &mockITeamMembers_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, teamID, options
func (_m *mockITeamMembers) Add(ctx context.Context, teamID string, options tfe.TeamMemberAddOptions) error {
	ret := _m.Called(ctx, teamID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.TeamMemberAddOptions) error); ok {
		r0 = rf(ctx, teamID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockITeamMembers_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type mockITeamMembers_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//  - ctx context.Context
//  - teamID string
//  - options tfe.TeamMemberAddOptions
func (_e *mockITeamMembers_Expecter) Add(ctx interface{}, teamID interface{}, options interface{}) *mockITeamMembers_Add_Call {
	return &mockITeamMembers_Add_Call{Call: _e.mock.On("Add", ctx, teamID, options)}
}

func (_c *mockITeamMembers_Add_Call) Run(run func(ctx context.Context, teamID string, options tfe.TeamMemberAddOptions)) *mockITeamMembers_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(tfe.TeamMemberAddOptions))
	})
	return _c
}

func (_c *mockITeamMembers_Add_Call) Return(_a0 error) *mockITeamMembers_Add_Call {
	_c.Call.Return(_a0)
	return _c
}

// List provides a mock function with given fields: ctx, teamID
func (_m *mockITeamMembers) List(ctx context.Context, teamID string) ([]*tfe.User, error) {
	ret := _m.Called(ctx, teamID)

	var r0 []*tfe.User
	if rf, ok := ret.Get(0).(func(context.Context, string) []*tfe.User); ok {
		r0 = rf(ctx, teamID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*tfe.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, teamID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockITeamMembers_List_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'List'
type mockITeamMembers_List_Call struct {
	*mock.Call
}

// List is a helper method to define mock.On call
//  - ctx context.Context
//  - teamID string
func (_e *mockITeamMembers_Expecter) List(ctx interface{}, teamID interface{}) *mockITeamMembers_List_Call {
	return &mockITeamMembers_List_Call{Call: _e.mock.On("List", ctx, teamID)}
}

func (_c *mockITeamMembers_List_Call) Run(run func(ctx context.Context, teamID string)) *mockITeamMembers_List_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockITeamMembers_List_Call) Return(_a0 []*tfe.User, _a1 error) *mockITeamMembers_List_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Remove provides a mock function with given fields: ctx, teamID, options
func (_m *mockITeamMembers) Remove(ctx context.Context, teamID string, options tfe.TeamMemberRemoveOptions) error {
	ret := _m.Called(ctx, teamID, options)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, string, tfe.TeamMemberRemoveOptions) error); ok {
		r0 = rf(ctx, teamID, options)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// mockITeamMembers_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type mockITeamMembers_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//  - ctx context.Context
//  - teamID string
//  - options tfe.TeamMemberRemoveOptions
func (_e *mockITeamMembers_Expecter) Remove(ctx interface{}, teamID interface{}, options interface{}) *mockITeamMembers_Remove_Call {
	return &mockITeamMembers_Remove_Call{Call: _e.mock.On("Remove", ctx, teamID, options)}
}

func (_c *mockITeamMembers_Remove_Call) Run(run func(ctx context.Context, teamID string, options tfe.TeamMemberRemoveOptions)) *mockITeamMembers_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(tfe.TeamMemberRemoveOptions))
	})
	return _c
}

func (_c *mockITeamMembers_Remove_Call) Return(_a0 error) *mockITeamMembers_Remove_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTnewMockITeamMembers interface {
	mock.TestingT
	Cleanup(func())
}

// newMockITeamMembers creates a new instance of mockITeamMembers. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func newMockITeamMembers(t mockConstructorTestingTnewMockITeamMembers) *mockITeamMembers {
	mock := &mockITeamMembers{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}