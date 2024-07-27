package handler

import (
	"ChaikaGoods/internal/service"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
)

// Endpoints holds all Go kit endpoints for all operations
type Endpoints struct {
	// For products
	GetAllProducts    endpoint.Endpoint
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
func MakeEndpoints(logger log.Logger, service service.Service) Endpoints {
	return Endpoints{
		// Products
		GetAllProducts:    makeGetAllProductsEndpoint(logger, service),
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
