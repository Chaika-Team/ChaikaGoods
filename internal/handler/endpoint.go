package handler

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Endpoints holds all Go kit endpoints for all operations
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

// MakeEndpoints initializes all Go kit endpoints for all operations
func MakeEndpoints(logger log.Logger, service service.Service) Endpoints {
	// Инициализируем мапперы
	productMapper := &schemas.ProductMapper{}
	productsMapper := &schemas.ProductsMapper{
		ProductMapper: productMapper,
	}
	packageContentMapper := &schemas.PackageContentMapper{}
	packageMapper := &schemas.PackageMapper{
		ContentMapper: packageContentMapper,
		ProductMapper: productMapper,
	}
	packagesMapper := &schemas.PackagesMapper{
		PackageMapper: packageMapper,
	}

	return Endpoints{
		// Products
		GetAllProducts: makeGetAllProductsEndpoint(logger, service, productsMapper),
		GetProductByID: makeGetProductByIDEndpoint(logger, service, productMapper),
		// Packets
		SearchPacket:  makeSearchPacketEndpoint(logger, service, packagesMapper),
		AddPacket:     makeAddPacketEndpoint(logger, service, packageMapper),
		GetPacketByID: makeGetPacketByIDEndpoint(logger, service, packageMapper),
		// Products (admin)
		CreateProduct: makeCreateProductEndpoint(logger, service, productMapper),
		UpdateProduct: makeUpdateProductEndpoint(logger, service, productMapper),
		DeleteProduct: makeDeleteProductEndpoint(logger, service),
	}
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
func makeGetAllProductsEndpoint(logger log.Logger, s service.Service, mapper *schemas.ProductsMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(schemas.GetAllProductsRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeGetProductByIDEndpoint(logger log.Logger, s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*schemas.GetProductByIDRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeSearchPacketEndpoint(logger log.Logger, s service.Service, mapper *schemas.PackagesMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.SearchPacketRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeAddPacketEndpoint(logger log.Logger, s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.AddPacketRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeGetPacketByIDEndpoint(logger log.Logger, s service.Service, mapper *schemas.PackageMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*schemas.GetPacketByIDRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeCreateProductEndpoint(logger log.Logger, s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.CreateProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
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
//	@Router			/api/v1/products [put]
func makeUpdateProductEndpoint(logger log.Logger, s service.Service, mapper *schemas.ProductMapper) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.UpdateProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

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
func makeDeleteProductEndpoint(logger log.Logger, s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(*schemas.DeleteProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

		err := s.DeleteProduct(ctx, req.ProductID)
		if err != nil {
			return nil, err
		}

		return schemas.DeleteProductResponse{}, nil
	}
}
