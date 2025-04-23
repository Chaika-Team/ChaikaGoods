package handler

import (
	"context"
	"fmt"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

const (
	invalidRequestType = "invalid request type"
)

// Endpoints содержит все Go kit эндпоинты для всех операций
type Endpoints struct {
	// For products
	GetAllProducts    endpoint.Endpoint
	GetProductByID    endpoint.Endpoint
	GetCurrentVersion endpoint.Endpoint
	GetDelta          endpoint.Endpoint
	// For Templates
	SearchTemplates endpoint.Endpoint
	AddTemplate     endpoint.Endpoint
	GetTemplateByID endpoint.Endpoint
	// For products (admin)
	CreateProduct endpoint.Endpoint
	UpdateProduct endpoint.Endpoint
	DeleteProduct endpoint.Endpoint
}

// MakeEndpoints инициализирует все Go kit эндпоинты для всех операций
func MakeEndpoints(logger log.Logger, svc service.Service) Endpoints {
	// Инициализируем мапперы через конструкторы
	productMapper := schemas.NewProductMapper()
	productsMapper := schemas.NewProductsMapper(productMapper)
	templateContentMapper := schemas.NewTemplateContentMapper()
	templateMapper := schemas.NewTemplateMapper(templateContentMapper, productMapper)
	templatesMapper := schemas.NewTemplatesMapper(templateMapper)

	// Создаем middleware для логирования и обработки ошибок
	logMiddleware := LoggingMiddleware(logger)

	return Endpoints{
		// Products
		GetAllProducts: logMiddleware(makeGetAllProductsEndpoint(svc, productsMapper)),
		GetProductByID: logMiddleware(makeGetProductByIDEndpoint(svc, productMapper)),
		// Templates
		SearchTemplates: logMiddleware(makeSearchTemplatesEndpoint(svc, templatesMapper)),
		AddTemplate:     logMiddleware(makeAddTemplateEndpoint(svc, templateMapper)),
		GetTemplateByID: logMiddleware(makeGetTemplateByIDEndpoint(svc, templateMapper)),
		// Products (admin)
		CreateProduct: logMiddleware(makeCreateProductEndpoint(svc, productMapper)),
		UpdateProduct: logMiddleware(makeUpdateProductEndpoint(svc, productMapper)),
		DeleteProduct: logMiddleware(makeDeleteProductEndpoint(svc)),
	}
}

// castRequest приводит запрос к нужному типу или возвращает ошибку
func castRequest[T any](req interface{}) (T, error) {
	casted, ok := req.(T)
	if !ok {
		return *new(T), fmt.Errorf("invalid request type: expected %T, got %T", *new(T), req)
	}
	return casted, nil
}

// makeGetAllProductsEndpoint constructs a GetAllProducts endpoint wrapping the service.
//
//	@Summary		Get all products
//	@Description	Get all products from the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.GetAllProductsResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/product [get]
func makeGetAllProductsEndpoint(s service.Service, mapper *schemas.ProductsMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		products, err := s.GetAllProducts(ctx)
		if err != nil {
			return nil, err
		}

		productsSchema := mapper.ToSchemas(products)
		return schemas.GetAllProductsResponse{Products: productsSchema}, nil
	}
}

// makeGetProductByIDEndpoint constructs a GetProductByID endpoint wrapping the service.
//
//	@Summary		Get product by ID
//	@Description	Get product details by its ID
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	schemas.GetProductByIDResponse
//	@Failure		404	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/product/{id} [get]
func makeGetProductByIDEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.GetProductByIDRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		product, err := s.GetProductByID(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}

		productSchema := mapper.ToSchema(product)
		return schemas.GetProductByIDResponse{Product: productSchema}, nil
	}
}

// makeSearchTemplatesEndpoint constructs a SearchTemplates endpoint wrapping the service.
//
//	@Summary		Search Template
//	@Description	Search for Templates
//	@Tags			Templates
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	false	"Search query"
//	@Param			limit	query		int		true	"Limit"
//	@Param			offset	query		int		true	"Offset"
//	@Success		200		{object}	schemas.SearchTemplatesResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/product/template/search [get]
func makeSearchTemplatesEndpoint(s service.Service, mapper *schemas.TemplatesMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.SearchTemplatesRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		Templates, err := s.SearchTemplates(ctx, req.Query, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}

		TemplatesSchema := mapper.ToSchemas(Templates)
		return schemas.SearchTemplatesResponse{Templates: TemplatesSchema}, nil
	}
}

// makeAddTemplateEndpoint constructs a AddTemplate endpoint wrapping the service.
//
//	@Summary		Add Template
//	@Description	Add a new Template of products to the database
//	@Tags			Templates
//	@Accept			json
//	@Produce		json
//	@Param			Template	body		schemas.AddTemplateRequest	true	"Template details"
//	@Success		200		{object}	schemas.AddTemplateResponse
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/product/template [post]
func makeAddTemplateEndpoint(s service.Service, mapper *schemas.TemplateMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.AddTemplateRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		TemplateModel := mapper.ToModel(req.Template)

		id, err := s.AddTemplate(ctx, &TemplateModel)
		if err != nil {
			return nil, err
		}

		return schemas.AddTemplateResponse{TemplateID: id}, nil
	}
}

// makeGetTemplateByIDEndpoint constructs a GetTemplateByID endpoint wrapping the service.
//
//	@Summary		Get Template by ID
//	@Description	Get Template details by its ID
//	@Tags			Templates
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Template ID"
//	@Success		200		{object}	schemas.GetTemplateByIDResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/product/template/{id} [get]
func makeGetTemplateByIDEndpoint(s service.Service, mapper *schemas.TemplateMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.GetTemplateByIDRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		Template, err := s.GetTemplateByID(ctx, req.TemplateID)
		if err != nil {
			return nil, err
		}

		TemplateSchema := mapper.ToSchema(Template)
		return schemas.GetTemplateByIDResponse{Template: TemplateSchema}, nil
	}
}

// makeCreateProductEndpoint constructs a CreateProduct endpoint wrapping the service.
//
//	@Summary		Add product
//	@Description	Add a new product to the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		schemas.CreateProductRequest	true	"Product details"
//	@Success		200		{object}	schemas.CreateProductResponse
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/product [post]
func makeCreateProductEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.CreateProductRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		productModel := mapper.ToModel(req.Product)

		id, err := s.CreateProduct(ctx, &productModel)
		if err != nil {
			return nil, err
		}

		return schemas.CreateProductResponse{ProductID: id}, nil
	}
}

// makeUpdateProductEndpoint constructs a UpdateProduct endpoint wrapping the service.
//
//	@Summary		Update product
//	@Description	Update product details in the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			product	body		schemas.UpdateProductRequest	true	"Product details"
//	@Success		200		{object}	schemas.UpdateProductResponse
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/product [put]
func makeUpdateProductEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.UpdateProductRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		productModel := mapper.ToModel(req.Product)

		err = s.UpdateProduct(ctx, &productModel)
		if err != nil {
			return nil, err
		}

		return schemas.UpdateProductResponse{}, nil
	}
}

// makeDeleteProductEndpoint constructs a DeleteProduct endpoint wrapping the service.
//
//	@Summary		Delete product
//	@Description	Delete a product from the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Product ID"
//	@Success		200	{object}	schemas.DeleteProductResponse
//	@Failure		404	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/product/{id} [delete]
func makeDeleteProductEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := castRequest[*schemas.DeleteProductRequest](request)
		if err != nil {
			return nil, myerr.Validation(invalidRequestType, err)
		}

		err = s.DeleteProduct(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}

		return schemas.DeleteProductResponse{}, nil
	}
}
