// Code generated by mockery. DO NOT EDIT.

package user

import (
	context "context"

	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	mock "github.com/stretchr/testify/mock"

	users "github.com/microsoftgraph/msgraph-sdk-go/users"
)

// mockIUser is an autogenerated mock type for the iUser type
type mockIUser struct {
	mock.Mock
}

type mockIUser_Expecter struct {
	mock *mock.Mock
}

func (_m *mockIUser) EXPECT() *mockIUser_Expecter {
	return &mockIUser_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: _a0, _a1
func (_m *mockIUser) Get(_a0 context.Context, _a1 *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error) {
	ret := _m.Called(_a0, _a1)

	var r0 models.UserCollectionResponseable
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) models.UserCollectionResponseable); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.UserCollectionResponseable)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockIUser_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockIUser_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *users.UsersRequestBuilderGetRequestConfiguration
func (_e *mockIUser_Expecter) Get(_a0 interface{}, _a1 interface{}) *mockIUser_Get_Call {
	return &mockIUser_Get_Call{Call: _e.mock.On("Get", _a0, _a1)}
}

func (_c *mockIUser_Get_Call) Run(run func(_a0 context.Context, _a1 *users.UsersRequestBuilderGetRequestConfiguration)) *mockIUser_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*users.UsersRequestBuilderGetRequestConfiguration))
	})
	return _c
}

func (_c *mockIUser_Get_Call) Return(_a0 models.UserCollectionResponseable, _a1 error) *mockIUser_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockIUser_Get_Call) RunAndReturn(run func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error)) *mockIUser_Get_Call {
	_c.Call.Return(run)
	return _c
}

// newMockIUser creates a new instance of mockIUser. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockIUser(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockIUser {
	mock := &mockIUser{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
