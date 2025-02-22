package mapper

import (
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

// ProductMapper Block
func TestProductMapperToSchema(t *testing.T) {
	product := models.Product{
		ID:          1,
		Name:        "Milk",
		Description: "2.5% milk Vologodskoye",
		Price:       149.99,
		ImageURL:    "http://example.com/milk.jpg",
	}
	m := schemas.NewProductMapper()

	schema := m.ToSchema(product)

	assert.Equal(t, product.ID, schema.ID)
	assert.Equal(t, product.Name, schema.Name)
	assert.Equal(t, product.Description, schema.Description)
	assert.Equal(t, product.Price, schema.Price)
	assert.Equal(t, product.ImageURL, schema.ImageURL)
}

func TestProductMapperToModel(t *testing.T) {
	productSchema := schemas.ProductSchema{
		ID:          1,
		Name:        "Milk",
		Description: "2.5% milk Vologodskoye",
		Price:       149.99,
		ImageURL:    "http://example.com/milk.jpg",
	}
	m := schemas.NewProductMapper()

	product := m.ToModel(productSchema)

	assert.Equal(t, productSchema.ID, product.ID)
	assert.Equal(t, productSchema.Name, product.Name)
	assert.Equal(t, productSchema.Description, product.Description)
	assert.Equal(t, productSchema.Price, product.Price)
	assert.Equal(t, productSchema.ImageURL, product.ImageURL)
}

// TemplateContentMapper Block
func TestTemplateContentMapperToSchema(t *testing.T) {
	contentModel := models.TemplateContent{
		ProductID: 1,
		Quantity:  10,
	}
	tcm := schemas.NewTemplateContentMapper()

	contentSchema := tcm.ToSchema(contentModel)

	assert.Equal(t, contentModel.ProductID, contentSchema.ProductID)
	assert.Equal(t, contentModel.Quantity, contentSchema.Quantity)
}

func TestTemplateContenrMapperToModel(t *testing.T) {
	contentSchema := schemas.TemplateContentSchema{
		ProductID: 2,
		Quantity:  5,
	}
	tcm := schemas.NewTemplateContentMapper()

	contentModel := tcm.ToModel(contentSchema)

	assert.Equal(t, contentSchema.ProductID, contentModel.ProductID)
	assert.Equal(t, contentSchema.Quantity, contentModel.Quantity)
}

// TemplateMapper Block
func TestTemplateMapperToSchema(t *testing.T) {
	templateModel := models.Template{
		ID:           1,
		TemplateName: "108A",
		Description:  "Main template for route 109A St.Petersburg -> Moscow",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 10},
			{ProductID: 2, Quantity: 20},
		},
	}
	contentMapper := schemas.NewTemplateContentMapper()
	productMapper := schemas.NewProductMapper()
	tm := schemas.NewTemplateMapper(contentMapper, productMapper)

	templateSchema := tm.ToSchema(templateModel)

	assert.Equal(t, templateModel.ID, templateSchema.ID)
	assert.Equal(t, templateModel.TemplateName, templateSchema.TemplateName)
	assert.Equal(t, templateModel.Description, templateSchema.Description)
	assert.Len(t, templateSchema.Content, 2)
	assert.Equal(t, templateModel.Content[0].ProductID, templateSchema.Content[0].ProductID)
	assert.Equal(t, templateModel.Content[0].Quantity, templateSchema.Content[0].Quantity)
	assert.Equal(t, templateModel.Content[1].ProductID, templateSchema.Content[1].ProductID)
	assert.Equal(t, templateModel.Content[1].Quantity, templateSchema.Content[1].Quantity)
}

func TestTemplateMapperToModel(t *testing.T) {
	templateSchema := schemas.TemplateSchema{
		ID:           2,
		TemplateName: "109A",
		Description:  "Main template for route 109A Moscow -> St.Petersburg",
		Content: []schemas.TemplateContentSchema{
			{ProductID: 1, Quantity: 5},
			{ProductID: 2, Quantity: 7},
		},
	}
	contentMapper := schemas.NewTemplateContentMapper()
	productMapper := schemas.NewProductMapper()
	ts := schemas.NewTemplateMapper(contentMapper, productMapper)

	templateModel := ts.ToModel(templateSchema)

	assert.Equal(t, templateSchema.ID, templateModel.ID)
	assert.Equal(t, templateSchema.TemplateName, templateModel.TemplateName)
	assert.Equal(t, templateSchema.Description, templateModel.Description)
	assert.Len(t, templateSchema.Content, 2)
	assert.Equal(t, templateSchema.Content[0].ProductID, templateModel.Content[0].ProductID)
	assert.Equal(t, templateSchema.Content[0].Quantity, templateModel.Content[0].Quantity)
	assert.Equal(t, templateSchema.Content[1].ProductID, templateModel.Content[1].ProductID)
	assert.Equal(t, templateSchema.Content[1].Quantity, templateModel.Content[1].Quantity)
}

// ProductsMapper Block
func TestProductsMapperToSchemas(t *testing.T) {
	productsModel := []models.Product{
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

	schemas := pm.ToSchemas(productsModel)

	assert.Len(t, schemas, 2)
	assert.Equal(t, productsModel[0].ID, schemas[0].ID)
	assert.Equal(t, productsModel[1].Name, schemas[1].Name)
}

func TestProductsMapperToModels(t *testing.T) {
	productSchemas := []schemas.ProductSchema{
		{ID: 1, Name: "Product 1", Description: "Description 1", Price: 199.99, ImageURL: "http://example.com/product1.jpg"},
		{ID: 2, Name: "Product 2", Description: "Description 2", Price: 299.99, ImageURL: "http://example.com/product2.jpg"},
	}
	pm := schemas.NewProductsMapper(schemas.NewProductMapper())

	products := pm.ToModels(productSchemas)

	assert.Len(t, products, 2)
	assert.Equal(t, productSchemas[0].ID, products[0].ID)
	assert.Equal(t, productSchemas[1].Price, products[1].Price)
}

// TemplatesMapper Block
func TestTemplatesMapperToSchema(t *testing.T) {
	templatesModel := []models.Template{
		{
			ID:           1,
			TemplateName: "108A",
			Description:  "Main template for route 109A St.Petersburg -> Moscow",
			Content: []models.TemplateContent{
				{ProductID: 1, Quantity: 10},
				{ProductID: 2, Quantity: 20},
			},
		},
		{
			ID:           2,
			TemplateName: "109A",
			Description:  "Main template for route 109A Moscow -> St.Petersburg",
			Content: []models.TemplateContent{
				{ProductID: 1, Quantity: 5},
				{ProductID: 2, Quantity: 7},
			},
		},
	}
	contentMapper := schemas.NewTemplateContentMapper()
	productMapper := schemas.NewProductMapper()
	templateMapper := schemas.NewTemplateMapper(contentMapper, productMapper)
	tm := schemas.NewTemplatesMapper(templateMapper)

	templatesSchema := tm.ToSchemas(templatesModel)

	assert.Len(t, templatesModel, 2)
	assert.Equal(t, templatesModel[0].ID, templatesSchema[0].ID)
	assert.Equal(t, templatesModel[1].ID, templatesSchema[1].ID)
	assert.Equal(t, templatesModel[0].Content[0].ProductID, templatesSchema[0].Content[0].ProductID)
	assert.Equal(t, templatesModel[1].Content[1].Quantity, templatesSchema[1].Content[1].Quantity)
}

func TestTemplatesMapperToModel(t *testing.T) {
	templatesSchema := []schemas.TemplateSchema{
		{
			ID:           1,
			TemplateName: "108A",
			Description:  "Main template for route 109A St.Petersburg -> Moscow",
			Content: []schemas.TemplateContentSchema{
				{ProductID: 1, Quantity: 10},
				{ProductID: 2, Quantity: 20},
			},
		},
		{
			ID:           2,
			TemplateName: "109A",
			Description:  "Main template for route 109A Moscow -> St.Petersburg",
			Content: []schemas.TemplateContentSchema{
				{ProductID: 1, Quantity: 5},
				{ProductID: 2, Quantity: 7},
			},
		},
	}
	contentMapper := schemas.NewTemplateContentMapper()
	productMapper := schemas.NewProductMapper()
	templateMapper := schemas.NewTemplateMapper(contentMapper, productMapper)
	tm := schemas.NewTemplatesMapper(templateMapper)

	templatesModel := tm.ToModels(templatesSchema)

	assert.Len(t, templatesModel, 2)
	assert.Equal(t, templatesSchema[0].ID, templatesModel[0].ID)
	assert.Equal(t, templatesSchema[1].ID, templatesModel[1].ID)
	assert.Equal(t, templatesSchema[0].Content[0].ProductID, templatesModel[0].Content[0].ProductID)
	assert.Equal(t, templatesSchema[1].Content[1].Quantity, templatesModel[1].Content[1].Quantity)
}
