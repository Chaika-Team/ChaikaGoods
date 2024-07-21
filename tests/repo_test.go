package tests

import (
	"ChaikaGoods/internal/config"
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/repository/postgresql"
	"context"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"testing"
)

var dbPool *pgxpool.Pool

func TestMain(m *testing.M) {
	// Setup
	ctx := context.Background()
	logger := log.NewNopLogger() // Using no-op logger for testing

	// Read configuration for test database
	cfg := config.GetConfigWithPath(logger, "C:\\Users\\Fasci\\GolandProjects\\ChaikaGoods\\config_test.yml")

	// Initialize database connection
	var err error
	dbPool, err = postgresql.NewClient(ctx, cfg.Storage, 3)
	if err != nil {
		panic("Failed to connect to the test database")
	}
	defer dbPool.Close()
	m.Run()

	// Teardown logic could be added here if necessary
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	// Add test packages
	pkg := &models.Package{PackageName: "Test Package", Description: "A test package"}
	err = repo.CreatePackage(context.Background(), pkg)
	if err != nil {
		panic("Failed to create test package")
	}
	err = repo.AddProductToPackage(context.Background(), pkg.ID, []models.PackageContent{
		{PackageID: pkg.ID, ProductID: 1, Quantity: 10}, {PackageID: pkg.ID, ProductID: 2, Quantity: 5},
	})
	if err != nil {
		panic("Failed to add products to test package")
	}
}

// TestGetProductByID tests the GetProductByID method of the GoodsRepository.
func TestGetProductByID(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	expected := &models.Product{Name: "Test Product", Description: "A test product", Price: 10.00, ImageURL: "images/image.jpg", SKU: "TESTSKU100"}
	err := repo.CreateProduct(context.Background(), expected)
	assert.NoError(t, err, "Failed to create test product")
	product, err := repo.GetProductByID(context.Background(), expected.ID)
	assert.NoError(t, err, "Failed to execute GetProductByID")
	assert.NotNil(t, product)
	assert.Equal(t, expected.Name, product.Name)
	assert.Equal(t, expected.Description, product.Description)
	assert.Equal(t, expected.Price, product.Price)
	assert.Equal(t, expected.ImageURL, product.ImageURL)
	assert.Equal(t, expected.SKU, product.SKU)

	err = repo.DeleteProduct(context.Background(), expected.ID)
	assert.NoError(t, err, "Failed to delete test product")
}

// TestGetAllProducts tests the GetAllProducts method of the GoodsRepository.
func TestGetAllProducts(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	products, err := repo.GetAllProducts(context.Background())
	assert.NoError(t, err, "Failed to execute GetAllProducts")
	assert.NotNil(t, products)
	assert.NotEmpty(t, products)
}

// TestUpdateProduct tests the UpdateProduct method of the GoodsRepository.
func TestUpdateProduct(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	// Create a product to update
	product := &models.Product{Name: "Update Product", Description: "Before update", Price: 20.00, ImageURL: "images/before.jpg", SKU: "UPDATESKU100"}
	err := repo.CreateProduct(context.Background(), product)
	assert.NoError(t, err, "Failed to create product for update test")

	// Update product details
	product.Name = "Updated Product"
	product.Description = "After update"
	product.Price = 25.00
	product.ImageURL = "images/after.jpg"
	err = repo.UpdateProduct(context.Background(), product)
	assert.NoError(t, err, "Failed to update product")

	// Retrieve the updated product to verify changes
	updatedProduct, err := repo.GetProductByID(context.Background(), product.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Updated Product", updatedProduct.Name)
	assert.Equal(t, "After update", updatedProduct.Description)
	assert.Equal(t, 25.00, updatedProduct.Price)
	assert.Equal(t, "images/after.jpg", updatedProduct.ImageURL)

	// Cleanup
	err = repo.DeleteProduct(context.Background(), product.ID)
	assert.NoError(t, err, "Failed to delete product after update test")
}

// TestDeleteProduct tests the DeleteProduct method of the GoodsRepository.
func TestDeleteProduct(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	// Create a product to delete
	product := &models.Product{Name: "Product to Delete", Description: "A product to be deleted", Price: 30.00, ImageURL: "images/delete.jpg", SKU: "DELETESKU100"}
	err := repo.CreateProduct(context.Background(), product)
	assert.NoError(t, err, "Failed to create product for deletion test")

	// Delete the product
	err = repo.DeleteProduct(context.Background(), product.ID)
	assert.NoError(t, err, "Failed to delete product")

	// Verify deletion by attempting to fetch it
	_, err = repo.GetProductByID(context.Background(), product.ID)
	assert.Error(t, err, "Expected an error fetching a deleted product")

}

// TestGetPackageByID tests the GetPackageByID method of the GoodsRepository.
func TestGetPackageByID(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	// create a package
	pkg := &models.Package{PackageName: "Test Package", Description: "A test package"}
	err := repo.CreatePackage(context.Background(), pkg)
	assert.NoError(t, err, "Failed to create test package")
	// add products to the package
	err = repo.AddProductToPackage(context.Background(), pkg.ID, []models.PackageContent{
		{PackageID: pkg.ID, ProductID: 1, Quantity: 10}, {PackageID: pkg.ID, ProductID: 2, Quantity: 5},
	})
	assert.NoError(t, err, "Failed to add products to test package")
	// get the package by ID
	pkg2, contents, err := repo.GetPackageByID(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to get package by ID")
	assert.NotNil(t, pkg2)
	assert.Equal(t, pkg.ID, pkg2.ID)
	assert.Equal(t, pkg.PackageName, pkg2.PackageName)
	assert.Equal(t, pkg.Description, pkg2.Description)
	assert.NotEmpty(t, contents, "No contents found in package")
	// delete the package
	err = repo.DeletePackage(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to delete test package")

}

// TestListPackages tests the ListPackages method of the GoodsRepository.
func TestListPackages(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	// get current packages
	packages, err := repo.ListPackages(context.Background())
	assert.NoError(t, err, "Failed to list packages")
	startLen := len(packages)
	// add packages to the database
	pkgs := []*models.Package{
		{PackageName: "Package 1", Description: "First package"},
		{PackageName: "Package 2", Description: "Second package"},
		{PackageName: "Package 3", Description: "Third package"},
	}
	for _, pkg := range pkgs {
		err := repo.CreatePackage(context.Background(), pkg)
		assert.NoError(t, err, "Failed to create package")
	}
	// list all packages
	packages, err = repo.ListPackages(context.Background())
	assert.NoError(t, err, "Failed to list packages")
	assert.NotNil(t, packages)
	assert.Len(t, packages, len(pkgs)+startLen)
	// delete the packages
	for _, pkg := range pkgs {
		err = repo.DeletePackage(context.Background(), pkg.ID)
		assert.NoError(t, err, "Failed to delete package")
	}
}

// TestCreatePackage tests the CreatePackage method of the GoodsRepository.
func TestCreatePackage(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	pkg := &models.Package{PackageName: "New Package", Description: "A new package for testing"}
	err := repo.CreatePackage(context.Background(), pkg)
	assert.NoError(t, err, "Failed to create package")

	// Verify creation by fetching it
	_, contents, err := repo.GetPackageByID(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to fetch the newly created package")
	assert.Empty(t, contents, "New package should have no contents")

	// Cleanup
	err = repo.DeletePackage(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to delete the newly created package")
}

// TestAddProductToPackage tests the AddProductToPackage method of the GoodsRepository.
func TestAddProductToPackage(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	// Create a package and product for testing
	pkg := &models.Package{PackageName: "Test Package", Description: "Package for product addition"}
	err := repo.CreatePackage(context.Background(), pkg)
	assert.NoError(t, err, "Failed to create package")

	product := &models.Product{Name: "Product for Package", Description: "A product to add to package", Price: 15.00, ImageURL: "images/package.jpg", SKU: "PACKAGEPROD100"}
	err = repo.CreateProduct(context.Background(), product)
	assert.NoError(t, err, "Failed to create product for package")

	// Add product to package
	err = repo.AddProductToPackage(context.Background(), pkg.ID, []models.PackageContent{
		{PackageID: pkg.ID, ProductID: product.ID, Quantity: 10},
	})
	assert.NoError(t, err, "Failed to add product to package")

	// Verify addition
	contents, err := repo.GetProductsByPackageID(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to get products by package ID")
	assert.NotEmpty(t, contents, "No products found in package")
	assert.Len(t, contents, 1, "Expected one product in package")

	// Cleanup
	err = repo.DeletePackage(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to delete package with products")
	err = repo.DeleteProduct(context.Background(), product.ID)
	assert.NoError(t, err, "Failed to delete product used in package")
}

// TestDeletePackage tests the DeletePackage method of the GoodsRepository.
func TestDeletePackage(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	// Create a package to delete
	pkg := &models.Package{PackageName: "Package to Delete", Description: "A package to be deleted"}
	err := repo.CreatePackage(context.Background(), pkg)
	assert.NoError(t, err, "Failed to create package for deletion")

	// Delete the package
	err = repo.DeletePackage(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to delete package")

	// Verify deletion by attempting to fetch it
	id, contents, err := repo.GetPackageByID(context.Background(), pkg.ID)
	assert.Error(t, err, "Expected an error fetching a deleted package")
	assert.Nil(t, id, "Expected nil package")
	assert.Nil(t, contents, "Expected nil contents for deleted package")
}
