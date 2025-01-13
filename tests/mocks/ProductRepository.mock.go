// Code generated by mockery v2.50.4. DO NOT EDIT.

package mocks

import (
	context "context"

	models "github.com/Chaika-Team/ChaikaGoods/internal/models"
	mock "github.com/stretchr/testify/mock"
)

// MockProductRepository is an autogenerated mock type for the ProductRepository type
type MockProductRepository struct {
	mock.Mock
}

type MockProductRepository_Expecter struct {
	mock *mock.Mock
}

func (_m *MockProductRepository) EXPECT() *MockProductRepository_Expecter {
	return &MockProductRepository_Expecter{mock: &_m.Mock}
}

// CreateProduct provides a mock function with given fields: ctx, p
func (_m *MockProductRepository) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 int64
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) (int64, error)); ok {
		return rf(ctx, p)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) int64); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Get(0).(int64)
	}

	if rf, ok := ret.Get(1).(func(context.Context, *models.Product) error); ok {
		r1 = rf(ctx, p)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProductRepository_CreateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateProduct'
type MockProductRepository_CreateProduct_Call struct {
	*mock.Call
}

// CreateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - p *models.Product
func (_e *MockProductRepository_Expecter) CreateProduct(ctx interface{}, p interface{}) *MockProductRepository_CreateProduct_Call {
	return &MockProductRepository_CreateProduct_Call{Call: _e.mock.On("CreateProduct", ctx, p)}
}

func (_c *MockProductRepository_CreateProduct_Call) Run(run func(ctx context.Context, p *models.Product)) *MockProductRepository_CreateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Product))
	})
	return _c
}

func (_c *MockProductRepository_CreateProduct_Call) Return(_a0 int64, _a1 error) *MockProductRepository_CreateProduct_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductRepository_CreateProduct_Call) RunAndReturn(run func(context.Context, *models.Product) (int64, error)) *MockProductRepository_CreateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// DeleteProduct provides a mock function with given fields: ctx, id
func (_m *MockProductRepository) DeleteProduct(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for DeleteProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductRepository_DeleteProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'DeleteProduct'
type MockProductRepository_DeleteProduct_Call struct {
	*mock.Call
}

// DeleteProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockProductRepository_Expecter) DeleteProduct(ctx interface{}, id interface{}) *MockProductRepository_DeleteProduct_Call {
	return &MockProductRepository_DeleteProduct_Call{Call: _e.mock.On("DeleteProduct", ctx, id)}
}

func (_c *MockProductRepository_DeleteProduct_Call) Run(run func(ctx context.Context, id int64)) *MockProductRepository_DeleteProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockProductRepository_DeleteProduct_Call) Return(_a0 error) *MockProductRepository_DeleteProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductRepository_DeleteProduct_Call) RunAndReturn(run func(context.Context, int64) error) *MockProductRepository_DeleteProduct_Call {
	_c.Call.Return(run)
	return _c
}

// GetAllProducts provides a mock function with given fields: ctx
func (_m *MockProductRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for GetAllProducts")
	}

	var r0 []models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) ([]models.Product, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) []models.Product); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProductRepository_GetAllProducts_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetAllProducts'
type MockProductRepository_GetAllProducts_Call struct {
	*mock.Call
}

// GetAllProducts is a helper method to define mock.On call
//   - ctx context.Context
func (_e *MockProductRepository_Expecter) GetAllProducts(ctx interface{}) *MockProductRepository_GetAllProducts_Call {
	return &MockProductRepository_GetAllProducts_Call{Call: _e.mock.On("GetAllProducts", ctx)}
}

func (_c *MockProductRepository_GetAllProducts_Call) Run(run func(ctx context.Context)) *MockProductRepository_GetAllProducts_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context))
	})
	return _c
}

func (_c *MockProductRepository_GetAllProducts_Call) Return(_a0 []models.Product, _a1 error) *MockProductRepository_GetAllProducts_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductRepository_GetAllProducts_Call) RunAndReturn(run func(context.Context) ([]models.Product, error)) *MockProductRepository_GetAllProducts_Call {
	_c.Call.Return(run)
	return _c
}

// GetProductByID provides a mock function with given fields: ctx, id
func (_m *MockProductRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	ret := _m.Called(ctx, id)

	if len(ret) == 0 {
		panic("no return value specified for GetProductByID")
	}

	var r0 models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) (models.Product, error)); ok {
		return rf(ctx, id)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int64) models.Product); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Product)
	}

	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockProductRepository_GetProductByID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetProductByID'
type MockProductRepository_GetProductByID_Call struct {
	*mock.Call
}

// GetProductByID is a helper method to define mock.On call
//   - ctx context.Context
//   - id int64
func (_e *MockProductRepository_Expecter) GetProductByID(ctx interface{}, id interface{}) *MockProductRepository_GetProductByID_Call {
	return &MockProductRepository_GetProductByID_Call{Call: _e.mock.On("GetProductByID", ctx, id)}
}

func (_c *MockProductRepository_GetProductByID_Call) Run(run func(ctx context.Context, id int64)) *MockProductRepository_GetProductByID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64))
	})
	return _c
}

func (_c *MockProductRepository_GetProductByID_Call) Return(_a0 models.Product, _a1 error) *MockProductRepository_GetProductByID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockProductRepository_GetProductByID_Call) RunAndReturn(run func(context.Context, int64) (models.Product, error)) *MockProductRepository_GetProductByID_Call {
	_c.Call.Return(run)
	return _c
}

// UpdateProduct provides a mock function with given fields: ctx, p
func (_m *MockProductRepository) UpdateProduct(ctx context.Context, p *models.Product) error {
	ret := _m.Called(ctx, p)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *models.Product) error); ok {
		r0 = rf(ctx, p)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockProductRepository_UpdateProduct_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'UpdateProduct'
type MockProductRepository_UpdateProduct_Call struct {
	*mock.Call
}

// UpdateProduct is a helper method to define mock.On call
//   - ctx context.Context
//   - p *models.Product
func (_e *MockProductRepository_Expecter) UpdateProduct(ctx interface{}, p interface{}) *MockProductRepository_UpdateProduct_Call {
	return &MockProductRepository_UpdateProduct_Call{Call: _e.mock.On("UpdateProduct", ctx, p)}
}

func (_c *MockProductRepository_UpdateProduct_Call) Run(run func(ctx context.Context, p *models.Product)) *MockProductRepository_UpdateProduct_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*models.Product))
	})
	return _c
}

func (_c *MockProductRepository_UpdateProduct_Call) Return(_a0 error) *MockProductRepository_UpdateProduct_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockProductRepository_UpdateProduct_Call) RunAndReturn(run func(context.Context, *models.Product) error) *MockProductRepository_UpdateProduct_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockProductRepository creates a new instance of MockProductRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockProductRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockProductRepository {
	mock := &MockProductRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
