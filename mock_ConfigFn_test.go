// Code generated by mockery. DO NOT EDIT.

package gosync

import mock "github.com/stretchr/testify/mock"

// MockConfigFn is an autogenerated mock type for the ConfigFn type
type MockConfigFn[T Adapter] struct {
	mock.Mock
}

type MockConfigFn_Expecter[T Adapter] struct {
	mock *mock.Mock
}

func (_m *MockConfigFn[T]) EXPECT() *MockConfigFn_Expecter[T] {
	return &MockConfigFn_Expecter[T]{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0
func (_m *MockConfigFn[T]) Execute(_a0 T) {
	_m.Called(_a0)
}

// MockConfigFn_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type MockConfigFn_Execute_Call[T Adapter] struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 T
func (_e *MockConfigFn_Expecter[T]) Execute(_a0 interface{}) *MockConfigFn_Execute_Call[T] {
	return &MockConfigFn_Execute_Call[T]{Call: _e.mock.On("Execute", _a0)}
}

func (_c *MockConfigFn_Execute_Call[T]) Run(run func(_a0 T)) *MockConfigFn_Execute_Call[T] {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(T))
	})
	return _c
}

func (_c *MockConfigFn_Execute_Call[T]) Return() *MockConfigFn_Execute_Call[T] {
	_c.Call.Return()
	return _c
}

func (_c *MockConfigFn_Execute_Call[T]) RunAndReturn(run func(T)) *MockConfigFn_Execute_Call[T] {
	_c.Call.Return(run)
	return _c
}

// NewMockConfigFn creates a new instance of MockConfigFn. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockConfigFn[T Adapter](t interface {
	mock.TestingT
	Cleanup(func())
}) *MockConfigFn[T] {
	mock := &MockConfigFn[T]{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
