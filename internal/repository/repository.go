package repository

import (
	"context"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

// ProductRepository defines methods for product-related database operations.
type ProductRepository interface {
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	CreateProduct(ctx context.Context, p *models.Product) (int64, error)
	UpdateProduct(ctx context.Context, p *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

// PackageRepository defines methods for package-related database operations.
type PackageRepository interface {
	GetPackageByID(ctx context.Context, id int64) (models.Package, error)
	GetProductsByPackageID(ctx context.Context, packageID int64) ([]models.PackageContent, error)
	ListPackages(ctx context.Context) ([]models.Package, error)
	CreatePackage(ctx context.Context, pkg *models.Package) error
	DeletePackage(ctx context.Context, packageID int64) error
	SearchPackages(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Package, error)
	GetAllPackages(ctx context.Context, limit int64, offset int64) ([]models.Package, error)
}

// GoodsRepository объединяет репозитории для продуктов и пакетов.
type GoodsRepository interface {
	ProductRepository
	PackageRepository
}
