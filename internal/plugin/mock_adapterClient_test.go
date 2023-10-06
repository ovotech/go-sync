// Code generated by mockery. DO NOT EDIT.

package plugin

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	grpc "google.golang.org/grpc"

	proto "github.com/ovotech/go-sync/internal/proto"
)

// mockAdapterClient is an autogenerated mock type for the adapterClient type
type mockAdapterClient struct {
	mock.Mock
}

type mockAdapterClient_Expecter struct {
	mock *mock.Mock
}

func (_m *mockAdapterClient) EXPECT() *mockAdapterClient_Expecter {
	return &mockAdapterClient_Expecter{mock: &_m.Mock}
}

// Add provides a mock function with given fields: ctx, in, opts
func (_m *mockAdapterClient) Add(ctx context.Context, in *proto.AddRequest, opts ...grpc.CallOption) (*proto.AddResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.AddResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AddRequest, ...grpc.CallOption) (*proto.AddResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.AddRequest, ...grpc.CallOption) *proto.AddResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.AddResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.AddRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockAdapterClient_Add_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Add'
type mockAdapterClient_Add_Call struct {
	*mock.Call
}

// Add is a helper method to define mock.On call
//   - ctx context.Context
//   - in *proto.AddRequest
//   - opts ...grpc.CallOption
func (_e *mockAdapterClient_Expecter) Add(ctx interface{}, in interface{}, opts ...interface{}) *mockAdapterClient_Add_Call {
	return &mockAdapterClient_Add_Call{Call: _e.mock.On("Add",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *mockAdapterClient_Add_Call) Run(run func(ctx context.Context, in *proto.AddRequest, opts ...grpc.CallOption)) *mockAdapterClient_Add_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*proto.AddRequest), variadicArgs...)
	})
	return _c
}

func (_c *mockAdapterClient_Add_Call) Return(_a0 *proto.AddResponse, _a1 error) *mockAdapterClient_Add_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockAdapterClient_Add_Call) RunAndReturn(run func(context.Context, *proto.AddRequest, ...grpc.CallOption) (*proto.AddResponse, error)) *mockAdapterClient_Add_Call {
	_c.Call.Return(run)
	return _c
}

// Get provides a mock function with given fields: ctx, in, opts
func (_m *mockAdapterClient) Get(ctx context.Context, in *proto.GetRequest, opts ...grpc.CallOption) (*proto.GetResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.GetResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetRequest, ...grpc.CallOption) (*proto.GetResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.GetRequest, ...grpc.CallOption) *proto.GetResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.GetResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.GetRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockAdapterClient_Get_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Get'
type mockAdapterClient_Get_Call struct {
	*mock.Call
}

// Get is a helper method to define mock.On call
//   - ctx context.Context
//   - in *proto.GetRequest
//   - opts ...grpc.CallOption
func (_e *mockAdapterClient_Expecter) Get(ctx interface{}, in interface{}, opts ...interface{}) *mockAdapterClient_Get_Call {
	return &mockAdapterClient_Get_Call{Call: _e.mock.On("Get",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *mockAdapterClient_Get_Call) Run(run func(ctx context.Context, in *proto.GetRequest, opts ...grpc.CallOption)) *mockAdapterClient_Get_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*proto.GetRequest), variadicArgs...)
	})
	return _c
}

func (_c *mockAdapterClient_Get_Call) Return(_a0 *proto.GetResponse, _a1 error) *mockAdapterClient_Get_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockAdapterClient_Get_Call) RunAndReturn(run func(context.Context, *proto.GetRequest, ...grpc.CallOption) (*proto.GetResponse, error)) *mockAdapterClient_Get_Call {
	_c.Call.Return(run)
	return _c
}

// Init provides a mock function with given fields: ctx, in, opts
func (_m *mockAdapterClient) Init(ctx context.Context, in *proto.InitRequest, opts ...grpc.CallOption) (*proto.InitResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.InitResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.InitRequest, ...grpc.CallOption) (*proto.InitResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.InitRequest, ...grpc.CallOption) *proto.InitResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.InitResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.InitRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockAdapterClient_Init_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Init'
type mockAdapterClient_Init_Call struct {
	*mock.Call
}

// Init is a helper method to define mock.On call
//   - ctx context.Context
//   - in *proto.InitRequest
//   - opts ...grpc.CallOption
func (_e *mockAdapterClient_Expecter) Init(ctx interface{}, in interface{}, opts ...interface{}) *mockAdapterClient_Init_Call {
	return &mockAdapterClient_Init_Call{Call: _e.mock.On("Init",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *mockAdapterClient_Init_Call) Run(run func(ctx context.Context, in *proto.InitRequest, opts ...grpc.CallOption)) *mockAdapterClient_Init_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*proto.InitRequest), variadicArgs...)
	})
	return _c
}

func (_c *mockAdapterClient_Init_Call) Return(_a0 *proto.InitResponse, _a1 error) *mockAdapterClient_Init_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockAdapterClient_Init_Call) RunAndReturn(run func(context.Context, *proto.InitRequest, ...grpc.CallOption) (*proto.InitResponse, error)) *mockAdapterClient_Init_Call {
	_c.Call.Return(run)
	return _c
}

// Remove provides a mock function with given fields: ctx, in, opts
func (_m *mockAdapterClient) Remove(ctx context.Context, in *proto.RemoveRequest, opts ...grpc.CallOption) (*proto.RemoveResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *proto.RemoveResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *proto.RemoveRequest, ...grpc.CallOption) (*proto.RemoveResponse, error)); ok {
		return rf(ctx, in, opts...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *proto.RemoveRequest, ...grpc.CallOption) *proto.RemoveResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*proto.RemoveResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *proto.RemoveRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// mockAdapterClient_Remove_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Remove'
type mockAdapterClient_Remove_Call struct {
	*mock.Call
}

// Remove is a helper method to define mock.On call
//   - ctx context.Context
//   - in *proto.RemoveRequest
//   - opts ...grpc.CallOption
func (_e *mockAdapterClient_Expecter) Remove(ctx interface{}, in interface{}, opts ...interface{}) *mockAdapterClient_Remove_Call {
	return &mockAdapterClient_Remove_Call{Call: _e.mock.On("Remove",
		append([]interface{}{ctx, in}, opts...)...)}
}

func (_c *mockAdapterClient_Remove_Call) Run(run func(ctx context.Context, in *proto.RemoveRequest, opts ...grpc.CallOption)) *mockAdapterClient_Remove_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]grpc.CallOption, len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(grpc.CallOption)
			}
		}
		run(args[0].(context.Context), args[1].(*proto.RemoveRequest), variadicArgs...)
	})
	return _c
}

func (_c *mockAdapterClient_Remove_Call) Return(_a0 *proto.RemoveResponse, _a1 error) *mockAdapterClient_Remove_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *mockAdapterClient_Remove_Call) RunAndReturn(run func(context.Context, *proto.RemoveRequest, ...grpc.CallOption) (*proto.RemoveResponse, error)) *mockAdapterClient_Remove_Call {
	_c.Call.Return(run)
	return _c
}

// newMockAdapterClient creates a new instance of mockAdapterClient. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func newMockAdapterClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *mockAdapterClient {
	mock := &mockAdapterClient{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
