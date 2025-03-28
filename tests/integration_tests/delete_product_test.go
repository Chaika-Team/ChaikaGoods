//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

func TestDeleteProduct(t *testing.T) {
	keyspace := "delete_product_test"

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

	err = svc.DeleteProduct(ctx, id)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestDeleteProductNotFound(t *testing.T) {
	keyspace := "delete_product_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Banana",
		Price:       99.99,
		Description: "From Australia",
	}

	id, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Cannot create product, got err: %v", err)
	}
	otherId := id + 1
	err = svc.DeleteProduct(ctx, otherId)
	if err == nil {
		t.Fatalf("Expected error, got nil: %v", err)
	}
}
