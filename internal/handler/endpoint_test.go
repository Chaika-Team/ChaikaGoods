package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
<<<<<<< HEAD
=======
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
>>>>>>> origin/master
	"github.com/Chaika-Team/ChaikaGoods/tests/mocks"
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
)

<<<<<<< HEAD
// Техника тест-дизайна: попарное тестирование
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию MakeEndpoints, которая инициализирует все эндпоинты.
//   - Используется мок-сервис (сгенерированный с помощью mockery) для задания ожидания вызова метода GetAllProducts.
//   - Проверяется, что все поля структуры Endpoints заполнены, а вызов одного из эндпоинтов (GetAllProducts) возвращает ожидаемый результат.
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию MakeEndpoints, которая инициализирует все эндпоинты
//   - Используется мок-сервис (сгенерированный с помощью mockery) для задания ожидания вызова метода GetAllProducts
//   - Проверяется, что все поля структуры Endpoints заполнены, а вызов одного из эндпоинтов (GetAllProducts) возвращает ожидаемый результат
>>>>>>> origin/master
func TestMakeEndpoints(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	logger := log.NewNopLogger()

	mockSvc.EXPECT().GetAllProducts(context.Background()).Return([]models.Product{}, nil)

	endpoints := MakeEndpoints(logger, mockSvc)

	assert.NotNil(t, endpoints.GetAllProducts, "GetAllProducts endpoint should not be nil")
	assert.NotNil(t, endpoints.GetProductByID, "GetProductByID endpoint should not be nil")
	assert.NotNil(t, endpoints.SearchTemplates, "SearchTemplates endpoint should not be nil")
	assert.NotNil(t, endpoints.AddTemplate, "AddTemplate endpoint should not be nil")
	assert.NotNil(t, endpoints.GetTemplateByID, "GetTemplateByID endpoint should not be nil")
	assert.NotNil(t, endpoints.CreateProduct, "CreateProduct endpoint should not be nil")
	assert.NotNil(t, endpoints.UpdateProduct, "UpdateProduct endpoint should not be nil")
	assert.NotNil(t, endpoints.DeleteProduct, "DeleteProduct endpoint should not be nil")

	// Additional test
	resp, err := endpoints.GetAllProducts(context.Background(), nil)
	assert.NoError(t, err)
	getAllResp, ok := resp.(schemas.GetAllProductsResponse)
	assert.True(t, ok, "response should be of type GetAllProductsResponse")
	assert.Empty(t, getAllResp.Products, "Products list should be empty")
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeGetAllProducts
//   - Идет проверка соответствия фактически полученных данных и тех, которые находились на эндпоинте
>>>>>>> origin/master
func TestMakeGetAllProductsEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	products := []models.Product{
		{ID: 1, Name: "Milk"},
		{ID: 2, Name: "Chocolate"},
	}
	mockSvc.EXPECT().GetAllProducts(context.Background()).Return(products, nil)

	mockProductMapper := schemas.NewProductMapper()
	mockProductsMapper := schemas.NewProductsMapper(mockProductMapper)

	ep := makeGetAllProductsEndpoint(mockSvc, mockProductsMapper)
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

// Техника тест-дизайна: Прогнозирование ошибок
<<<<<<< HEAD
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест проверяет негативный сценарий для эндпоинта GetAllProducts.
//   - Прогнозирование ошибок: если сервис возвращает ошибку, эндпоинт должен вернуть её и не формировать корректный ответ.
=======
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для эндпоинта GetAllProducts
//   - Прогнозирование ошибок: если сервис возвращает ошибку, эндпоинт должен вернуть её и не формировать корректный ответ
>>>>>>> origin/master
func TestMakeGetAllProductsEndpointFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	errMsg := "Database is unreachable!"
	mockSvc.EXPECT().GetAllProducts(context.Background()).Return(nil, errors.New(errMsg))

	mockProductMapper := schemas.NewProductMapper()
	mockProductsMapper := schemas.NewProductsMapper(mockProductMapper)

	ep := makeGetAllProductsEndpoint(mockSvc, mockProductsMapper)
	resp, err := ep(context.Background(), nil)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeGetProductByIDEndpoint
//   - Идет проверка соответствия фактически полученных данных по ID и тех, которые находились на эндпоинте ID
>>>>>>> origin/master
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

<<<<<<< HEAD
func TestMakeGetProductByIDEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeGetProductByIDEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockProductMapper := schemas.NewProductMapper()

	ep := makeGetProductByIDEndpoint(mockSvc, mockProductMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeGetProductByIDEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "ID not found"
	mockSvc.EXPECT().GetProductByID(context.Background(), int64(5)).Return(models.Product{}, errors.New(errMsg))
	mockProductMapper := schemas.NewProductMapper()

	ep := makeGetProductByIDEndpoint(mockSvc, mockProductMapper)
	req := &schemas.GetProductByIDRequest{ProductID: 5}
	resp, err := ep(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeSearchTemplatesEndpoint
//   - Идет проверка соответствия фактически полученных данных по ID и тех, которые находились на эндпоинте ID
>>>>>>> origin/master
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
	req := &schemas.SearchTemplatesRequest{Query: "Default", Limit: 10, Offset: 0}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	searchTemplates, ok := resp.(schemas.SearchTemplatesResponse)
	assert.True(t, ok, "response should be of type SearchTemplatesResponse")
	assert.Equal(t, int64(33), searchTemplates.Templates[0].ID)
	assert.Equal(t, "Default", searchTemplates.Templates[0].TemplateName)
	assert.Equal(t, "None", searchTemplates.Templates[0].Description)
	assert.Len(t, searchTemplates.Templates[0].Content, 2)
}

<<<<<<< HEAD
func TestMakeSearchTemplatesEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeSearchTemplatesEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())
	mockTemplatesMapper := schemas.NewTemplatesMapper(mockTemplateMapper)

	ep := makeSearchTemplatesEndpoint(mockSvc, mockTemplatesMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeSearchTemplatesEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "No such templates!"
	mockSvc.EXPECT().SearchTemplates(context.Background(), string("Default"), int64(10), int64(0)).Return([]models.Template{}, errors.New(errMsg))
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())
	mockTemplatesMapper := schemas.NewTemplatesMapper(mockTemplateMapper)

	ep := makeSearchTemplatesEndpoint(mockSvc, mockTemplatesMapper)
	req := &schemas.SearchTemplatesRequest{Query: "Default", Limit: 10, Offset: 0}
	resp, err := ep(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeAddTemplateEndpoint
//   - Идет проверка соответствия фактически полученного айди от добавления и айди, который должен был добавиться
>>>>>>> origin/master
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
	req := &schemas.AddTemplateRequest{Template: reqTemplateSchema}
	resp, err := ep(context.Background(), req)

	assert.NoError(t, err)
	addTempResp, ok := resp.(schemas.AddTemplateResponse)
	assert.True(t, ok, "response should be of type AddTemplateResponse")
	assert.Equal(t, int64(42), addTempResp.TemplateID)
}

<<<<<<< HEAD
func TestMakeAddTemplateEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeAddTemplateEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())

	ep := makeAddTemplateEndpoint(mockSvc, mockTemplateMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeAddTemplateEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "Template already exists!"
	reqTemplateSchema := schemas.TemplateSchema{
		ID:           55,
		TemplateName: "Varnop",
		Description:  "Exec!",
		Content: []schemas.TemplateContentSchema{
			{ProductID: 242, Quantity: 11},
			{ProductID: 1, Quantity: 58},
		},
	}
	expectedModel := models.Template{
		ID:           reqTemplateSchema.ID,
		TemplateName: reqTemplateSchema.TemplateName,
		Description:  reqTemplateSchema.Description,
		Content: []models.TemplateContent{
			{ProductID: 242, Quantity: 11},
			{ProductID: 1, Quantity: 58},
		},
	}
	mockSvc.EXPECT().AddTemplate(context.Background(), &expectedModel).Return(expectedModel.ID, errors.New(errMsg))
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())

	ep := makeAddTemplateEndpoint(mockSvc, mockTemplateMapper)
	req := &schemas.AddTemplateRequest{Template: reqTemplateSchema}
	resp, err := ep(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeGetTemplateByIDEndpoint
//   - Идет проверка соответствия фактически полученного шаблона и того, который должен был быть получен
>>>>>>> origin/master
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

<<<<<<< HEAD
func TestMakeGetTemplateByIDEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeGetTemplateByIDEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockTemplateMapper := schemas.NewTemplateMapper(schemas.NewTemplateContentMapper(), schemas.NewProductMapper())

	ep := makeGetTemplateByIDEndpoint(mockSvc, mockTemplateMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeGetTemplateByIDEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "ID not found"
	mockSvc.EXPECT().GetTemplateByID(context.Background(), int64(51)).Return(models.Template{}, errors.New(errMsg))
	mockProductMapper := schemas.NewProductMapper()
	mockContentMapper := schemas.NewTemplateContentMapper()
	mockTemplateMapper := schemas.NewTemplateMapper(mockContentMapper, mockProductMapper)

	ep := makeGetTemplateByIDEndpoint(mockSvc, mockTemplateMapper)
	req := &schemas.GetTemplateByIDRequest{TemplateID: 51}
	resp, err := ep(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeCreateProductEndpoint
//   - Идет проверка соответствия фактически полученного продукта и того, который должен был быть получен
>>>>>>> origin/master
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
	req := &schemas.CreateProductRequest{Product: reqProductSchema}
	resp, err := ep(context.Background(), req)
	assert.NoError(t, err)

	cpResp, ok := resp.(schemas.CreateProductResponse)
	assert.True(t, ok, "response should be of type CreateProductResponse")
	assert.Equal(t, int64(1), cpResp.ProductID)
}

<<<<<<< HEAD
func TestMakeCreateProductEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeCreateProductEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockProductMapper := schemas.NewProductMapper()

	ep := makeCreateProductEndpoint(mockSvc, mockProductMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeCreateProductEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "Product already exists!"
	reqProductSchema := schemas.ProductSchema{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       89.99,
		ImageURL:    "http://example.com/updatedimage.png",
	}
	expectedModel := models.Product{
		ID:          reqProductSchema.ID,
		Name:        reqProductSchema.Name,
		Description: reqProductSchema.Description,
		Price:       reqProductSchema.Price,
		ImageURL:    reqProductSchema.ImageURL,
	}
	mockProductMapper := schemas.NewProductMapper()
	mockSvc.EXPECT().CreateProduct(context.Background(), &expectedModel).Return(expectedModel.ID, errors.New(errMsg))

	ep := makeCreateProductEndpoint(mockSvc, mockProductMapper)
	req := &schemas.CreateProductRequest{Product: reqProductSchema}
	resp, err := ep(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeUpdateProductEndpoint
//   - Идет проверка соответствия фактически полученного ответа об обновлении и того, который должен был быть получен
>>>>>>> origin/master
func TestMakeUpdateProductEndpointSuccess(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockProductMapper := schemas.NewProductMapper()
	reqProductSchema := schemas.ProductSchema{
		ID:   3,
		Name: "Updated chocolate",
	}
	reqUpdate := &schemas.UpdateProductRequest{Product: reqProductSchema}

	expectedModel := models.Product{
		ID:   3,
		Name: "Updated chocolate",
	}
	mockSvc.EXPECT().UpdateProduct(context.Background(), &expectedModel).Return(nil)

	ep := makeUpdateProductEndpoint(mockSvc, mockProductMapper)
	resp, err := ep(context.Background(), reqUpdate)
	assert.NoError(t, err)
	_, ok := resp.(schemas.UpdateProductResponse)
	assert.True(t, ok, "response should be of type UpdateProductResponse")
}

<<<<<<< HEAD
func TestMakeUpdateProductEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeUpdateProductEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)
	mockProductMapper := schemas.NewProductMapper()

	ep := makeUpdateProductEndpoint(mockSvc, mockProductMapper)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeUpdateProductEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "Product to update not found!"
	reqProductSchema := schemas.ProductSchema{
		ID:   5,
		Name: "Updated milk 3.2%",
	}
	expectedModel := models.Product{
		ID:   5,
		Name: "Updated milk 3.2%",
	}
	mockSvc.EXPECT().UpdateProduct(context.Background(), &expectedModel).Return(errors.New(errMsg))
	mockProductMapper := schemas.NewProductMapper()
	reqUpdate := &schemas.UpdateProductRequest{Product: reqProductSchema}

	ep := makeUpdateProductEndpoint(mockSvc, mockProductMapper)
	resp, err := ep(context.Background(), reqUpdate)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}

<<<<<<< HEAD
=======
// Техника тест-дизайна: Классы эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет функцию makeUpdateProductEndpoint
//   - Идет проверка соответствия фактически полученного ответа об удалении и того, который должен был быть получен
>>>>>>> origin/master
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

<<<<<<< HEAD
func TestMakeDeleteProductEndpointFailed(t *testing.T) {
=======
// Техника тест-дизайна: Прогнозирование ошибок + таблица эквивалентности
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для неверного типа request
//   - Прогнозирование ошибок: если запрос возвращает ошибку, то он должен вернуть ошибку валидации
func TestMakeDeleteProductEndpointCastFailed(t *testing.T) {
	mockSvc := mocks.NewMockService(t)

	ep := makeDeleteProductEndpoint(mockSvc)
	req := "invalid request type"
	resp, err := ep(context.Background(), req)
	assert.Error(t, err)
	assert.Nil(t, resp)

	appErr, ok := myerr.IsAppError(err)
	assert.True(t, ok, "error is AppError")
	assert.Contains(t, appErr.Message, req)
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Фоломкин Дмитрий
// Описание:
//   - Тест проверяет негативный сценарий для ошибок на уровне сервиса во время запроса
//   - Прогнозирование ошибок: если сервис возвращает ошибку, то запрос неудачен и сервис недоступен
func TestMakeDeleteProductEndpointServiceFailed(t *testing.T) {
>>>>>>> origin/master
	mockSvc := mocks.NewMockService(t)
	errMsg := "Product to delete not found!"
	reqDelete := &schemas.DeleteProductRequest{
		ProductID: 19,
	}
	mockSvc.EXPECT().DeleteProduct(context.Background(), int64(19)).Return(errors.New(errMsg))

	ep := makeDeleteProductEndpoint(mockSvc)
	resp, err := ep(context.Background(), reqDelete)

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, errMsg, err.Error())
}
