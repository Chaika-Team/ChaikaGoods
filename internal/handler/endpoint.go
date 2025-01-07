package handler

import (
	"context"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

// Endpoints содержит все Go kit эндпоинты для всех операций
type Endpoints struct {
	// For products
	GetAllProducts    endpoint.Endpoint
	GetProductByID    endpoint.Endpoint
	GetCurrentVersion endpoint.Endpoint
	GetDelta          endpoint.Endpoint
	// For Packages
	SearchPackages endpoint.Endpoint
	AddPackage     endpoint.Endpoint
	GetPackageByID endpoint.Endpoint
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
	packageContentMapper := schemas.NewPackageContentMapper()
	packageMapper := schemas.NewPackageMapper(packageContentMapper, productMapper)
	packagesMapper := schemas.NewPackagesMapper(packageMapper)

	// Создаем middleware для логирования и обработки ошибок
	logMiddleware := LoggingMiddleware(logger)

	return Endpoints{
		// Products
		GetAllProducts: logMiddleware(makeGetAllProductsEndpoint(svc, productsMapper)),
		GetProductByID: logMiddleware(makeGetProductByIDEndpoint(svc, productMapper)),
		// Packages
		SearchPackages: logMiddleware(makeSearchPackagesEndpoint(svc, packagesMapper)),
		AddPackage:     logMiddleware(makeAddPackageEndpoint(svc, packageMapper)),
		GetPackageByID: logMiddleware(makeGetPackageByIDEndpoint(svc, packageMapper)),
		// Products (admin)
		CreateProduct: logMiddleware(makeCreateProductEndpoint(svc, productMapper)),
		UpdateProduct: logMiddleware(makeUpdateProductEndpoint(svc, productMapper)),
		DeleteProduct: logMiddleware(makeDeleteProductEndpoint(svc)),
	}
}

// Типизированная функция для приведения типов без проверки
func mustCast[T any](req interface{}) T {
	return req.(T) // Паникует, если приведение неудачно
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
//	@Router			/api/v1/products [get]
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
//	@Router			/api/v1/products/{id} [get]
func makeGetProductByIDEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[*schemas.GetProductByIDRequest](request)

		product, err := s.GetProductByID(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}

		productSchema := mapper.ToSchema(product)
		return schemas.GetProductByIDResponse{Product: productSchema}, nil
	}
}

// makeSearchPackagesEndpoint constructs a SearchPackages endpoint wrapping the service.
//
//	@Summary		Search Package
//	@Description	Search for Packages
//	@Tags			Packages
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	false	"Search query"
//	@Param			limit	query		int		true	"Limit"
//	@Param			offset	query		int		true	"Offset"
//	@Success		200		{object}	schemas.SearchPackagesResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/Packages/search [get]
func makeSearchPackagesEndpoint(s service.Service, mapper *schemas.PackagesMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.SearchPackagesRequest](request)

		Packages, err := s.SearchPackages(ctx, req.Query, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}

		PackagesSchema := mapper.ToSchemas(Packages)
		return schemas.SearchPackagesResponse{Packages: PackagesSchema}, nil
	}
}

// makeAddPackageEndpoint constructs a AddPackage endpoint wrapping the service.
//
//	@Summary		Add Package
//	@Description	Add a new Package of products to the database
//	@Tags			Packages
//	@Accept			json
//	@Produce		json
//	@Param			Package	body		schemas.AddPackageRequest	true	"Package details"
//	@Success		200		{object}	schemas.AddPackageResponse
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/Packages [post]
func makeAddPackageEndpoint(s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.AddPackageRequest](request)

		PackageModel := mapper.ToModel(req.Package)

		id, err := s.AddPackage(ctx, &PackageModel)
		if err != nil {
			return nil, err
		}

		return schemas.AddPackageResponse{PackageID: id}, nil
	}
}

// makeGetPackageByIDEndpoint constructs a GetPackageByID endpoint wrapping the service.
//
//	@Summary		Get Package by ID
//	@Description	Get Package details by its ID
//	@Tags			Packages
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Package ID"
//	@Success		200		{object}	schemas.GetPackageByIDResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/Packages/{id} [get]
func makeGetPackageByIDEndpoint(s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[*schemas.GetPackageByIDRequest](request)

		Package, err := s.GetPackageByID(ctx, req.PackageID)
		if err != nil {
			return nil, err
		}

		PackageSchema := mapper.ToSchema(Package)
		return schemas.GetPackageByIDResponse{Package: PackageSchema}, nil
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
//	@Router			/api/v1/products [post]
func makeCreateProductEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.CreateProductRequest](request)

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
//	@Router			/api/v1/products [put]
func makeUpdateProductEndpoint(s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.UpdateProductRequest](request)

		productModel := mapper.ToModel(req.Product)

		err := s.UpdateProduct(ctx, &productModel)
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
//	@Router			/api/v1/products/{id} [delete]
func makeDeleteProductEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[*schemas.DeleteProductRequest](request)

		err := s.DeleteProduct(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}

		return schemas.DeleteProductResponse{}, nil
	}
}
