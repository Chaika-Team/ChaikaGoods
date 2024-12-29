package handler

import (
	"ChaikaGoods/internal/handler/schemas"
	"context"
	"encoding/json"
	"github.com/go-kit/kit/endpoint"
	httpGoKit "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"net/http"
	"reflect"
)

// NewHTTPServer initializes and returns a new HTTP server with all routes defined.
func NewHTTPServer(logger log.Logger, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(commonMiddleware(logger))

	// Swagger UI
	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	// Apply logging middleware to all endpoints
	wrapEndpointsWithLogging(logger, &endpoints)

	// Register file and folder routes
	registerRoutes(logger, r, endpoints)

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

// wrapEndpointsWithLogging applies logging middleware to all fields in the Endpoints struct using reflection.
func wrapEndpointsWithLogging(logger log.Logger, endpoints interface{}) {
	loggingMiddleware := makeLoggingMiddleware(logger)
	endpointsVal := reflect.ValueOf(endpoints).Elem()

	for i := 0; i < endpointsVal.NumField(); i++ {
		field := endpointsVal.Field(i)
		if field.CanInterface() {
			ep, ok := field.Interface().(endpoint.Endpoint)
			if ok {
				wrapped := loggingMiddleware(ep)
				if field.CanSet() {
					field.Set(reflect.ValueOf(wrapped))
				}
			}
		}
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
func registerRoutes(logger log.Logger, r *mux.Router, endpoints Endpoints) {
	api := r.PathPrefix("/api/v1").Subrouter()

	// Register Product Endpoints
	api.Methods("GET").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.GetAllProducts,
		decodeGetAllProductsRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	api.Methods("GET").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetProductByID,
		decodeGetProductByIDRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	api.Methods("POST").Path("/products").Handler(httpGoKit.NewServer(
		endpoints.AddProduct,
		decodeAddProductRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	api.Methods("PUT").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.UpdateProduct,
		decodeUpdateProductRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	api.Methods("DELETE").Path("/products/{id}").Handler(httpGoKit.NewServer(
		endpoints.DeleteProduct,
		decodeDeleteProductRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Register additional routes similarly
}

// Example encoder function for responses
func encodeResponse(logger log.Logger) httpGoKit.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		// Common response encoding logic here
		return json.NewEncoder(w).Encode(response)
	}
}

// Example error encoder function
func encodeErrorResponse(logger log.Logger) httpGoKit.ErrorEncoder {
	return func(ctx context.Context, err error, w http.ResponseWriter) {
		// Common error encoding logic here
		w.WriteHeader(determineHTTPError(err))
		err = json.NewEncoder(w).Encode(map[string]interface{}{
			"error": err.Error(),
		})
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to encode error", "err", err)
		}
	}
}

// Helper function to determine HTTP status code from error
func determineHTTPError(err error) int {
	// Implement error to HTTP status code mapping
	return http.StatusInternalServerError
}

// decodeGetAllProductsRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeGetAllProductsRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.GetAllProductsRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

// decodeGetProductByIDRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeGetProductByIDRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.GetProductByIDRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

// decodeAddProductRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeAddProductRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.AddProductRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

// decodeUpdateProductRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeUpdateProductRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.UpdateProductRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}

// decodeDeleteProductRequest is a transport/http.DecodeRequestFunc that decodes a JSON-encoded request from the HTTP request body.
func decodeDeleteProductRequest(ctx context.Context, req *http.Request) (request interface{}, err error) {
	var r schemas.DeleteProductRequest
	err = json.NewDecoder(req.Body).Decode(&r)
	return r, err
}
