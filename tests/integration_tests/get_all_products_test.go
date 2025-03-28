//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestGetAllProductSuccess(t *testing.T) {
	keyspace := "get_all_products_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	productOne := models.Product{
		Name:        "Milk",
		Price:       99.99,
		Description: "Basic 2.5% milk",
		SKU:         "SF93N30A",
	}

	productTwo := models.Product{
		Name:        "Chocolate",
		Price:       199.99,
		Description: "80 percent cacao",
		SKU:         "SF93N30B",
	}

	id, err := svc.CreateProduct(ctx, &productOne)
	if err != nil {
		t.Fatalf("Cannot create product one, got err: %v", err)
	}
	assert.Equal(t, id, int64(1))

	id, err = svc.CreateProduct(ctx, &productTwo)
	if err != nil {
		t.Fatalf("Cannot create product two, got err: %v", err)
	}
	assert.Equal(t, id, int64(2))

	productsInfo, err := svc.GetAllProducts(ctx)
	if err != nil {
		t.Fatalf("Cannot get all products, got err: %v", err)
	}
	assert.Len(t, productsInfo, 2)
}

func TestGetAllProductEmpty(t *testing.T) {
	keyspace := "get_all_products_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	productsInfo, err := svc.GetAllProducts(ctx)
	if err != nil {
		t.Fatalf("Cannot get all products, got err: %v", err)
	}
	assert.Len(t, productsInfo, 0)
}
