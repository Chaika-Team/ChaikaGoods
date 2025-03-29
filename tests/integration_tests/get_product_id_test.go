//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetProductID(t *testing.T) {
	keyspace := "get_product_id_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Milk",
		Price:       99.99,
		Description: "Basic 2.5% milk",
	}

	id, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Cannot create product, got err: %v", err)
	}

	productInfo, err := svc.GetProductByID(ctx, id)
	if err != nil {
		t.Fatalf("Expected product, got err: %v", err)
	}
	assert.Equal(t, productInfo.ID, id)
	assert.Equal(t, productInfo.Price, product.Price)
	assert.Equal(t, productInfo.Name, product.Name)

	productInfo, err = svc.GetProductByID(ctx, 3)
	if err == nil {
		t.Fatalf("Expected error, got product info: %v", productInfo)
	}
	assert.Equal(t, productInfo.ID, int64(0))
	assert.Contains(t, err.Error(), "not found")
}
