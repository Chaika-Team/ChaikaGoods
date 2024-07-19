package dao

import (
	"ChaikaGoods/internal/models"
	"context"
)

// ProductDAO определяет интерфейс для работы с товарами.
type ProductDAO interface {
	GetProductByID(ctx context.Context, id int64) (*models.Product, error)
	CreateProduct(ctx context.Context, product *models.Product) error
	UpdateProduct(ctx context.Context, product *models.Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

// PackageDAO определяет интерфейс для работы с пакетами товаров.
type PackageDAO interface {
	GetPackageByID(ctx context.Context, id int64) (*models.Package, error)
	CreatePackage(ctx context.Context, pkg *models.Package) error
	UpdatePackage(ctx context.Context, pkg *models.Package) error
	DeletePackage(ctx context.Context, id int64) error
}

// PackageContentDAO определяет интерфейс для работы с содержимым пакетов.
type PackageContentDAO interface {
	GetPackageContentsByPackageID(ctx context.Context, packageID int64) ([]models.PackageContent, error)
	AddProductToPackage(ctx context.Context, content *models.PackageContent) error
	UpdatePackageContent(ctx context.Context, content *models.PackageContent) error
	RemoveProductFromPackage(ctx context.Context, id int64) error
}
