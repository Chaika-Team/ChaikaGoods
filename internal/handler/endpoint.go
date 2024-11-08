package handler

import (
	"ChaikaGoods/internal/handler/schemas"
	"ChaikaGoods/internal/service"
	"context"
	"errors"
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
	SearchPacket endpoint.Endpoint
	AddPacket    endpoint.Endpoint
	// For products (admin)
	AddProduct    endpoint.Endpoint
	UpdateProduct endpoint.Endpoint
	DeleteProduct endpoint.Endpoint
}

// MakeEndpoints initializes all Go kit endpoints for all operations
func MakeEndpoints(logger log.Logger, service service.GoodsService) Endpoints {
	return Endpoints{
		// Products
		GetAllProducts:    makeGetAllProductsEndpoint(logger, service),
		GetProductByID:    makeGetProductByIDEndpoint(logger, service),
		GetCurrentVersion: makeGetCurrentVersionEndpoint(logger, service),
		GetDelta:          makeGetDeltaEndpoint(logger, service),
		// Packets
		SearchPacket: makeSearchPacketEndpoint(logger, service),
		AddPacket:    makeAddPacketEndpoint(logger, service),
		// Products (admin)
		AddProduct:    makeAddProductEndpoint(logger, service),
		UpdateProduct: makeUpdateProductEndpoint(logger, service),
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
func makeGetAllProductsEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(schemas.GetAllProductsRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		products, err := s.GetAllProducts(ctx)
		productsSchema := make([]schemas.ProductSchema, 0)
		for _, product := range products {
			p := schemas.ProductSchema{
				ID:          product.ID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				ImageURL:    product.ImageURL,
			}
			productsSchema = append(productsSchema, p)
		}
		return schemas.GetAllProductsResponse{Products: productsSchema}, err

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
func makeGetProductByIDEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.GetProductByIDRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		product, err := s.GetProductByID(ctx, req.ProductID)
		// Convert product model to product schema
		productSchema := schemas.ProductSchema{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			ImageURL:    product.ImageURL,
		}
		return schemas.GetProductByIDResponse{Product: productSchema}, err
	}

}

// makeGetCurrentVersionEndpoint constructs a GetCurrentVersion endpoint wrapping the service.
//
//	@Summary		Get current version
//	@Description	Get the current version of the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.GetCurrentVersionResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/version [get]
func makeGetCurrentVersionEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(schemas.GetCurrentVersionRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		version, err := s.GetCurrentVersion(ctx)
		return schemas.GetCurrentVersionResponse{Version: version}, err
	}
}

// makeGetDeltaEndpoint constructs a GetDelta endpoint wrapping the service.
//
//	@Summary		Get delta
//	@Description	Get the changes in the database between versions
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Param			version_id	query		int	true	"Version ID"
//	@Success		200			{object}	schemas.GetDeltaResponse
//	@Failure		404			{object}	schemas.ErrorResponse
//	@Failure		500			{object}	schemas.ErrorResponse
//	@Router			/api/v1/delta [get]
func makeGetDeltaEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.GetDeltaRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		changes, err := s.GetDelta(ctx, req.VersionID)
		var changeSchemas []schemas.ChangeSchema
		// Convert change model to change schema
		for _, change := range changes {
			c := schemas.ChangeSchema{
				ID:            change.ID,
				VersionID:     change.VersionID,
				OperationType: int(change.OperationType),
				// Convert json raw data to string
				NewValue:        string(change.NewValue),
				ChangeTimestamp: change.ChangeTimestamp,
				Considered:      change.Considered,
			}
			changeSchemas = append(changeSchemas, c)
		}
		return schemas.GetDeltaResponse{Changes: changeSchemas}, err
	}
}

// makeSearchPacketEndpoint constructs a SearchPacket endpoint wrapping the service.
//
//	@Summary		Search packet
//	@Description	Search for packets
//	@Tags			packets
//	@Accept			json
//	@Produce		json
//	@Param			query	query		string	true	"Search query"
//	@Param			limit	query		int		true	"Limit"
//	@Param			offset	query		int		true	"Offset"
//	@Success		200		{object}	schemas.SearchPacketResponse
//	@Failure		500		{object}	schemas.ErrorResponse
//	@Router			/api/v1/packets/search [get]
func makeSearchPacketEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.SearchPacketRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		packets, err := s.SearchPacket(ctx, req.Query, req.Limit, req.Offset)
		return schemas.SearchPacketResponse{Packets: packets}, err
	}
}

// makeAddPacketEndpoint constructs a AddPacket endpoint wrapping the service.
//
//	@Summary		Add packet
//	@Description	Add a new packet of products to the database
//	@Tags			packets
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.AddPacketResponse
//	@Failure		400	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/packets [post]
func makeAddPacketEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.AddPacketRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		id, err := s.AddPacket(ctx, &req.Packet, req.PacketContent)
		return schemas.AddPacketResponse{PacketID: id}, err
	}
}

// makeAddProductEndpoint constructs a AddProduct endpoint wrapping the service.
//
//	@Summary		Add product
//	@Description	Add a new product to the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.AddProductResponse
//	@Failure		400	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/products [post]
func makeAddProductEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.AddProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

		changeID, err := s.AddProduct(ctx, &req.ProductData)
		return schemas.AddProductResponse{ChangeID: changeID}, err
	}
}

// makeUpdateProductEndpoint constructs a UpdateProduct endpoint wrapping the service.
//
//	@Summary		Update product
//	@Description	Update product details in the database
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	schemas.UpdateProductResponse
//	@Failure		400	{object}	schemas.ErrorResponse
//	@Failure		500	{object}	schemas.ErrorResponse
//	@Router			/api/v1/products [put]
func makeUpdateProductEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.UpdateProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

		changeID, err := s.UpdateProduct(ctx, &req.ProductData)
		return schemas.UpdateProductResponse{ChangeID: changeID}, err
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
func makeDeleteProductEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.DeleteProductRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}

		changeID, err := s.DeleteProduct(ctx, req.ProductID)
		return schemas.DeleteProductResponse{ChangeID: changeID}, err
	}
}
