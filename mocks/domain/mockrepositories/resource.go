// Code generated by mockery. DO NOT EDIT.

package mockrepositories

import (
	context "context"

	entities "github.com/mattdowdell/sandbox/internal/domain/entities"
	mock "github.com/stretchr/testify/mock"

	uuid "github.com/google/uuid"
)

// Resource is an autogenerated mock type for the Resource type
type Resource struct {
	mock.Mock
}

type Resource_Expecter struct {
	mock *mock.Mock
}

func (_m *Resource) EXPECT() *Resource_Expecter {
	return &Resource_Expecter{mock: &_m.Mock}
}

// CreateResource provides a mock function with given fields: _a0, _a1
func (_m *Resource) CreateResource(_a0 context.Context, _a1 *entities.Resource) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateResource")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Resource) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Resource_CreateResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateResource'
type Resource_CreateResource_Call struct {
	*mock.Call
}

// CreateResource is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *entities.Resource
func (_e *Resource_Expecter) CreateResource(_a0 interface{}, _a1 interface{}) *Resource_CreateResource_Call {
	return &Resource_CreateResource_Call{Call: _e.mock.On("CreateResource", _a0, _a1)}
}

func (_c *Resource_CreateResource_Call) Run(run func(_a0 context.Context, _a1 *entities.Resource)) *Resource_CreateResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.Resource))
	})
	return _c
}

func (_c *Resource_CreateResource_Call) Return(_a0 error) *Resource_CreateResource_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Resource_CreateResource_Call) RunAndReturn(run func(context.Context, *entities.Resource) error) *Resource_CreateResource_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteResource provides a mock function with given fields: _a0, _a1
func (_m *Resource) DeleteResource(_a0 context.Context, _a1 uuid.UUID) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for DeleteResource")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Resource_DeleteResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteResource'
type Resource_DeleteResource_Call struct {
	*mock.Call
}

// DeleteResource is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 uuid.UUID
func (_e *Resource_Expecter) DeleteResource(_a0 interface{}, _a1 interface{}) *Resource_DeleteResource_Call {
	return &Resource_DeleteResource_Call{Call: _e.mock.On("DeleteResource", _a0, _a1)}
}

func (_c *Resource_DeleteResource_Call) Run(run func(_a0 context.Context, _a1 uuid.UUID)) *Resource_DeleteResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *Resource_DeleteResource_Call) Return(_a0 error) *Resource_DeleteResource_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Resource_DeleteResource_Call) RunAndReturn(run func(context.Context, uuid.UUID) error) *Resource_DeleteResource_Call {
	_c.Call.Return(run)
	return _c
}

// GetResource provides a mock function with given fields: _a0, _a1
func (_m *Resource) GetResource(_a0 context.Context, _a1 uuid.UUID) (*entities.Resource, error) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for GetResource")
	}

	var r0 *entities.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) (*entities.Resource, error)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uuid.UUID) *entities.Resource); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entities.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uuid.UUID) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Resource_GetResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetResource'
type Resource_GetResource_Call struct {
	*mock.Call
}

// GetResource is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 uuid.UUID
func (_e *Resource_Expecter) GetResource(_a0 interface{}, _a1 interface{}) *Resource_GetResource_Call {
	return &Resource_GetResource_Call{Call: _e.mock.On("GetResource", _a0, _a1)}
}

func (_c *Resource_GetResource_Call) Run(run func(_a0 context.Context, _a1 uuid.UUID)) *Resource_GetResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(uuid.UUID))
	})
	return _c
}

func (_c *Resource_GetResource_Call) Return(_a0 *entities.Resource, _a1 error) *Resource_GetResource_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Resource_GetResource_Call) RunAndReturn(run func(context.Context, uuid.UUID) (*entities.Resource, error)) *Resource_GetResource_Call {
	_c.Call.Return(run)
	return _c
}

// ListResources provides a mock function with given fields: _a0
func (_m *Resource) ListResources(_a0 context.Context) ([]*entities.Resource, error) {
	ret := _m.Called(_a0)

	if len(ret) == 0 {
		panic("no return value specified for ListResources")
	}

	var r0 []*entities.Resource
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]*entities.Resource, error)); ok {
		return rf(_a0)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []*entities.Resource); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*entities.Resource)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Resource_ListResources_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ListResources'
type Resource_ListResources_Call struct {
	*mock.Call
}

// ListResources is a helper method to define mock.On call
//   - _a0 context.Context
func (_e *Resource_Expecter) ListResources(_a0 interface{}) *Resource_ListResources_Call {
	return &Resource_ListResources_Call{Call: _e.mock.On("ListResources", _a0)}
}

func (_c *Resource_ListResources_Call) Run(run func(_a0 context.Context)) *Resource_ListResources_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *Resource_ListResources_Call) Return(_a0 []*entities.Resource, _a1 error) *Resource_ListResources_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Resource_ListResources_Call) RunAndReturn(run func(context.Context) ([]*entities.Resource, error)) *Resource_ListResources_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateResource provides a mock function with given fields: _a0, _a1
func (_m *Resource) UpdateResource(_a0 context.Context, _a1 *entities.Resource) error {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for UpdateResource")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *entities.Resource) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Resource_UpdateResource_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateResource'
type Resource_UpdateResource_Call struct {
	*mock.Call
}

// UpdateResource is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *entities.Resource
func (_e *Resource_Expecter) UpdateResource(_a0 interface{}, _a1 interface{}) *Resource_UpdateResource_Call {
	return &Resource_UpdateResource_Call{Call: _e.mock.On("UpdateResource", _a0, _a1)}
}

func (_c *Resource_UpdateResource_Call) Run(run func(_a0 context.Context, _a1 *entities.Resource)) *Resource_UpdateResource_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*entities.Resource))
	})
	return _c
}

func (_c *Resource_UpdateResource_Call) Return(_a0 error) *Resource_UpdateResource_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *Resource_UpdateResource_Call) RunAndReturn(run func(context.Context, *entities.Resource) error) *Resource_UpdateResource_Call {
	_c.Call.Return(run)
	return _c
}

// NewResource creates a new instance of Resource. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewResource(t interface {
	mock.TestingT
	Cleanup(func())
}) *Resource {
	mock := &Resource{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
