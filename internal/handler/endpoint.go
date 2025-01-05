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
	// For packets
	SearchPacket  endpoint.Endpoint
	AddPacket     endpoint.Endpoint
	GetPacketByID endpoint.Endpoint
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
		// Packets
		SearchPacket:  logMiddleware(makeSearchPacketEndpoint(svc, packagesMapper)),
		AddPacket:     logMiddleware(makeAddPacketEndpoint(svc, packageMapper)),
		GetPacketByID: logMiddleware(makeGetPacketByIDEndpoint(svc, packageMapper)),
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

// makeSearchPacketEndpoint constructs a SearchPacket endpoint wrapping the service.
//
//	@Summary		Search packet
//	@Description	Search for packets
//	@Tags			packets
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	false	"Search query"
//	@Param			limit	query		int		true	"Limit"
//	@Param			offset	query		int		true	"Offset"
//	@Success		200		{object}	schemas.SearchPacketResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/packets/search [get]
func makeSearchPacketEndpoint(s service.Service, mapper *schemas.PackagesMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.SearchPacketRequest](request)

		packets, err := s.SearchPacket(ctx, req.Query, req.Limit, req.Offset)
		if err != nil {
			return nil, err
		}

		packetsSchema := mapper.ToSchemas(packets)
		return schemas.SearchPacketResponse{Packets: packetsSchema}, nil
	}
}

// makeAddPacketEndpoint constructs a AddPacket endpoint wrapping the service.
//
//	@Summary		Add packet
//	@Description	Add a new packet of products to the database
//	@Tags			packets
//	@Accept			json
//	@Produce		json
//	@Param			packet	body		schemas.AddPacketRequest	true	"Packet details"
//	@Success		200		{object}	schemas.AddPacketResponse
//	@Failure		400		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/packets [post]
func makeAddPacketEndpoint(s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[schemas.AddPacketRequest](request)

		packetModel := mapper.ToModel(req.Packet)

		id, err := s.AddPacket(ctx, &packetModel)
		if err != nil {
			return nil, err
		}

		return schemas.AddPacketResponse{PacketID: id}, nil
	}
}

// makeGetPacketByIDEndpoint constructs a GetPacketByID endpoint wrapping the service.
//
//	@Summary		Get packet by ID
//	@Description	Get packet details by its ID
//	@Tags			packets
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Packet ID"
//	@Success		200		{object}	schemas.GetPacketByIDResponse
//	@Failure		404		{object}	schemas.ErrorResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/packets/{id} [get]
func makeGetPacketByIDEndpoint(s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := mustCast[*schemas.GetPacketByIDRequest](request)

		packet, err := s.GetPacketByID(ctx, req.PacketID)
		if err != nil {
			return nil, err
		}

		packetSchema := mapper.ToSchema(packet)
		return schemas.GetPacketByIDResponse{Packet: packetSchema}, nil
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
