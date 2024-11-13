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
	CreateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error)
	UpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error)
	DeleteProduct(ctx context.Context, id int64) (changeID int64, err error)
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
