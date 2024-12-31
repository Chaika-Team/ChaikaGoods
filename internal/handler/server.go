package handler

import (
	_ "ChaikaGoods/docs"
	"ChaikaGoods/internal/handler/schemas"
	"ChaikaGoods/internal/myerr"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-kit/kit/endpoint"
	httpGoKit "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewHTTPServer initializes and returns a new HTTP server with all routes defined.
func NewHTTPServer(logger log.Logger, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware(logger))
	// Swagger UI
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Register file and folder routes
	registerRoutes(logger, r.PathPrefix("/api/v1").Subrouter(), endpoints)

	return r
}

// commonMiddleware adds common HTTP headers to all responses.
func commonMiddleware(logger log.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			_ = level.Info(logger).Log(
				"msg", "received request",
				"method", r.Method,
				"url", r.URL.String(),
			)
			w.Header().Set("Content-Type", "application/json")
			next.ServeHTTP(w, r)
			_ = level.Info(logger).Log(
				"msg", "handled request",
				"method", r.Method,
				"url", r.URL.String(),
			)
		})
	}
}

// makeLoggingMiddleware creates a logging middleware for an endpoint.
func makeLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			_ = level.Info(logger).Log("msg", "calling endpoint", "request", request)
			response, err = next(ctx, request)
			_ = level.Info(logger).Log("msg", "called endpoint", "response", response, "err", err)
			return
		}
	}
}

// registerRoutes registers all routes for the service using go-kit transport.
func registerRoutes(logger log.Logger, api *mux.Router, endpoints Endpoints) {
	// GetAllProducts    endpoint.Endpoint
	// GetProductByID    endpoint.Endpoint
	// For packets
	// SearchPacket endpoint.Endpoint
	// AddPacket    endpoint.Endpoint
	// For products (admin)
	// AddProduct    endpoint.Endpoint
	// UpdateProduct endpoint.Endpoint
	// DeleteProduct endpoint.Endpoint

	// Get all products
	api.Methods("GET").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.GetAllProducts,
		decodeGetAllProductsRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Get product by ID
	api.Methods("GET").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetProductByID,
		decodeRequestWithID(logger, "id", &schemas.GetProductByIDRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Search packet
	api.Methods("GET").Path("/packets/search").Handler(httpGoKit.NewServer(
		endpoints.SearchPacket,
		decodeSearchPacketRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add packet
	api.Methods("POST").Path("/packets").Handler(httpGoKit.NewServer(
		endpoints.AddPacket,
		decodeAddPacketRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Get packet by ID
	api.Methods("GET").Path("/packets/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetPacketByID,
		decodeRequestWithID(logger, "id", &schemas.GetPacketByIDRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add product
	api.Methods("POST").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.AddProduct,
		decodeAddProductRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Update product
	api.Methods("PUT").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.UpdateProduct,
		decodeUpdateProductRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Delete product
	api.Methods("DELETE").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.DeleteProduct,
		decodeRequestWithID(logger, "id", &schemas.DeleteProductRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))
}

// Example encoder function for responses
func encodeResponse(_ log.Logger) httpGoKit.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		// Common response encoding logic here
		return json.NewEncoder(w).Encode(response)
	}
}

// Example error encoder function
func encodeErrorResponse(logger log.Logger) httpGoKit.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		var status int
		var response interface{}

		var e *myerr.AppError
		switch {
		case errors.As(err, &e):
			switch e.Type {
			case myerr.ErrorTypeNotFound:
				status = http.StatusNotFound
			case myerr.ErrorTypeValidation:
				status = http.StatusBadRequest
			case myerr.ErrorTypeDuplicate:
				status = http.StatusConflict
			default:
				status = http.StatusInternalServerError
			}
			response = map[string]string{"error": e.Message}
		default:
			status = http.StatusInternalServerError
			_ = level.Error(logger).Log("msg", "handling error", "err", err)
			response = map[string]string{"error": "internal server error"}
		}

		w.WriteHeader(status)
		_ = json.NewEncoder(w).Encode(response)
	}
}

// decodeRequestWithID генерирует DecodeRequestFunc для запросов с ID в пути.
func decodeRequestWithID(logger log.Logger, paramName string, schema interface{}) httpGoKit.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		id, err := extractID(req, paramName)
		if err != nil {
			return nil, err
		}

		switch s := schema.(type) {
		case *schemas.GetProductByIDRequest:
			s.ProductID = id
			_ = level.Debug(logger).Log("msg", "decoder returning", "type", fmt.Sprintf("%T", s))
			return s, nil
		case *schemas.GetPacketByIDRequest:
			s.PacketID = id
			_ = level.Debug(logger).Log("msg", "decoder returning", "type", fmt.Sprintf("%T", s))
			return s, nil
		case *schemas.DeleteProductRequest:
			s.ProductID = id
			_ = level.Debug(logger).Log("msg", "decoder returning", "type", fmt.Sprintf("%T", s))
			return s, nil
		default:
			return nil, errors.New("unsupported schema type")
		}
	}
}

// decodeGetAllProductsRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeGetAllProductsRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	// Since GET requests typically do not have a body, return an empty request
	return schemas.GetAllProductsRequest{}, nil
}

// decodeAddProductRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeAddProductRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.AddProductRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

// decodeUpdateProductRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeUpdateProductRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.UpdateProductRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

func decodeSearchPacketRequest(_ context.Context, req *http.Request) (interface{}, error) {
	query := req.URL.Query()

	searchString := query.Get("query")
	limit, err := strconv.ParseInt(query.Get("limit"), 10, 64)
	if err != nil || limit <= 0 {
		return nil, errors.New("invalid or missing limit parameter")
	}

	offset, err := strconv.ParseInt(query.Get("offset"), 10, 64)
	if err != nil || offset < 0 {
		return nil, errors.New("invalid or missing offset parameter")
	}

	return schemas.SearchPacketRequest{
		Query:  searchString,
		Limit:  limit,
		Offset: offset,
	}, nil
}

// decodeAddPacketRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeAddPacketRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.AddPacketRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}
