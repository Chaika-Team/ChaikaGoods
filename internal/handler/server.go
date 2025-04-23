package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/Chaika-Team/ChaikaGoods/docs"
	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	httpGoKit "github.com/go-kit/kit/transport/http"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/gorilla/mux"
)

// Constants for logging messages and endpoint
const (
	// apiPrefix is the common prefix for all API routes
	apiPrefix = "/api"
	// v1Prefix is the prefix for API v1 routes
	v1Prefix            = apiPrefix + "/v1/product"
	decoderReturningMsg = "decoder returning"
)

// NewHTTPServer initializes and returns a new HTTP server with all the necessary routes and middleware.
// It sets up the router with middleware and registers all API endpoints.
//
// Parameters:
//   - logger: Logger instance for logging HTTP requests and errors
//   - endpoints: Collection of service endpoints to be exposed via HTTP
//
// Returns:
//   - http.Handler: Configured HTTP handler with all routes and middleware
func NewHTTPServer(logger log.Logger, endpoints Endpoints) http.Handler {
	r := mux.NewRouter()
	r.Use(HTTPLoggingMiddleware(logger), HeaderMiddleware)

	// API information endpoint
	r.HandleFunc(apiPrefix, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(map[string]interface{}{
			"versions":      []string{"v1"},
			"documentation": v1Prefix + "/docs",
		})
		if err != nil {
			_ = level.Error(logger).Log("msg", "failed to encode API info", "err", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}).Methods("GET")

	// Register API v1 routes
	registerV1Routes(logger, r, endpoints)

	// Future versions can be added here as needed
	// registerV2Routes(logger, r, endpoints)

	return r
}

// registerV1Routes registers all routes for the service using go-kit transport.
func registerV1Routes(logger log.Logger, router *mux.Router, endpoints Endpoints) {
	// GetAllProducts    endpoint.Endpoint
	// GetProductByID    endpoint.Endpoint
	// For templates
	// SearchTemplates endpoint.Endpoint
	// AddTemplate    endpoint.Endpoint
	// For products (admin)
	// CreateProduct    endpoint.Endpoint
	// UpdateProduct endpoint.Endpoint
	// DeleteProduct endpoint.Endpoint
	v1 := router.PathPrefix(v1Prefix).Subrouter()

	// Swagger UI at /api/v1/product/docs/
	v1.PathPrefix("/docs/").Handler(httpSwagger.Handler(
		httpSwagger.URL("./doc.json"), // The URL points to API definition (relative path)
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("list"),
		httpSwagger.DomID("swagger-ui"),
	))

	// Get all products
	v1.Methods("GET").Path("").Handler(httpGoKit.NewServer(
		endpoints.GetAllProducts,
		decodeEmptyRequest[schemas.GetAllProductsRequest](),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Get product by ID
	v1.Methods("GET").Path("/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetProductByID,
		decodeRequestWithID(logger, "id", &schemas.GetProductByIDRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Search Template
	v1.Methods("GET").Path("/template/search").Handler(httpGoKit.NewServer(
		endpoints.SearchTemplates,
		decodeSearchTemplatesRequest,
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add Template
	v1.Methods("POST").Path("/template").Handler(httpGoKit.NewServer(
		endpoints.AddTemplate,
		decodeJSONRequest(&schemas.AddTemplateRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Get Template by ID
	v1.Methods("GET").Path("/template/{id}").Handler(httpGoKit.NewServer(
		endpoints.GetTemplateByID,
		decodeRequestWithID(logger, "id", &schemas.GetTemplateByIDRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Add product
	v1.Methods("POST").Path("").Handler(httpGoKit.NewServer(
		endpoints.CreateProduct,
		decodeJSONRequest(&schemas.CreateProductRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Update product
	v1.Methods("PUT").Path("/{id}").Handler(httpGoKit.NewServer(
		endpoints.UpdateProduct,
		decodeJSONRequest(&schemas.UpdateProductRequest{}),
		encodeResponse(logger),
		httpGoKit.ServerErrorEncoder(encodeErrorResponse(logger)),
	))

	// Delete product
	v1.Methods("DELETE").Path("/{id}").Handler(httpGoKit.NewServer(
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
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				_ = level.Error(log.With(ctx.Value("logger").(log.Logger))).Log("msg", "failed to close request body", "err", err)
			}
		}(req.Body)
		err := json.NewDecoder(req.Body).Decode(schema)
		if err == io.EOF {
			return nil, errors.New("empty request body")
		} else if err != nil {
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
		case *schemas.GetTemplateByIDRequest:
			s.TemplateID = id
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

// decodeEmptyRequest безопасно возвращает пустую структуру запроса.
func decodeEmptyRequest[T any]() func(context.Context, *http.Request) (interface{}, error) {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		return new(T), nil
	}
}

// decodeSearchTemplatesRequest декодирует GET запрос с параметрами query, limit и offset.
func decodeSearchTemplatesRequest(_ context.Context, req *http.Request) (interface{}, error) {
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

	return &schemas.SearchTemplatesRequest{
		Query:  searchString,
		Limit:  limit,
		Offset: offset,
	}, nil
}
