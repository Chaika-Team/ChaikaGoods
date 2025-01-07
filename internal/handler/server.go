package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	_ "github.com/Chaika-Team/ChaikaGoods/docs"
	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	httpGoKit "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Constants for logging messages
const (
	decoderReturningMsg = "decoder returning"
)

// NewHTTPServer initializes and returns a new HTTP server with all routes defined.
func NewHTTPServer(logger log.Logger, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(HTTPLoggingMiddleware(logger), HeaderMiddleware)
	// Swagger UI
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Register file and folder routes
	registerRoutes(logger, r.PathPrefix("/api/v1").Subrouter(), endpoints)

	return r
}

// registerRoutes registers all routes for the service using go-kit transport.
func registerRoutes(logger log.Logger, api *mux.Router, endpoints Endpoints) {
	// GetAllProducts    endpoint.Endpoint
	// GetProductByID    endpoint.Endpoint
	// For Packages
	// SearchPackages endpoint.Endpoint
	// AddPackage    endpoint.Endpoint
	// For products (admin)
	// CreateProduct    endpoint.Endpoint
	// UpdateProduct endpoint.Endpoint
	// DeleteProduct endpoint.Endpoint

	// Get all products
	api.Methods("GET").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.GetAllProducts,
		decodeJSONRequest(&schemas.GetAllProductsRequest{}),
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

	// Search Package
	api.Methods("GET").Path("/Packages/search").Handler(httpGoKit.NewServer(
		endpoints.SearchPackages,
		decodeSearchPackagesRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add Package
	api.Methods("POST").Path("/Packages").Handler(httpGoKit.NewServer(
		endpoints.AddPackage,
		decodeJSONRequest(&schemas.AddPackageRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Get Package by ID
	api.Methods("GET").Path("/Packages/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetPackageByID,
		decodeRequestWithID(logger, "id", &schemas.GetPackageByIDRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add product
	api.Methods("POST").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.CreateProduct,
		decodeJSONRequest(&schemas.CreateProductRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Update product
	api.Methods("PUT").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.UpdateProduct,
		decodeJSONRequest(&schemas.UpdateProductRequest{}),
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

// encodeResponse encodes the response as JSON.
func encodeResponse(_ log.Logger) httpGoKit.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		// Common response encoding logic here
		return json.NewEncoder(w).Encode(response)
	}
}

// encodeErrorResponse encodes the error response as JSON with appropriate status code.
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

func decodeJSONRequest(schema interface{}) httpGoKit.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		if req.Body == nil {
			return nil, errors.New("empty request body")
		}
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				_ = level.Error(log.NewNopLogger()).Log("msg", "error closing request body", "err", err)
			}
		}(req.Body)

		if err := json.NewDecoder(req.Body).Decode(schema); err != nil {
			return nil, err
		}
		return schema, nil
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
			_ = level.Debug(logger).Log("msg", decoderReturningMsg, "type", fmt.Sprintf("%T", s))
			return s, nil
		case *schemas.GetPackageByIDRequest:
			s.PackageID = id
			_ = level.Debug(logger).Log("msg", decoderReturningMsg, "type", fmt.Sprintf("%T", s))
			return s, nil
		case *schemas.DeleteProductRequest:
			s.ProductID = id
			_ = level.Debug(logger).Log("msg", decoderReturningMsg, "type", fmt.Sprintf("%T", s))
			return s, nil
		default:
			return nil, errors.New("unsupported schema type")
		}
	}
}

// decodeSearchPackagesRequest декодирует GET запрос с параметрами query, limit и offset.
func decodeSearchPackagesRequest(_ context.Context, req *http.Request) (interface{}, error) {
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

	return schemas.SearchPackagesRequest{
		Query:  searchString,
		Limit:  limit,
		Offset: offset,
	}, nil
}
