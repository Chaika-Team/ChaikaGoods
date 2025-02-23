package handler

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/tests/mocks"
	"github.com/stretchr/testify/assert"
)

// TestMakeGetAllProductsEndpointSuccess проверяет эндпоинт GetAllProducts.
func TestMakeGetAllProductsEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	products := []models.Product{
		{ID: 1, Name: "Milk"},
		{ID: 2, Name: "Chocolate"},
	}
	mockSvc.EXPECT().GetAllProducts(context.Background()).Return(products, nil)

	mockProductMapper := schemas.NewProductMapper()

	productsMapper := schemas.NewProductsMapper(mockProductMapper)
	ep := makeGetAllProductsEndpoint(mockSvc, productsMapper)

	resp, err := ep(context.Background(), nil)
	assert.NoError(t, err)
	getAllResp, ok := resp.(schemas.GetAllProductsResponse)
	assert.True(t, ok, "response should be of type GetAllProductsResponse")
	assert.Len(t, getAllResp.Products, 2)
	assert.Equal(t, int64(1), getAllResp.Products[0].ID)
	assert.Equal(t, "Milk", getAllResp.Products[0].Name)
	assert.Equal(t, int64(2), getAllResp.Products[1].ID)
	assert.Equal(t, "Chocolate", getAllResp.Products[1].Name)
}

func TestMakeGetProductByIDEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	product := models.Product{ID: 33, Name: "Doshirak"}
	mockSvc.EXPECT().GetProductByID(context.Background(), int64(33)).Return(product, nil)

	mockProductMapper := schemas.NewProductMapper()

	ep := makeGetProductByIDEndpoint(mockSvc, mockProductMapper)
	req := &schemas.GetProductByIDRequest{ProductID: 33}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	getByIDResp, ok := resp.(schemas.GetProductByIDResponse)
	assert.True(t, ok, "response should be of type GetProductByIDResponse")
	assert.Equal(t, int64(33), getByIDResp.Product.ID)
	assert.Equal(t, "Doshirak", getByIDResp.Product.Name)
}

func TestMakeSearchTemplatesEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	templates := []models.Template{
		{
			ID:           33,
			TemplateName: "Default",
			Description:  "None",
			Content: []models.TemplateContent{
				{ProductID: 6, Quantity: 50},
				{ProductID: 5, Quantity: 10},
			},
		},
	}
	mockSvc.EXPECT().SearchTemplates(context.Background(), string("Default"), int64(10), int64(0)).Return(templates, nil)
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())
	mockTemplatesMapper := schemas.NewTemplatesMapper(mockTemplateMapper)

	ep := makeSearchTemplatesEndpoint(mockSvc, mockTemplatesMapper)
	req := schemas.SearchTemplatesRequest{Query: "Default", Limit: 10, Offset: 0}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	searchTemplates, ok := resp.(schemas.SearchTemplatesResponse)
	assert.True(t, ok, "response should be of type SearchTemplatesResponse")
	assert.Equal(t, int64(33), searchTemplates.Templates[0].ID)
	assert.Equal(t, "Default", searchTemplates.Templates[0].TemplateName)
	assert.Equal(t, "None", searchTemplates.Templates[0].Description)
	assert.Len(t, searchTemplates.Templates[0].Content, 2)
}

func TestMakeAddTemplateEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	reqTemplateSchema := schemas.TemplateSchema{
		ID:           42,
		TemplateName: "Premium",
		Description:  "Test only!",
		Content: []schemas.TemplateContentSchema{
			{ProductID: 22, Quantity: 66},
			{ProductID: 3, Quantity: 52},
		},
	}
	mockProductMapper := schemas.NewProductMapper()
	mockContentMapper := schemas.NewTemplateContentMapper()
	mockTemplateMapper := schemas.NewTemplateMapper(mockContentMapper, mockProductMapper)
	expectedModel := models.Template{
		ID:           reqTemplateSchema.ID,
		TemplateName: reqTemplateSchema.TemplateName,
		Description:  reqTemplateSchema.Description,
		Content: []models.TemplateContent{
			{ProductID: 22, Quantity: 66},
			{ProductID: 3, Quantity: 52},
		},
	}
	mockSvc.EXPECT().AddTemplate(context.Background(), &expectedModel).Return(reqTemplateSchema.ID, nil)

	ep := makeAddTemplateEndpoint(mockSvc, mockTemplateMapper)
	req := schemas.AddTemplateRequest{Template: reqTemplateSchema}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	addTempResp, ok := resp.(schemas.AddTemplateResponse)
	assert.True(t, ok, "response should be of type AddTemplateResponse")
	assert.Equal(t, int64(42), addTempResp.TemplateID)
}

func TestMakeGetTemplateByIDEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	template := models.Template{
		ID:           33,
		TemplateName: "Default",
		Description:  "None",
		Content: []models.TemplateContent{
			{ProductID: 6, Quantity: 50},
			{ProductID: 5, Quantity: 10},
		},
	}
	mockSvc.EXPECT().GetTemplateByID(context.Background(), int64(33)).Return(template, nil)
	mockProductMapper := schemas.NewProductMapper()
	mockContentMapper := schemas.NewTemplateContentMapper()
	mockTemplateMapper := schemas.NewTemplateMapper(mockContentMapper, mockProductMapper)

	ep := makeGetTemplateByIDEndpoint(mockSvc, mockTemplateMapper)
	req := &schemas.GetTemplateByIDRequest{TemplateID: 33}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	getByIDResp, ok := resp.(schemas.GetTemplateByIDResponse)
	assert.True(t, ok, "response should be of type GetTemplateByIDResponse")
	assert.Equal(t, int64(33), getByIDResp.Template.ID)
	assert.Equal(t, "Default", getByIDResp.Template.TemplateName)
	assert.Equal(t, "None", getByIDResp.Template.Description)
	assert.Len(t, getByIDResp.Template.Content, 2)
}

// TestMakeCreateProductEndpoint_Success проверяет эндпоинт CreateProduct.
func TestMakeCreateProductEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	reqProductSchema := schemas.ProductSchema{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       89.99,
		ImageURL:    "http://example.com/updatedimage.png",
	}
	mockProductMapper := schemas.NewProductMapper()
	expectedModel := models.Product{
		ID:          reqProductSchema.ID,
		Name:        reqProductSchema.Name,
		Description: reqProductSchema.Description,
		Price:       reqProductSchema.Price,
		ImageURL:    reqProductSchema.ImageURL,
	}
	mockSvc.EXPECT().CreateProduct(context.Background(), &expectedModel).
		Return(reqProductSchema.ID, nil)

	ep := makeCreateProductEndpoint(mockSvc, mockProductMapper)
	req := schemas.CreateProductRequest{Product: reqProductSchema}
	resp, err := ep(context.Background(), req)
	assert.NoError(t, err)

	cpResp, ok := resp.(schemas.CreateProductResponse)
	assert.True(t, ok, "response should be of type CreateProductResponse")
	assert.Equal(t, int64(1), cpResp.ProductID)
}

func TestMakeUpdateProductEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockProductMapper := schemas.NewProductMapper()
	reqProductSchema := schemas.ProductSchema{
		ID:   5,
		Name: "Updated milk 3.2%",
	}
	reqUpdate := schemas.UpdateProductRequest{Product: reqProductSchema}

	expectedModel := models.Product{
		ID:   5,
		Name: "Updated milk 3.2%",
	}
	mockSvc.EXPECT().UpdateProduct(context.Background(), &expectedModel).Return(nil)

	ep := makeUpdateProductEndpoint(mockSvc, mockProductMapper)
	resp, err := ep(context.Background(), reqUpdate)
	assert.NoError(t, err)
	_, ok := resp.(schemas.UpdateProductResponse)
	assert.True(t, ok, "response should be of type UpdateProductResponse")
}

func TestMakeDeleteProductEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	reqDelete := &schemas.DeleteProductRequest{
		ProductID: 10,
	}
	mockSvc.EXPECT().DeleteProduct(context.Background(), int64(10)).Return(nil)

	ep := makeDeleteProductEndpoint(mockSvc)
	resp, err := ep(context.Background(), reqDelete)

	assert.NoError(t, err)
	_, ok := resp.(schemas.DeleteProductResponse)
	assert.True(t, ok, "response should be of type DeleteProductResponse")
}

// Техника тест-дизайна: Попарное тестирование и таблица принятия решений
// Автор: Джон Дое
// Описание:
//   - Тесты эндпоинтов используют сгенерированные моки для Service и Mapper,
//     что позволяет изолировать логику эндпоинта от реальной реализации зависимостей.
//   - Для каждого эндпоинта проверяются успешные сценарии и корректное преобразование входных данных.
