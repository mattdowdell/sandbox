// Code generated by mockery. DO NOT EDIT.

package mockusecasefacades

import (
	context "context"

	entities "github.com/mattdowdell/sandbox/internal/domain/entities"
	mock "github.com/stretchr/testify/mock"

	repositories "github.com/mattdowdell/sandbox/internal/domain/repositories"
)

// ResourceLister is an autogenerated mock type for the ResourceLister type
type ResourceLister struct {
	mock.Mock
}

type ResourceLister_Expecter struct {
	mock *mock.Mock
}

func (_m *ResourceLister) EXPECT() *ResourceLister_Expecter {
	return &ResourceLister_Expecter{mock: &_m.Mock}
}

// Execute provides a mock function with given fields: _a0, _a1
func (_m *ResourceLister) Execute(_a0 context.Context, _a1 repositories.Resource) ([]*entities.Resource, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for Execute")
	}

	var r0 []*entities.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, repositories.Resource) ([]*entities.Resource, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, repositories.Resource) []*entities.Resource); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, repositories.Resource) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ResourceLister_Execute_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Execute'
type ResourceLister_Execute_Call struct {
	*mock.Call
}

// Execute is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 repositories.Resource
func (_e *ResourceLister_Expecter) Execute(_a0 interface{}, _a1 interface{}) *ResourceLister_Execute_Call {
	return &ResourceLister_Execute_Call{Call: _e.mock.On("Execute", _a0, _a1)}
}

func (_c *ResourceLister_Execute_Call) Run(run func(_a0 context.Context, _a1 repositories.Resource)) *ResourceLister_Execute_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(repositories.Resource))
	})
	return _c
}

func (_c *ResourceLister_Execute_Call) Return(_a0 []*entities.Resource, _a1 error) *ResourceLister_Execute_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ResourceLister_Execute_Call) RunAndReturn(run func(context.Context, repositories.Resource) ([]*entities.Resource, error)) *ResourceLister_Execute_Call {
	_c.Call.Return(run)
	return _c
}

// NewResourceLister creates a new instance of ResourceLister. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResourceLister(t interface {
	mock.TestingT
	Cleanup(func())
}) *ResourceLister {
	mock := &ResourceLister{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
