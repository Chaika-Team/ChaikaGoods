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

// createTestTemplate создает и возвращает тестовый пакет с заданным ID и именем.
func createTestTemplate(id int64, name string) models.Template {
	return models.Template{
		ID:           id,
		TemplateName: name,
		Description:  "Description for " + name,
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
}
