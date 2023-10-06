// Code generated by mockery. DO NOT EDIT.

package proto

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"
)

// MockAdapterClient is an autogenerated mock type for the AdapterClient type
type MockAdapterClient struct {
	mock.Mock
}

type MockAdapterClient_Expecter struct {
	mock *mock.Mock
}

func (_m *MockAdapterClient) EXPECT() *MockAdapterClient_Expecter {
	return &MockAdapterClient_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, in, opts
func (_m *MockAdapterClient) Add(ctx context.Context, in *AddRequest, opts ...grpc.CallOption) (*AddResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *AddResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *AddRequest, ...grpc.CallOption) (*AddResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *AddRequest, ...grpc.CallOption) *AddResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*AddResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *AddRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapterClient_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type MockAdapterClient_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - in *AddRequest
//   - opts ...grpc.CallOption
func (_e *MockAdapterClient_Expecter) Add(ctx interface{}, in interface{}, opts ...interface{}) *MockAdapterClient_Add_Call {
	return &MockAdapterClient_Add_Call{Call: _e.mock.On("Add",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAdapterClient_Add_Call) Run(run func(ctx context.Context, in *AddRequest, opts ...grpc.CallOption)) *MockAdapterClient_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*AddRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAdapterClient_Add_Call) Return(_a0 *AddResponse, _a1 error) *MockAdapterClient_Add_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapterClient_Add_Call) RunAndReturn(run func(context.Context, *AddRequest, ...grpc.CallOption) (*AddResponse, error)) *MockAdapterClient_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *MockAdapterClient) Get(ctx context.Context, in *GetRequest, opts ...grpc.CallOption) (*GetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *GetRequest, ...grpc.CallOption) (*GetResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *GetRequest, ...grpc.CallOption) *GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapterClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type MockAdapterClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - in *GetRequest
//   - opts ...grpc.CallOption
func (_e *MockAdapterClient_Expecter) Get(ctx interface{}, in interface{}, opts ...interface{}) *MockAdapterClient_Get_Call {
	return &MockAdapterClient_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAdapterClient_Get_Call) Run(run func(ctx context.Context, in *GetRequest, opts ...grpc.CallOption)) *MockAdapterClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*GetRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAdapterClient_Get_Call) Return(_a0 *GetResponse, _a1 error) *MockAdapterClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapterClient_Get_Call) RunAndReturn(run func(context.Context, *GetRequest, ...grpc.CallOption) (*GetResponse, error)) *MockAdapterClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with given fields: ctx, in, opts
func (_m *MockAdapterClient) Init(ctx context.Context, in *InitRequest, opts ...grpc.CallOption) (*InitResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *InitResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *InitRequest, ...grpc.CallOption) (*InitResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *InitRequest, ...grpc.CallOption) *InitResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*InitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *InitRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapterClient_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type MockAdapterClient_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - ctx context.Context
//   - in *InitRequest
//   - opts ...grpc.CallOption
func (_e *MockAdapterClient_Expecter) Init(ctx interface{}, in interface{}, opts ...interface{}) *MockAdapterClient_Init_Call {
	return &MockAdapterClient_Init_Call{Call: _e.mock.On("Init",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAdapterClient_Init_Call) Run(run func(ctx context.Context, in *InitRequest, opts ...grpc.CallOption)) *MockAdapterClient_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*InitRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAdapterClient_Init_Call) Return(_a0 *InitResponse, _a1 error) *MockAdapterClient_Init_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapterClient_Init_Call) RunAndReturn(run func(context.Context, *InitRequest, ...grpc.CallOption) (*InitResponse, error)) *MockAdapterClient_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, in, opts
func (_m *MockAdapterClient) Remove(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption) (*RemoveResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *RemoveResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *RemoveRequest, ...grpc.CallOption) (*RemoveResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *RemoveRequest, ...grpc.CallOption) *RemoveResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*RemoveResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *RemoveRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockAdapterClient_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type MockAdapterClient_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - in *RemoveRequest
//   - opts ...grpc.CallOption
func (_e *MockAdapterClient_Expecter) Remove(ctx interface{}, in interface{}, opts ...interface{}) *MockAdapterClient_Remove_Call {
	return &MockAdapterClient_Remove_Call{Call: _e.mock.On("Remove",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *MockAdapterClient_Remove_Call) Run(run func(ctx context.Context, in *RemoveRequest, opts ...grpc.CallOption)) *MockAdapterClient_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*RemoveRequest), variadicArgs...)
	})
	return _c
}

func (_c *MockAdapterClient_Remove_Call) Return(_a0 *RemoveResponse, _a1 error) *MockAdapterClient_Remove_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockAdapterClient_Remove_Call) RunAndReturn(run func(context.Context, *RemoveRequest, ...grpc.CallOption) (*RemoveResponse, error)) *MockAdapterClient_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockAdapterClient creates a new instance of MockAdapterClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockAdapterClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockAdapterClient {
	mock := &MockAdapterClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
