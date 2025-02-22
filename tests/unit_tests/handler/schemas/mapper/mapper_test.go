package mapper

import (
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestProductMapperToSchema(t *testing.T) {
	// Arrange
	product := models.Product{
		ID:          1,
		Name:        "Test Product",
		Description: "This is a test product.",
		Price:       199.99,
		ImageURL:    "http://example.com/product.jpg",
	}
	m := schemas.NewProductMapper()

	// Act
	schema := m.ToSchema(product)

	// Assert
	assert.Equal(t, product.ID, schema.ID)
	assert.Equal(t, product.Name, schema.Name)
	assert.Equal(t, product.Description, schema.Description)
	assert.Equal(t, product.Price, schema.Price)
	assert.Equal(t, product.ImageURL, schema.ImageURL)
}

func TestProductMapperToModel(t *testing.T) {
	// Arrange
	productSchema := schemas.ProductSchema{
		ID:          1,
		Name:        "Test Product",
		Description: "This is a test product.",
		Price:       199.99,
		ImageURL:    "http://example.com/product.jpg",
	}
	m := schemas.NewProductMapper()

	// Act
	product := m.ToModel(productSchema)

	// Assert
	assert.Equal(t, productSchema.ID, product.ID)
	assert.Equal(t, productSchema.Name, product.Name)
	assert.Equal(t, productSchema.Description, product.Description)
	assert.Equal(t, productSchema.Price, product.Price)
	assert.Equal(t, productSchema.ImageURL, product.ImageURL)
}

func TestProductsMapperToSchemas(t *testing.T) {
	// Arrange
	products := []models.Product{
		{
			ID:          1,
			Name:        "Product 1",
			Description: "Description 1",
			Price:       199.99,
			ImageURL:    "http://example.com/product1.jpg",
		},
		{
			ID:          2,
			Name:        "Product 2",
			Description: "Description 2",
			Price:       299.99,
			ImageURL:    "http://example.com/product2.jpg",
		},
	}
	pm := schemas.NewProductsMapper(schemas.NewProductMapper())

	// Act
	schemas := pm.ToSchemas(products)

	// Assert
	assert.Len(t, schemas, 2)
	assert.Equal(t, products[0].ID, schemas[0].ID)
	assert.Equal(t, products[1].Name, schemas[1].Name)
}

func TestProductsMapperToModels(t *testing.T) {
	// Arrange
	productSchemas := []schemas.ProductSchema{
		{ID: 1, Name: "Product 1", Description: "Description 1", Price: 199.99, ImageURL: "http://example.com/product1.jpg"},
		{ID: 2, Name: "Product 2", Description: "Description 2", Price: 299.99, ImageURL: "http://example.com/product2.jpg"},
	}
	pm := schemas.NewProductsMapper(schemas.NewProductMapper())

	// Act
	products := pm.ToModels(productSchemas)

	// Assert
	assert.Len(t, products, 2)
	assert.Equal(t, productSchemas[0].ID, products[0].ID)
	assert.Equal(t, productSchemas[1].Price, products[1].Price)
}

func TestTemplateMapperToSchema(t *testing.T) {
	// Arrange
	template := models.Template{
		ID:           1,
		TemplateName: "Template 1",
		Description:  "Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 10},
			{ProductID: 2, Quantity: 20},
		},
	}
	contentMapper := schemas.NewTemplateContentMapper()
	productMapper := schemas.NewProductMapper()
	tm := schemas.NewTemplateMapper(contentMapper, productMapper)

	// Act
	templateSchema := tm.ToSchema(template)

	// Assert
	assert.Equal(t, template.ID, templateSchema.ID)
	assert.Equal(t, template.TemplateName, templateSchema.TemplateName)
	assert.Len(t, templateSchema.Content, 2)
	assert.Equal(t, template.Content[0].ProductID, templateSchema.Content[0].ProductID)
}

func TestTemplateContentMapperToSchema(t *testing.T) {
	// Arrange
	content := models.TemplateContent{
		ProductID: 1,
		Quantity:  10,
	}
	m := schemas.NewTemplateContentMapper()

	// Act
	contentSchema := m.ToSchema(content)

	// Assert
	assert.Equal(t, content.ProductID, contentSchema.ProductID)
	assert.Equal(t, content.Quantity, contentSchema.Quantity)
}
