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
func makeGetAllProductsEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(schemas.GetAllProductsRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		products, err := s.GetAllProducts(ctx)
		return schemas.GetAllProductsResponse{Products: products}, err

	}
}

// makeGetProductByIDEndpoint constructs a GetProductByID endpoint wrapping the service.
func makeGetProductByIDEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.GetProductByIDRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		product, err := s.GetProductByID(ctx, req.ProductID)
		return schemas.GetProductByIDResponse{Product: product}, err
	}

}

// makeGetCurrentVersionEndpoint constructs a GetCurrentVersion endpoint wrapping the service.
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
func makeGetDeltaEndpoint(logger log.Logger, s service.GoodsService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(schemas.GetDeltaRequest)
		if !ok {
			_ = level.Error(logger).Log("msg", "invalid request type")
			return nil, errors.New("invalid request type")
		}
		changes, err := s.GetDelta(ctx, req.VersionID)
		return schemas.GetDeltaResponse{Changes: changes}, err
	}
}

// makeSearchPacketEndpoint constructs a SearchPacket endpoint wrapping the service.
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
