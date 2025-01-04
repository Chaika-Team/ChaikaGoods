package tests

import (
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestCheckMap(t *testing.T) {
	// Valid product data
	productData := map[string]interface{}{
		"id":          1,
		"name":        "Test Product",
		"description": "A product for testing",
		"price":       99.99,
		"imageurl":    "images/product.jpg",
		"sku":         "TP-001",
	}
	err := utils.VerifyMapFields[models.Product](productData)
	assert.NoError(t, err)

	// Not full product data, but still valid
	productData = map[string]interface{}{
		"id":   1,
		"name": "Test Product",
	}
	err = utils.VerifyMapFields[models.Product](productData)
	assert.NoError(t, err)

	// Invalid product data
	productData["invalid_field"] = "invalid"
	err = utils.VerifyMapFields[models.Product](productData)
	assert.Error(t, err)
}
