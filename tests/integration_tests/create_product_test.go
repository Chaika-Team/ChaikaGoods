//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCreateProductSuccess(t *testing.T) {
	keyspace := "create_product_test"

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
	assert.Equal(t, id, int64(1))
}

func TestCreateProductWithSameSKU(t *testing.T) {
	keyspace := "create_product_test"

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

	productReplica := models.Product{
		Name:        "Tea",
		Price:       99.99,
		Description: "Black Knight",
		SKU:         "AE41MV9ZCHKGDA",
	}
	id, err := svc.CreateProduct(ctx, &productReplica)

	if err == nil {
		t.Fatalf("Expected err, got: %v", id)
	}
	assert.Error(t, err)
}
