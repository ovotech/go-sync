// Code generated by mockery. DO NOT EDIT.

package groupmembership

import (
	context "context"

	models "github.com/microsoftgraph/msgraph-sdk-go/models"
	users "github.com/microsoftgraph/msgraph-sdk-go/users"
	mock "github.com/stretchr/testify/mock"
)

// mockIUserClient is an autogenerated mock type for the iUserClient type
type mockIUserClient struct {
	mock.Mock
}

type mockIUserClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockIUserClient) EXPECT() *mockIUserClient_Expecter {
	return &mockIUserClient_Expecter{mock: &_m.Mock}
}

// Get provides a mock function with given fields: ctx, config
func (_m *mockIUserClient) Get(ctx context.Context, config *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error) {
	ret := _m.Called(ctx, config)

	if len(ret) == 0 {
		panic("no return value specified for Get")
	}

	var r0 models.UserCollectionResponseable
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error)); ok {
		return rf(ctx, config)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) models.UserCollectionResponseable); ok {
		r0 = rf(ctx, config)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(models.UserCollectionResponseable)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) error); ok {
		r1 = rf(ctx, config)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockIUserClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockIUserClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - config *users.UsersRequestBuilderGetRequestConfiguration
func (_e *mockIUserClient_Expecter) Get(ctx interface{}, config interface{}) *mockIUserClient_Get_Call {
	return &mockIUserClient_Get_Call{Call: _e.mock.On("Get", ctx, config)}
}

func (_c *mockIUserClient_Get_Call) Run(run func(ctx context.Context, config *users.UsersRequestBuilderGetRequestConfiguration)) *mockIUserClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*users.UsersRequestBuilderGetRequestConfiguration))
	})
	return _c
}

func (_c *mockIUserClient_Get_Call) Return(_a0 models.UserCollectionResponseable, _a1 error) *mockIUserClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockIUserClient_Get_Call) RunAndReturn(run func(context.Context, *users.UsersRequestBuilderGetRequestConfiguration) (models.UserCollectionResponseable, error)) *mockIUserClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// newMockIUserClient creates a new instance of mockIUserClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockIUserClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockIUserClient {
	mock := &mockIUserClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}