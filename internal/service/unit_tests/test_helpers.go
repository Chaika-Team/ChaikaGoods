package unit_tests

import (
	"strconv"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

// createTestProduct создает и возвращает тестовый продукт с заданным ID и именем.
func createTestProduct(id int64, name string) models.Product {
	return models.Product{
		ID:          id,
		Name:        name,
		Description: "Description for " + name,
		Price:       99.99,
		ImageURL:    "http://example.com/" + name + ".png",
		SKU:         "SKU" + strconv.FormatInt(id, 10),
	}
}

// createTestPackage создает и возвращает тестовый пакет с заданным ID и именем.
func createTestPackage(id int64, name string) models.Package {
	return models.Package{
		ID:          id,
		PackageName: name,
		Description: "Description for " + name,
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
}
