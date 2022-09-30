// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	slack "github.com/slack-go/slack"
	mock "github.com/stretchr/testify/mock"
)

// ISlackUserGroup is an autogenerated mock type for the iSlackUserGroup type
type ISlackUserGroup struct {
	mock.Mock
}

type ISlackUserGroup_Expecter struct {
	mock *mock.Mock
}

func (_m *ISlackUserGroup) EXPECT() *ISlackUserGroup_Expecter {
	return &ISlackUserGroup_Expecter{mock: &_m.Mock}
}

// DisableUserGroup provides a mock function with given fields: userGroup
func (_m *ISlackUserGroup) DisableUserGroup(userGroup string) (slack.UserGroup, error) {
	ret := _m.Called(userGroup)

	var r0 slack.UserGroup
	if rf, ok := ret.Get(0).(func(string) slack.UserGroup); ok {
		r0 = rf(userGroup)
	} else {
		r0 = ret.Get(0).(slack.UserGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userGroup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_DisableUserGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DisableUserGroup'
type ISlackUserGroup_DisableUserGroup_Call struct {
	*mock.Call
}

// DisableUserGroup is a helper method to define mock.On call
//   - userGroup string
func (_e *ISlackUserGroup_Expecter) DisableUserGroup(userGroup interface{}) *ISlackUserGroup_DisableUserGroup_Call {
	return &ISlackUserGroup_DisableUserGroup_Call{Call: _e.mock.On("DisableUserGroup", userGroup)}
}

func (_c *ISlackUserGroup_DisableUserGroup_Call) Run(run func(userGroup string)) *ISlackUserGroup_DisableUserGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISlackUserGroup_DisableUserGroup_Call) Return(_a0 slack.UserGroup, _a1 error) *ISlackUserGroup_DisableUserGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// EnableUserGroup provides a mock function with given fields: userGroup
func (_m *ISlackUserGroup) EnableUserGroup(userGroup string) (slack.UserGroup, error) {
	ret := _m.Called(userGroup)

	var r0 slack.UserGroup
	if rf, ok := ret.Get(0).(func(string) slack.UserGroup); ok {
		r0 = rf(userGroup)
	} else {
		r0 = ret.Get(0).(slack.UserGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userGroup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_EnableUserGroup_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'EnableUserGroup'
type ISlackUserGroup_EnableUserGroup_Call struct {
	*mock.Call
}

// EnableUserGroup is a helper method to define mock.On call
//   - userGroup string
func (_e *ISlackUserGroup_Expecter) EnableUserGroup(userGroup interface{}) *ISlackUserGroup_EnableUserGroup_Call {
	return &ISlackUserGroup_EnableUserGroup_Call{Call: _e.mock.On("EnableUserGroup", userGroup)}
}

func (_c *ISlackUserGroup_EnableUserGroup_Call) Run(run func(userGroup string)) *ISlackUserGroup_EnableUserGroup_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISlackUserGroup_EnableUserGroup_Call) Return(_a0 slack.UserGroup, _a1 error) *ISlackUserGroup_EnableUserGroup_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *ISlackUserGroup) GetUserByEmail(email string) (*slack.User, error) {
	ret := _m.Called(email)

	var r0 *slack.User
	if rf, ok := ret.Get(0).(func(string) *slack.User); ok {
		r0 = rf(email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*slack.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_GetUserByEmail_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserByEmail'
type ISlackUserGroup_GetUserByEmail_Call struct {
	*mock.Call
}

// GetUserByEmail is a helper method to define mock.On call
//   - email string
func (_e *ISlackUserGroup_Expecter) GetUserByEmail(email interface{}) *ISlackUserGroup_GetUserByEmail_Call {
	return &ISlackUserGroup_GetUserByEmail_Call{Call: _e.mock.On("GetUserByEmail", email)}
}

func (_c *ISlackUserGroup_GetUserByEmail_Call) Run(run func(email string)) *ISlackUserGroup_GetUserByEmail_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISlackUserGroup_GetUserByEmail_Call) Return(_a0 *slack.User, _a1 error) *ISlackUserGroup_GetUserByEmail_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetUserGroupMembers provides a mock function with given fields: userGroup
func (_m *ISlackUserGroup) GetUserGroupMembers(userGroup string) ([]string, error) {
	ret := _m.Called(userGroup)

	var r0 []string
	if rf, ok := ret.Get(0).(func(string) []string); ok {
		r0 = rf(userGroup)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userGroup)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_GetUserGroupMembers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUserGroupMembers'
type ISlackUserGroup_GetUserGroupMembers_Call struct {
	*mock.Call
}

// GetUserGroupMembers is a helper method to define mock.On call
//   - userGroup string
func (_e *ISlackUserGroup_Expecter) GetUserGroupMembers(userGroup interface{}) *ISlackUserGroup_GetUserGroupMembers_Call {
	return &ISlackUserGroup_GetUserGroupMembers_Call{Call: _e.mock.On("GetUserGroupMembers", userGroup)}
}

func (_c *ISlackUserGroup_GetUserGroupMembers_Call) Run(run func(userGroup string)) *ISlackUserGroup_GetUserGroupMembers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string))
	})
	return _c
}

func (_c *ISlackUserGroup_GetUserGroupMembers_Call) Return(_a0 []string, _a1 error) *ISlackUserGroup_GetUserGroupMembers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetUsersInfo provides a mock function with given fields: users
func (_m *ISlackUserGroup) GetUsersInfo(users ...string) (*[]slack.User, error) {
	_va := make([]interface{}, len(users))
	for _i := range users {
		_va[_i] = users[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *[]slack.User
	if rf, ok := ret.Get(0).(func(...string) *[]slack.User); ok {
		r0 = rf(users...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*[]slack.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...string) error); ok {
		r1 = rf(users...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_GetUsersInfo_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetUsersInfo'
type ISlackUserGroup_GetUsersInfo_Call struct {
	*mock.Call
}

// GetUsersInfo is a helper method to define mock.On call
//   - users ...string
func (_e *ISlackUserGroup_Expecter) GetUsersInfo(users ...interface{}) *ISlackUserGroup_GetUsersInfo_Call {
	return &ISlackUserGroup_GetUsersInfo_Call{Call: _e.mock.On("GetUsersInfo",
		append([]interface{}{}, users...)...)}
}

func (_c *ISlackUserGroup_GetUsersInfo_Call) Run(run func(users ...string)) *ISlackUserGroup_GetUsersInfo_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]string, len(args)-0)
		for i, a := range args[0:] {
			if a != nil {
				variadicArgs[i] = a.(string)
			}
		}
		run(variadicArgs...)
	})
	return _c
}

func (_c *ISlackUserGroup_GetUsersInfo_Call) Return(_a0 *[]slack.User, _a1 error) *ISlackUserGroup_GetUsersInfo_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// UpdateUserGroupMembers provides a mock function with given fields: userGroup, members
func (_m *ISlackUserGroup) UpdateUserGroupMembers(userGroup string, members string) (slack.UserGroup, error) {
	ret := _m.Called(userGroup, members)

	var r0 slack.UserGroup
	if rf, ok := ret.Get(0).(func(string, string) slack.UserGroup); ok {
		r0 = rf(userGroup, members)
	} else {
		r0 = ret.Get(0).(slack.UserGroup)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(userGroup, members)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ISlackUserGroup_UpdateUserGroupMembers_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateUserGroupMembers'
type ISlackUserGroup_UpdateUserGroupMembers_Call struct {
	*mock.Call
}

// UpdateUserGroupMembers is a helper method to define mock.On call
//   - userGroup string
//   - members string
func (_e *ISlackUserGroup_Expecter) UpdateUserGroupMembers(userGroup interface{}, members interface{}) *ISlackUserGroup_UpdateUserGroupMembers_Call {
	return &ISlackUserGroup_UpdateUserGroupMembers_Call{Call: _e.mock.On("UpdateUserGroupMembers", userGroup, members)}
}

func (_c *ISlackUserGroup_UpdateUserGroupMembers_Call) Run(run func(userGroup string, members string)) *ISlackUserGroup_UpdateUserGroupMembers_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(string), args[1].(string))
	})
	return _c
}

func (_c *ISlackUserGroup_UpdateUserGroupMembers_Call) Return(_a0 slack.UserGroup, _a1 error) *ISlackUserGroup_UpdateUserGroupMembers_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

type mockConstructorTestingTNewISlackUserGroup interface {
	mock.TestingT
	Cleanup(func())
}

// NewISlackUserGroup creates a new instance of ISlackUserGroup. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewISlackUserGroup(t mockConstructorTestingTNewISlackUserGroup) *ISlackUserGroup {
	mock := &ISlackUserGroup{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}