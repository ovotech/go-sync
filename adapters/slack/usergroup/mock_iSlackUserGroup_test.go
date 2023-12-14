// Code generated by mockery. DO NOT EDIT.

package usergroup

import (
	context "context"

	slack "github.com/slack-go/slack"
	mock "github.com/stretchr/testify/mock"
)

// mockISlackUserGroup is an autogenerated mock type for the iSlackUserGroup type
type mockISlackUserGroup struct {
	mock.Mock
}

type mockISlackUserGroup_Expecter struct {
	mock *mock.Mock
}

func (_m *mockISlackUserGroup) EXPECT() *mockISlackUserGroup_Expecter {
	return &mockISlackUserGroup_Expecter{mock: &_m.Mock}
}

// GetUserByEmailContext provides a mock function with given fields: ctx, email
func (_m *mockISlackUserGroup) GetUserByEmailContext(ctx context.Context, email string) (*slack.User, error) {
	ret := _m.Called(ctx, email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmailContext")
	}

	var r0 *slack.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (*slack.User, error)); ok {
		return rf(ctx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *slack.User); ok {
		r0 = rf(ctx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*slack.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockISlackUserGroup_GetUserByEmailContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmailContext'
type mockISlackUserGroup_GetUserByEmailContext_Call struct {
	*mock.Call
}

// GetUserByEmailContext is a helper method to define mock.On call
//   - ctx context.Context
//   - email string
func (_e *mockISlackUserGroup_Expecter) GetUserByEmailContext(ctx interface{}, email interface{}) *mockISlackUserGroup_GetUserByEmailContext_Call {
	return &mockISlackUserGroup_GetUserByEmailContext_Call{Call: _e.mock.On("GetUserByEmailContext", ctx, email)}
}

func (_c *mockISlackUserGroup_GetUserByEmailContext_Call) Run(run func(ctx context.Context, email string)) *mockISlackUserGroup_GetUserByEmailContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockISlackUserGroup_GetUserByEmailContext_Call) Return(_a0 *slack.User, _a1 error) *mockISlackUserGroup_GetUserByEmailContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockISlackUserGroup_GetUserByEmailContext_Call) RunAndReturn(run func(context.Context, string) (*slack.User, error)) *mockISlackUserGroup_GetUserByEmailContext_Call {
	_c.Call.Return(run)
	return _c
}

// GetUserGroupMembersContext provides a mock function with given fields: ctx, userGroup
func (_m *mockISlackUserGroup) GetUserGroupMembersContext(ctx context.Context, userGroup string) ([]string, error) {
	ret := _m.Called(ctx, userGroup)

	if len(ret) == 0 {
		panic("no return value specified for GetUserGroupMembersContext")
	}

	var r0 []string
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) ([]string, error)); ok {
		return rf(ctx, userGroup)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) []string); ok {
		r0 = rf(ctx, userGroup)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, userGroup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockISlackUserGroup_GetUserGroupMembersContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserGroupMembersContext'
type mockISlackUserGroup_GetUserGroupMembersContext_Call struct {
	*mock.Call
}

// GetUserGroupMembersContext is a helper method to define mock.On call
//   - ctx context.Context
//   - userGroup string
func (_e *mockISlackUserGroup_Expecter) GetUserGroupMembersContext(ctx interface{}, userGroup interface{}) *mockISlackUserGroup_GetUserGroupMembersContext_Call {
	return &mockISlackUserGroup_GetUserGroupMembersContext_Call{Call: _e.mock.On("GetUserGroupMembersContext", ctx, userGroup)}
}

func (_c *mockISlackUserGroup_GetUserGroupMembersContext_Call) Run(run func(ctx context.Context, userGroup string)) *mockISlackUserGroup_GetUserGroupMembersContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string))
	})
	return _c
}

func (_c *mockISlackUserGroup_GetUserGroupMembersContext_Call) Return(_a0 []string, _a1 error) *mockISlackUserGroup_GetUserGroupMembersContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockISlackUserGroup_GetUserGroupMembersContext_Call) RunAndReturn(run func(context.Context, string) ([]string, error)) *mockISlackUserGroup_GetUserGroupMembersContext_Call {
	_c.Call.Return(run)
	return _c
}

// GetUsersInfoContext provides a mock function with given fields: ctx, users
func (_m *mockISlackUserGroup) GetUsersInfoContext(ctx context.Context, users ...string) (*[]slack.User, error) {
	_va := make([]interface{}, len(users))
	for _i := range users {
		_va[_i] = users[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for GetUsersInfoContext")
	}

	var r0 *[]slack.User
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, ...string) (*[]slack.User, error)); ok {
		return rf(ctx, users...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, ...string) *[]slack.User); ok {
		r0 = rf(ctx, users...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]slack.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, ...string) error); ok {
		r1 = rf(ctx, users...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockISlackUserGroup_GetUsersInfoContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUsersInfoContext'
type mockISlackUserGroup_GetUsersInfoContext_Call struct {
	*mock.Call
}

// GetUsersInfoContext is a helper method to define mock.On call
//   - ctx context.Context
//   - users ...string
func (_e *mockISlackUserGroup_Expecter) GetUsersInfoContext(ctx interface{}, users ...interface{}) *mockISlackUserGroup_GetUsersInfoContext_Call {
	return &mockISlackUserGroup_GetUsersInfoContext_Call{Call: _e.mock.On("GetUsersInfoContext",
		append([]interface{}{ctx}, users...)...)}
}

func (_c *mockISlackUserGroup_GetUsersInfoContext_Call) Run(run func(ctx context.Context, users ...string)) *mockISlackUserGroup_GetUsersInfoContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(args[0].(context.Context), variadicArgs...)
	})
	return _c
}

func (_c *mockISlackUserGroup_GetUsersInfoContext_Call) Return(_a0 *[]slack.User, _a1 error) *mockISlackUserGroup_GetUsersInfoContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockISlackUserGroup_GetUsersInfoContext_Call) RunAndReturn(run func(context.Context, ...string) (*[]slack.User, error)) *mockISlackUserGroup_GetUsersInfoContext_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateUserGroupMembersContext provides a mock function with given fields: ctx, userGroup, members
func (_m *mockISlackUserGroup) UpdateUserGroupMembersContext(ctx context.Context, userGroup string, members string) (slack.UserGroup, error) {
	ret := _m.Called(ctx, userGroup, members)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUserGroupMembersContext")
	}

	var r0 slack.UserGroup
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string) (slack.UserGroup, error)); ok {
		return rf(ctx, userGroup, members)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string) slack.UserGroup); ok {
		r0 = rf(ctx, userGroup, members)
	} else {
		r0 = ret.Get(0).(slack.UserGroup)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string) error); ok {
		r1 = rf(ctx, userGroup, members)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockISlackUserGroup_UpdateUserGroupMembersContext_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUserGroupMembersContext'
type mockISlackUserGroup_UpdateUserGroupMembersContext_Call struct {
	*mock.Call
}

// UpdateUserGroupMembersContext is a helper method to define mock.On call
//   - ctx context.Context
//   - userGroup string
//   - members string
func (_e *mockISlackUserGroup_Expecter) UpdateUserGroupMembersContext(ctx interface{}, userGroup interface{}, members interface{}) *mockISlackUserGroup_UpdateUserGroupMembersContext_Call {
	return &mockISlackUserGroup_UpdateUserGroupMembersContext_Call{Call: _e.mock.On("UpdateUserGroupMembersContext", ctx, userGroup, members)}
}

func (_c *mockISlackUserGroup_UpdateUserGroupMembersContext_Call) Run(run func(ctx context.Context, userGroup string, members string)) *mockISlackUserGroup_UpdateUserGroupMembersContext_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(string), args[2].(string))
	})
	return _c
}

func (_c *mockISlackUserGroup_UpdateUserGroupMembersContext_Call) Return(_a0 slack.UserGroup, _a1 error) *mockISlackUserGroup_UpdateUserGroupMembersContext_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockISlackUserGroup_UpdateUserGroupMembersContext_Call) RunAndReturn(run func(context.Context, string, string) (slack.UserGroup, error)) *mockISlackUserGroup_UpdateUserGroupMembersContext_Call {
	_c.Call.Return(run)
	return _c
}

// newMockISlackUserGroup creates a new instance of mockISlackUserGroup. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockISlackUserGroup(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockISlackUserGroup {
	mock := &mockISlackUserGroup{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
