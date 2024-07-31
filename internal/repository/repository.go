package repository

import (
	"ChaikaGoods/internal/models"
	"context"
)

type GoodsRepository interface {
	// GetProductByID returns a product by its ID.
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	// GetAllProducts returns all products.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// AddQueryToCreateProduct adds a query to create a product.
	AddQueryToCreateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error)
	// AddQueryToUpdateProduct adds a query to update a product.
	AddQueryToUpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error)
	// AddQueryToDeleteProduct adds a query to delete a product.
	AddQueryToDeleteProduct(ctx context.Context, id int64) (changeID int64, err error)
	// ApplyChanges applies changes to the database.
	ApplyChanges(ctx context.Context, version *models.Version) error
	// CreateNewDevVersion creates a new development version of product data.
	CreateNewDevVersion(ctx context.Context) (models.Version, error)
	// GetAllChanges returns all changes for a given version.
	GetAllChanges(ctx context.Context, version models.Version) ([]models.Change, error)
	// GetCurrentDevVersion returns the current development version.
	GetCurrentDevVersion(ctx context.Context) (models.Version, error)
	// GetCurrentActualVersion returns the current actual version.
	GetCurrentActualVersion(ctx context.Context) (models.Version, error)
	// GetVersionsBetween returns all versions between two versions.
	GetVersionsBetween(ctx context.Context, from, to int) ([]models.Version, error)
	// DeleteChange deletes a change by its ID.
	DeleteChange(ctx context.Context, id int64) error
	// GetPackageByID returns a package by its ID.
	GetPackageByID(ctx context.Context, p *models.Package) ([]models.PackageContent, error)
	// GetProductsByPackageID returns all products in a package.
	GetProductsByPackageID(ctx context.Context, p *models.Package) ([]models.PackageContent, error)
	// ListPackages returns all packages.
	ListPackages(ctx context.Context) ([]models.Package, error)
	// CreatePackage creates a new package.
	CreatePackage(ctx context.Context, pkg *models.Package) error
	// AddProductToPackage adds a product to a package.
	AddProductToPackage(ctx context.Context, packageID int64, products []models.PackageContent) error
	// DeletePackage deletes a package by its ID.
	DeletePackage(ctx context.Context, packageID int64) error
	// SearchPacket searches for a packet by its name.
	SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error)
}
