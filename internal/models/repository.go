package models

import (
	"context"
)

// ProductRepository defines methods for product-related database operations.
type ProductRepository interface {
	GetProductByID(ctx context.Context, id int64) (Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	CreateProduct(ctx context.Context, p *Product) (int64, error)
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

// PackageRepository defines methods for package-related database operations.
type PackageRepository interface {
	GetPackageByID(ctx context.Context, id int64) (Package, error)
	GetProductsByPackageID(ctx context.Context, packageID int64) ([]PackageContent, error)
	ListPackages(ctx context.Context) ([]Package, error)
	CreatePackage(ctx context.Context, pkg *Package) error
	DeletePackage(ctx context.Context, packageID int64) error
	SearchPackages(ctx context.Context, searchString string, limit int64, offset int64) ([]Package, error)
	GetAllPackages(ctx context.Context, limit int64, offset int64) ([]Package, error)
}

// GoodsRepository объединяет репозитории для продуктов и пакетов.
type GoodsRepository interface {
	ProductRepository
	PackageRepository
}
