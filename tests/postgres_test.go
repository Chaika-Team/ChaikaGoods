package tests

import (
	"ChaikaGoods/internal/config"
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/repository/postgresql"
	"context"
	"encoding/json"
	"fmt"
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
	cfg := config.GetConfigWithPath(logger, "..\\config_test.yml")

	dumpFile := "..\\db\\chaikagoods_test_dump.sql"
	// Create test database
	if err := createDatabase(ctx, cfg.Storage); err != nil {
		_ = logger.Log("error", fmt.Sprintf("Failed to create database: %v", err))
		return
	}

	if err := restoreDatabase(ctx, cfg.Storage, dumpFile); err != nil {
		_ = logger.Log("error", fmt.Sprintf("Failed to restore database: %v", err))
		return
	}

	defer func(ctx context.Context, storageCfg config.StorageConfig) {
		var err = deleteDatabase(ctx, storageCfg)
		if err != nil {
			_ = logger.Log("error", fmt.Sprintf("Failed to delete database: %v", err))
			panic("Failed to delete test database")
		}
	}(ctx, cfg.Storage)

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

// TestAddQueryToCreateProduct tests the AddQueryToCreateProduct method of the GoodsRepository.
func TestAddQueryToCreateProduct(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	ctx := context.Background()
	version, err := repo.GetCurrentDevVersion(ctx)
	assert.NoError(t, err, "Failed to get current dev version")
	// Clear all changes before adding a new one
	changes, err := repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	for _, change := range changes {
		err = repo.DeleteChange(ctx, change.ID)
		assert.NoError(t, err, "Failed to delete change")
	}
	// Add a new product data
	data := map[string]interface{}{
		"name":        "Test Product",
		"description": "A product for testing",
		"price":       99.99,
		"imageurl":    "images/product.jpg",
		"sku":         "TP-001",
	}
	_, err = repo.AddQueryToCreateProduct(context.Background(), &data)
	assert.NoError(t, err, "Failed to add query to create product")

	// Get all changes again
	changes, err = repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	assert.NotNil(t, changes)
	assert.Len(t, changes, 1, "Expected one new change")
	// Check the change
	change := &changes[0]
	assert.Equal(t, models.OperationTypeInsert, int(change.OperationType))
	// unmarshal the new value
	var newData map[string]interface{}
	err = json.Unmarshal(change.NewValue, &newData)
	assert.Equal(t, data, newData, "New product data should be equal to the added product data")
	// Cleanup
	err = repo.DeleteChange(ctx, change.ID)
	assert.NoError(t, err, "Failed to delete change")
}

// TestAddQueryToUpdateProduct tests the AddQueryToUpdateProduct method of the GoodsRepository.
func TestAddQueryToUpdateProduct(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())

	ctx := context.Background()
	version, err := repo.GetCurrentDevVersion(ctx)
	assert.NoError(t, err, "Failed to get current dev version")
	// Get all changes to count it before adding a new one
	changes, err := repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	startLen := len(changes)

	// get some product
	products, err := repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	product := products[0]

	// Create a updated product data
	data := map[string]interface{}{
		"id":    product.ID,
		"price": product.Price + 10.0,
	}
	// Update the product
	_, err = repo.AddQueryToUpdateProduct(ctx, &data)
	assert.NoError(t, err, "Failed to add query to update product")
	// Execute the query
	changes, err = repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	assert.NotNil(t, changes)
	assert.Len(t, changes, 1+startLen, "Expected one new change")
	// Check the last change
	change := changes[len(changes)-1]
	assert.Equal(t, models.OperationTypeUpdate, int(change.OperationType))
	// unmarshal the new value
	var updatedData map[string]interface{}
	err = json.Unmarshal(change.NewValue, &updatedData)
	assert.NoError(t, err)
	if idFloat, ok := updatedData["id"].(float64); ok {
		updatedData["id"] = int64(idFloat) // Приведение типа float64 к int64
	} else {
		assert.Fail(t, "Failed to convert id to int64")
	}

	assert.Equal(t, data, updatedData, "Price should be equal to the updated price")
	// Cleanup
	err = repo.DeleteChange(ctx, change.ID)
	assert.NoError(t, err, "Failed to delete change")
}

// TestAddQueryToDeleteProduct tests the AddQueryToDeleteProduct method of the GoodsRepository.
func TestGetAllChanges(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	ctx := context.Background()
	version, err := repo.GetCurrentDevVersion(ctx)
	// Get all changes to count it before adding a new one
	changes, err := repo.GetAllChanges(ctx, version)
	startlen := len(changes)
	assert.NoError(t, err, "Failed to get all changes")
	// Add a new change
	data := map[string]interface{}{
		"name":        "Test Product",
		"description": "A product for testing",
		"price":       99.99,
	}
	_, err = repo.AddQueryToCreateProduct(context.Background(), &data)
	assert.NoError(t, err, "Failed to add query to create product")
	// Get all changes again
	changes, err = repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	assert.NotNil(t, changes)
	assert.Len(t, changes, 1+startlen, "Expected one new change")
	// Cleanup
	err = repo.DeleteChange(ctx, changes[0].ID)
	assert.NoError(t, err, "Failed to delete change")
}

// TestGetCurrentDevVersion tests the GetCurrentDevVersion method of the GoodsRepository.
func TestGetCurrentDevVersion(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	ctx := context.Background()
	version, err := repo.GetCurrentDevVersion(ctx)
	assert.NoError(t, err, "Failed to get current dev version")
	assert.NotNil(t, version)
	assert.True(t, version.IsDev)
}

// TestApplyChanges tests the ApplyChanges method of the GoodsRepository.
func TestApplyChanges_Simple(t *testing.T) {
	// Test create, update and delete queries with applying changes
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	ctx := context.Background()
	version, err := repo.GetCurrentDevVersion(ctx)
	assert.NoError(t, err, "Failed to get current dev version")

	// Clear all changes
	changes, err := repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	for _, change := range changes {
		err = repo.DeleteChange(ctx, change.ID)
		assert.NoError(t, err, "Failed to delete change")
	}
	// Get all products
	products, err := repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	startLenProducts := len(products)

	//Create a new product
	data := map[string]interface{}{
		"name":        "Test Product",
		"description": "A product for testing",
		"price":       99.99,
		"imageurl":    "images/product.jpg",
		"sku":         "TP-001",
	}

	_, err = repo.AddQueryToCreateProduct(ctx, &data)
	assert.NoError(t, err, "Failed to add query to create product")
	// Check that the number of products not changed, because we didn't apply changes
	products, err = repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	assert.Len(t, products, startLenProducts, "Number of products should not change")
	// check that the number of changes is 1
	changes, err = repo.GetAllChanges(ctx, version)
	assert.NoError(t, err, "Failed to get all changes")
	assert.Len(t, changes, 1, "Expected one change")
	// Apply changes
	err = repo.ApplyChanges(ctx, &version)
	assert.NoError(t, err, "Failed to apply changes")
	// Get all products again
	products, err = repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	assert.NotNil(t, products)
	assert.Equal(t, startLenProducts+1, len(products), "Expected one new product")
	// Check that dev version is changed
	newVersion, err := repo.CreateNewDevVersion(ctx)
	assert.NoError(t, err, "Failed to create new dev version")
	assert.NotEqual(t, version, newVersion, "Dev version should be changed")

	//Check that the number of changes is zero, because we applied them
	changes, err = repo.GetAllChanges(ctx, newVersion)
	assert.NoError(t, err, "Failed to get all changes")
	assert.Len(t, changes, 0, "Expected no changes. All changes should be applied")

	// Get last product from products
	product := &products[len(products)-1]
	// Update the product
	updatedData := map[string]interface{}{
		"id":    product.ID,
		"price": 2077.00,
	}

	_, err = repo.AddQueryToUpdateProduct(ctx, &updatedData)
	assert.NoError(t, err, "Failed to add query to update product")
	// Check that price is not changed, because we didn't apply changes
	products, err = repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	assert.NotNil(t, products)
	assert.NotEqual(t, 2077.00, products[len(products)-1].Price, "Price should not be changed")
	// Check number of changes, should be 1
	changes, err = repo.GetAllChanges(ctx, newVersion)
	assert.NoError(t, err, "Failed to get all changes")
	assert.Len(t, changes, 1, "Expected one change")
	// Apply changes
	err = repo.ApplyChanges(ctx, &newVersion)
	assert.NoError(t, err, "Failed to apply changes")
	// Check that price is changed
	updatedProduct, err := repo.GetProductByID(ctx, product.ID)
	assert.NoError(t, err, "Failed to get product by ID")
	assert.NotNil(t, updatedProduct)
	assert.Equal(t, 2077.00, updatedProduct.Price, "Price should be updated")
	newVersion, err = repo.CreateNewDevVersion(ctx)

	// Check that the number of changes is zero, because we applied them
	changes, err = repo.GetAllChanges(ctx, newVersion)
	assert.NoError(t, err, "Failed to get all changes")
	assert.Len(t, changes, 0, "Expected no changes. All changes should be applied")

	// Delete the product
	_, err = repo.AddQueryToDeleteProduct(ctx, product.ID)
	assert.NoError(t, err, "Failed to add query to delete product")
	// Apply changes
	err = repo.ApplyChanges(ctx, &newVersion)
	assert.NoError(t, err, "Failed to apply changes")

	// Check that the number of products is the same as before
	products, err = repo.GetAllProducts(ctx)
	assert.NoError(t, err, "Failed to get all products")
	assert.Len(t, products, startLenProducts, "Number of products should be the same")

	_, err = repo.CreateNewDevVersion(ctx)
	assert.NoError(t, err, "Failed to create new dev version")
}

// TestGetAllProducts tests the GetAllProducts method of the GoodsRepository.
func TestGetAllProducts(t *testing.T) {
	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
	products, err := repo.GetAllProducts(context.Background())
	assert.NoError(t, err, "Failed to execute GetAllProducts")
	assert.NotNil(t, products)
	assert.NotEmpty(t, products)
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
	pkg2 := &models.Package{ID: pkg.ID}
	contents, err := repo.GetPackageByID(context.Background(), pkg2)
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
	pkg2 := &models.Package{ID: pkg.ID}
	contents, err := repo.GetPackageByID(context.Background(), pkg2)
	assert.NoError(t, err, "Failed to fetch the newly created package")
	assert.Empty(t, contents, "New package should have no contents")
	assert.Equal(t, pkg.PackageName, pkg2.PackageName, "Package name should match")
	assert.Equal(t, pkg.Description, pkg2.Description, "Package description should match")

	// Cleanup
	err = repo.DeletePackage(context.Background(), pkg.ID)
	assert.NoError(t, err, "Failed to delete the newly created package")
}

//// TestAddProductToPackage tests the AddProductToPackage method of the GoodsRepository.
//func TestAddProductToPackage(t *testing.T) {
//	repo := postgresql.NewGoodsRepository(dbPool, log.NewNopLogger())
//
//	// Create a package and product for testing
//	pkg := &models.Package{PackageName: "Test Package", Description: "Package for product addition"}
//	err := repo.CreatePackage(context.Background(), pkg)
//	assert.NoError(t, err, "Failed to create package")
//
//	product := &models.Product{Name: "Product for Package", Description: "A product to add to package", Price: 15.00, ImageURL: "images/package.jpg", SKU: "PACKAGEPROD100"}
//	err = repo.CreateProduct(context.Background(), product)
//	assert.NoError(t, err, "Failed to create product for package")
//
//	// Add product to package
//	err = repo.AddProductToPackage(context.Background(), pkg.ID, []models.PackageContent{
//		{PackageID: pkg.ID, ProductID: product.ID, Quantity: 10},
//	})
//	assert.NoError(t, err, "Failed to add product to package")
//
//	// Verify addition
//	contents, err := repo.GetProductsByPackageID(context.Background(), pkg.ID)
//	assert.NoError(t, err, "Failed to get products by package ID")
//	assert.NotEmpty(t, contents, "No products found in package")
//	assert.Len(t, contents, 1, "Expected one product in package")
//
//	// Cleanup
//	err = repo.DeletePackage(context.Background(), pkg.ID)
//	assert.NoError(t, err, "Failed to delete package with products")
//	err = repo.DeleteProduct(context.Background(), product.ID)
//	assert.NoError(t, err, "Failed to delete product used in package")
//}

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

	contents, err := repo.GetPackageByID(context.Background(), pkg)
	assert.Error(t, err, "Package should not be found after deletion")
	assert.Empty(t, contents, "Package should not be found after deletion")
}
