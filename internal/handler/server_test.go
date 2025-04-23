package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/handler/schemas"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const (
	dummyEndpoint       = "/dummy"
	internalServerError = "internal server error"
)

// Техника тест-дизайна: Таблица принятия решений и попарное тестирование
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функций NewHTTPServer и RegisterRoutes.
//   - Таблица решений: для каждого маршрута (метод, URL, тело запроса) проверяется ожидаемый обработчик и статус.
//   - Попарное тестирование: комбинации HTTP-методов, URL и тела запроса покрывают различные варианты вызова маршрутов.
func TestNewHTTPServerRoutes(t *testing.T) {
	dummyEndpoints := Endpoints{
		GetAllProducts: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "GetAllProducts"}, nil
		},
		GetProductByID: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "GetProductByID"}, nil
		},
		SearchTemplates: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "SearchTemplates"}, nil
		},
		AddTemplate: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "AddTemplate"}, nil
		},
		GetTemplateByID: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "GetTemplateByID"}, nil
		},
		CreateProduct: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "CreateProduct"}, nil
		},
		UpdateProduct: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "UpdateProduct"}, nil
		},
		DeleteProduct: func(ctx context.Context, request interface{}) (interface{}, error) {
			return map[string]string{"handler": "DeleteProduct"}, nil
		},
	}
	logger := log.NewNopLogger()
	server := NewHTTPServer(logger, dummyEndpoints)

	tests := []struct {
		name       string
		method     string
		url        string
		body       string
		expHandler string
		expStatus  int
	}{
		{
			name:       "Get Product By ID",
			method:     "GET",
			url:        "/api/v1/product/123",
			body:       "",
			expHandler: "GetProductByID",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Get All Products",
			method:     "GET",
			url:        "/api/v1/product",
			body:       "",
			expHandler: "GetAllProducts",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Search Templates",
			method:     "GET",
			url:        "/api/v1/product/template/search?query=test&limit=1&offset=0",
			body:       "",
			expHandler: "SearchTemplates",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Add Template",
			method:     "POST",
			url:        "/api/v1/product/template",
			body:       `{"dummy":"data"}`,
			expHandler: "AddTemplate",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Get Template By ID",
			method:     "GET",
			url:        "/api/v1/product/template/456",
			body:       "",
			expHandler: "GetTemplateByID",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Create Product",
			method:     "POST",
			url:        "/api/v1/product",
			body:       `{"dummy":"data"}`,
			expHandler: "CreateProduct",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Update Product",
			method:     "PUT",
			url:        "/api/v1/product/789",
			body:       `{"dummy":"data"}`,
			expHandler: "UpdateProduct",
			expStatus:  http.StatusOK,
		},
		{
			name:       "Delete Product",
			method:     "DELETE",
			url:        "/api/v1/product/101",
			body:       "",
			expHandler: "DeleteProduct",
			expStatus:  http.StatusOK,
		},
		{
			name:      "Swagger Docs",
			method:    "GET",
			url:       "/api/v1/product/docs/index.html",
			body:      "",
			expStatus: http.StatusOK,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest(tc.method, tc.url, strings.NewReader(tc.body))
				req.Header.Set("Content-Type", "application/json")
			} else {
				req = httptest.NewRequest(tc.method, tc.url, nil)
			}

			rec := httptest.NewRecorder()
			server.ServeHTTP(rec, req)
			assert.Equal(t, tc.expStatus, rec.Code)

			if tc.expHandler != "" {
				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tc.expHandler, resp["handler"])
			}
		})
	}
}

// -----------------------
// Тесты для encodeResponse
// -----------------------

// Техника тест-дизайна: Причинно-следственный анализ
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции encodeResponse.
//   - Причинно-следственный анализ: корректный вход приводит к ожидаемому JSON-выходу.
func TestEncodeResponseCauseEffect(t *testing.T) {
	rec := httptest.NewRecorder()
	ctx := context.Background()
	responseObj := map[string]interface{}{"result": "success"}
	encoder := encodeResponse(log.NewNopLogger())

	err := encoder(ctx, rec, responseObj)
	assert.NoError(t, err)

	expected, err := json.Marshal(responseObj)
	assert.NoError(t, err)
	assert.JSONEq(t, string(expected), rec.Body.String())
}

// -----------------------------------
// Тест для EncodeErrorResponse
// -----------------------------------

// Техника тест-дизайна: Таблица принятия решений
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции EncodeErrorResponse с различными типами ошибок.
//   - Таблица принятия решений: для каждого типа ошибки определяются ожидаемые статус и сообщение.
func TestEncodeErrorResponseDecisionTable(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "NotFound error",
			err:            &myerr.AppError{Type: myerr.ErrorTypeNotFound, Message: "resource not found"},
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "resource not found",
		},
		{
			name:           "Validation error",
			err:            &myerr.AppError{Type: myerr.ErrorTypeValidation, Message: "validation failed"},
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "validation failed",
		},
		{
			name:           "Duplicate error",
			err:            &myerr.AppError{Type: myerr.ErrorTypeDuplicate, Message: "duplicate resource"},
			expectedStatus: http.StatusConflict,
			expectedMsg:    "duplicate resource",
		},
		{
			name:           "Unknown error",
			err:            &myerr.AppError{Type: myerr.ErrorTypeUnknown, Message: internalServerError},
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    internalServerError,
		},
		{
			name:           "Generic error",
			err:            errors.New("generic failure"),
			expectedStatus: http.StatusInternalServerError,
			expectedMsg:    internalServerError,
		},
	}

	logger := log.NewNopLogger()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			encoder := encodeErrorResponse(logger)
			encoder(context.Background(), tc.err, rec)
			assert.Equal(t, tc.expectedStatus, rec.Code)

			var resp map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedMsg, resp["error"])
		})
	}
}

// -----------------------------------
// Тесты для DecodeJSONRequest
// -----------------------------------

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции DecodeJSONRequest при отсутствии тела запроса.
//   - Прогнозирование ошибок: пустой запрос приводит к ошибке.
func TestDecodeJSONRequestErrorEmptyBody(t *testing.T) {
	req := httptest.NewRequest("POST", dummyEndpoint, nil)
	decoder := decodeJSONRequest(&struct{ Field string }{})
	_, err := decoder(context.Background(), req)
	assert.Error(t, err)
	assert.Equal(t, "empty request body", err.Error())
}

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции DecodeJSONRequest с некорректным JSON.
//   - Прогнозирование ошибок: синтаксическая ошибка в JSON должна вызвать ошибку.
func TestDecodeJSONRequestErrorInvalidJSON(t *testing.T) {
	invalidJSON := `{"Field": "value",}` // Некорректный JSON: лишняя запятая
	req := httptest.NewRequest("POST", dummyEndpoint, strings.NewReader(invalidJSON))
	decoder := decodeJSONRequest(&struct{ Field string }{})
	_, err := decoder(context.Background(), req)
	assert.Error(t, err)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции DecodeJSONRequest с корректными данными JSON.
//   - Классы эквивалентности: валидный запрос не порождает ошибки
func TestDecodeJSONRequestSuccess(t *testing.T) {
	validJSON := `{"Field": "value"}`
	req := httptest.NewRequest("POST", dummyEndpoint, strings.NewReader(validJSON))
	decoder := decodeJSONRequest(&struct{ Field string }{})
	_, err := decoder(context.Background(), req)
	assert.NoError(t, err)
}

// -----------------------------------
// Тесты для decodeRequestWithID
// -----------------------------------

// Техника тест-дизайна: Прогнозирование ошибок
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции decodeRequestWithID с нечисловым идентификатором.
//   - Прогнозирование ошибок: передача строки вместо числа приводит к ошибке.
func TestDecodeRequestWithIDBoundaryNonNumericID(t *testing.T) {
	logger := log.NewNopLogger()
	req := httptest.NewRequest("GET", "/product/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	schema := &schemas.GetProductByIDRequest{}
	decoder := decodeRequestWithID(logger, "id", schema)
	_, err := decoder(context.Background(), req)
	assert.Error(t, err)
}

// Техника тест-дизайна: Таблица принятия решений
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции decodeRequestWithID для поддерживаемых схем.
//   - Таблица принятия решений: проверка соответствия запросов каждой схеме
func TestDecodeRequestWithIDDecisionTableValidSchemas(t *testing.T) {
	logger := log.NewNopLogger()
	tests := []struct {
		name       string
		url        string
		schema     interface{}
		expectedID int64
	}{
		{
			name:       "GetProductByIDRequest",
			url:        "/product/101",
			schema:     &schemas.GetProductByIDRequest{},
			expectedID: 101,
		},
		{
			name:       "GetTemplateByIDRequest",
			url:        "/product/template/202",
			schema:     &schemas.GetTemplateByIDRequest{},
			expectedID: 202,
		},
		{
			name:       "DeleteProductRequest",
			url:        "/product/303",
			schema:     &schemas.DeleteProductRequest{},
			expectedID: 303,
		},
		{
			name:       "AddTemplateRequest",
			url:        "/template/404",
			schema:     &schemas.AddTemplateRequest{},
			expectedID: 404,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tc.url, nil)
			parts := strings.Split(tc.url, "/")
			idStr := parts[len(parts)-1]
			req = mux.SetURLVars(req, map[string]string{"id": idStr})
			decoder := decodeRequestWithID(logger, "id", tc.schema)
			result, err := decoder(context.Background(), req)

			switch s := result.(type) {
			case *schemas.GetProductByIDRequest:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, s.ProductID)
			case *schemas.GetTemplateByIDRequest:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, s.TemplateID)
			case *schemas.DeleteProductRequest:
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedID, s.ProductID)
			default:
				assert.Error(t, err)
			}
		})
	}
}

// -----------------------------------
// Тесты для decodeSearchTemplatesRequest
// -----------------------------------

// Техника тест-дизайна: Таблица принятия решений
// Автор: Дмитрий Фоломкин
// Описание:
//   - Тест для функции decodeSearchTemplatesRequest с различными комбинациями параметров запроса.
//   - Таблица принятия решений: проверка граничных и ошибочных значений для параметров query, limit и offset.
func TestDecodeSearchTemplatesRequestDecisionTable(t *testing.T) {
	tests := []struct {
		name            string
		queryParams     string
		expectedError   string
		expectedRequest *schemas.SearchTemplatesRequest
	}{
		{
			name:          "Valid parameters",
			queryParams:   "query=example&limit=5&offset=0",
			expectedError: "",
			expectedRequest: &schemas.SearchTemplatesRequest{
				Query:  "example",
				Limit:  5,
				Offset: 0,
			},
		},
		{
			name:            "Invalid limit (zero)",
			queryParams:     "query=example&limit=0&offset=10",
			expectedError:   "invalid or missing limit parameter",
			expectedRequest: nil,
		},
		{
			name:            "Invalid offset (negative)",
			queryParams:     "query=example&limit=5&offset=-1",
			expectedError:   "invalid or missing offset parameter",
			expectedRequest: nil,
		},
		{
			name:            "Non-numeric limit",
			queryParams:     "query=example&limit=abc&offset=10",
			expectedError:   "invalid or missing limit parameter",
			expectedRequest: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/template/search?"+tc.queryParams, nil)
			result, err := decodeSearchTemplatesRequest(context.Background(), req)
			if tc.expectedError != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.expectedError, err.Error())
			} else {
				assert.NoError(t, err)
				stReq, ok := result.(*schemas.SearchTemplatesRequest)
				assert.True(t, ok)
				assert.Equal(t, tc.expectedRequest.Query, stReq.Query)
				assert.Equal(t, tc.expectedRequest.Limit, stReq.Limit)
				assert.Equal(t, tc.expectedRequest.Offset, stReq.Offset)
			}
		})
	}
}
