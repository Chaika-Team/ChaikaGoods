package mocks

import (
	"ChaikaGoods/internal/models"
	"context"

	"github.com/stretchr/testify/mock"
)

// MockGoodsRepository - это мок реализации интерфейса GoodsRepository
type MockGoodsRepository struct {
	mock.Mock
}

func (m *MockGoodsRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Product), args.Error(1)
}

func (m *MockGoodsRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Product), args.Error(1)
}

func (m *MockGoodsRepository) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	args := m.Called(ctx, p)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockGoodsRepository) UpdateProduct(ctx context.Context, p *models.Product) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockGoodsRepository) DeleteProduct(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockGoodsRepository) GetPackageByID(ctx context.Context, p *models.Package) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

func (m *MockGoodsRepository) GetProductsByPackageID(ctx context.Context, p *models.Package) ([]models.PackageContent, error) {
	args := m.Called(ctx, p)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.PackageContent), args.Error(1)
}

func (m *MockGoodsRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockGoodsRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	args := m.Called(ctx, pkg)
	return args.Error(0)
}

func (m *MockGoodsRepository) DeletePackage(ctx context.Context, packageID int64) error {
	args := m.Called(ctx, packageID)
	return args.Error(0)
}

func (m *MockGoodsRepository) SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error) {
	args := m.Called(ctx, searchString, quantity, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Package), args.Error(1)
}

func (m *MockGoodsRepository) GetAllPackages(ctx context.Context, quantity int64, offset int64) ([]models.Package, error) {
	args := m.Called(ctx, quantity, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Package), args.Error(1)
}
