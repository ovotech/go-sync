// Code generated by mockery. DO NOT EDIT.

package gosync

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockInitFn is an autogenerated mock type for the InitFn type
type MockInitFn[T Adapter] struct {
	mock.Mock
}

type MockInitFn_Expecter[T Adapter] struct {
	mock *mock.Mock
}

func (_m *MockInitFn[T]) EXPECT() *MockInitFn_Expecter[T] {
	return &MockInitFn_Expecter[T]{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0, _a1, _a2
func (_m *MockInitFn[T]) Execute(_a0 context.Context, _a1 map[string]string, _a2 ...ConfigFn[T]) (T, error) {
	_va := make([]interface{}, len(_a2))
	for _i := range _a2 {
		_va[_i] = _a2[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, _a0, _a1)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 T
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string, ...ConfigFn[T]) (T, error)); ok {
		return rf(_a0, _a1, _a2...)
	}
	if rf, ok := ret.Get(0).(func(context.Context, map[string]string, ...ConfigFn[T]) T); ok {
		r0 = rf(_a0, _a1, _a2...)
	} else {
		r0 = ret.Get(0).(T)
	}

	if rf, ok := ret.Get(1).(func(context.Context, map[string]string, ...ConfigFn[T]) error); ok {
		r1 = rf(_a0, _a1, _a2...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockInitFn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockInitFn_Execute_Call[T Adapter] struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 map[string]string
//   - _a2 ...ConfigFn[T]
func (_e *MockInitFn_Expecter[T]) Execute(_a0 interface{}, _a1 interface{}, _a2 ...interface{}) *MockInitFn_Execute_Call[T] {
	return &MockInitFn_Execute_Call[T]{Call: _e.mock.On("Execute",
		append([]interface{}{_a0, _a1}, _a2...)...)}
}

func (_c *MockInitFn_Execute_Call[T]) Run(run func(_a0 context.Context, _a1 map[string]string, _a2 ...ConfigFn[T])) *MockInitFn_Execute_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]ConfigFn[T], len(args)-2)
		for i, a := range args[2:] {
			if a != nil {
				variadicArgs[i] = a.(ConfigFn[T])
			}
		}
		run(args[0].(context.Context), args[1].(map[string]string), variadicArgs...)
	})
	return _c
}

func (_c *MockInitFn_Execute_Call[T]) Return(_a0 T, _a1 error) *MockInitFn_Execute_Call[T] {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockInitFn_Execute_Call[T]) RunAndReturn(run func(context.Context, map[string]string, ...ConfigFn[T]) (T, error)) *MockInitFn_Execute_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewMockInitFn creates a new instance of MockInitFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockInitFn[T Adapter](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockInitFn[T] {
	mock := &MockInitFn[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
