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
	// CreateProduct creates a new product.
	CreateProduct(ctx context.Context, p *models.Product) (int64, error)
	// UpdateProduct updates a product by its ID.
	UpdateProduct(ctx context.Context, p *models.Product) error
	// DeleteProduct deletes a product by its ID.
	DeleteProduct(ctx context.Context, id int64) error
	// GetPackageByID returns a package by its ID.
	GetPackageByID(ctx context.Context, p *models.Package) error
	// GetProductsByPackageID returns all products in a package.
	GetProductsByPackageID(ctx context.Context, p *models.Package) ([]models.PackageContent, error)
	// ListPackages returns all packages.
	ListPackages(ctx context.Context) ([]models.Package, error)
	// CreatePackage creates a new package.
	CreatePackage(ctx context.Context, pkg *models.Package) error
	// DeletePackage deletes a package by its ID.
	DeletePackage(ctx context.Context, packageID int64) error
	// SearchPacket searches for a packet by its name.
	SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error)
	// GetAllPackages returns all packages with pagination.
	GetAllPackages(ctx context.Context, quantity int64, offset int64) ([]models.Package, error)
}
