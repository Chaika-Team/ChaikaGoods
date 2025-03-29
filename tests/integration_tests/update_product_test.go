//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestUpdateProductSuccess(t *testing.T) {
	keyspace := "update_product_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Coffee",
		Price:       399.99,
		Description: "100 percent arabic",
	}
	id, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}

	product.Price = 499.99
	err = svc.UpdateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Unexpected err: %v", err)
	}

	productResponse, err := svc.GetProductByID(ctx, id)
	if err != nil {
		t.Fatalf("Expected product, got err: %v", err)
	}
	assert.Equal(t, productResponse.Price, product.Price)
	assert.Equal(t, productResponse.Name, product.Name)
}

func TestUpdateProductSKU(t *testing.T) {
	keyspace := "update_product_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Coffee",
		Price:       399.99,
		Description: "100 percent Arabic",
		SKU:         "AE41MV9ZCHKGDA",
	}
	id, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}

	productNew := models.Product{
		Name:        "Tea",
		Price:       99.99,
		Description: "Black Knight",
		SKU:         "AE41MV9ZCHKGDB",
	}
	idNew, err := svc.CreateProduct(ctx, &productNew)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}

	product.SKU = "AE41MV9ZCHKGDB"
	err = svc.UpdateProduct(ctx, &product)
	if err == nil {
		t.Fatalf("Expect error, got nil: %v", err)
	}

	productResponse1, err := svc.GetProductByID(ctx, id)
	if err != nil {
		t.Fatalf("Expected product, got err: %v", err)
	}

	productResponse2, err := svc.GetProductByID(ctx, idNew)
	if err != nil {
		t.Fatalf("Expected product, got err: %v", err)
	}
	assert.NotEqual(t, productResponse1.SKU, productResponse2.SKU)
}

func TestUpdateProductNotFound(t *testing.T) {
	keyspace := "update_product_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Coffee",
		Price:       399.99,
		Description: "100 percent Arabic",
		SKU:         "AE41MV9ZCHKGDA",
	}
	_, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}

	product.ID = 10
	product.Price = 299.99
	err = svc.UpdateProduct(ctx, &product)
	if err == nil {
		t.Fatalf("Expect error, got nil: %v", err)
	}
	assert.Error(t, err)
}
