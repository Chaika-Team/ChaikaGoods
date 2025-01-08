package schemas

import (
	_ "github.com/Chaika-Team/ChaikaGoods/docs"
)

type ProductSchema struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageurl"`
}

type TemplateSchema struct {
	ID           int64                   `json:"id"`
	TemplateName string                  `json:"template_name"`
	Description  string                  `json:"description"`
	Content      []TemplateContentSchema `json:"content"`
}

type TemplateContentSchema struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// GetAllProductsRequest представляет собой запрос на получение всех продуктов
// @Description Запрос на получение всех продуктов
type GetAllProductsRequest struct {
	// TODO: Фильтры, пагинация и другие параметры могут быть тут
}

// GetAllProductsResponse представляет собой ответ на запрос на получение всех продуктов
// @Description Ответ на запрос на получение всех продуктов
type GetAllProductsResponse struct {
	Products []ProductSchema `json:"products"`
}

// GetProductByIDRequest представляет собой запрос на получение продукта по его ID
// @Description Запрос на получение продукта по его ID
type GetProductByIDRequest struct {
	ProductID int64 `json:"id"`
}

// GetProductByIDResponse представляет собой ответ на запрос на получение продукта по его ID
// @Description Ответ на запрос на получение продукта по его ID
type GetProductByIDResponse struct {
	Product ProductSchema `json:"product"`
}

// SearchTemplatesRequest представляет собой запрос на поиск пакетов
// @Description Запрос на поиск пакетов
type SearchTemplatesRequest struct {
	Query  string `json:"query,omitempty"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

// SearchTemplatesResponse представляет собой ответ на запрос на поиск пакетов
// @Description Ответ на запрос на поиск пакетов
type SearchTemplatesResponse struct {
	Templates []TemplateSchema `json:"Templates"`
}

// AddTemplateRequest представляет собой запрос на добавление пакета
// @Description Запрос на добавление пакета
type AddTemplateRequest struct {
	Template TemplateSchema `json:"Template"` // Сведения о новом пакете
}

// AddTemplateResponse представляет собой ответ на запрос на добавление пакета
// @Description Ответ на запрос на добавление пакета
type AddTemplateResponse struct {
	TemplateID int64 `json:"id"` // ID созданного пакета
}

// GetTemplateByIDRequest представляет собой запрос на получение пакета по его ID
// @Description Запрос на получение пакета по его ID
type GetTemplateByIDRequest struct {
	TemplateID int64 `json:"id"`
}

// GetTemplateByIDResponse представляет собой ответ на запрос на получение пакета по его ID
// @Description Ответ на запрос на получение пакета по его ID
type GetTemplateByIDResponse struct {
	Template TemplateSchema `json:"Template"`
}

// CreateProductRequest представляет собой запрос на добавление продукта
// @Description Запрос на добавление продукта
type CreateProductRequest struct {
	Product ProductSchema `json:"product"`
}

// CreateProductResponse представляет собой ответ на запрос на добавление продукта
// @Description Ответ на запрос на добавление продукта
type CreateProductResponse struct {
	ProductID int64 `json:"id"`
}

// UpdateProductRequest представляет собой запрос на обновление продукта
// @Description Запрос на обновление продукта
type UpdateProductRequest struct {
	Product ProductSchema `json:"product"`
}

// UpdateProductResponse представляет собой ответ на запрос на обновление продукта
// @Description Ответ на запрос на обновление продукта
type UpdateProductResponse struct {
}

// DeleteProductRequest представляет собой запрос на удаление продукта
// @Description Запрос на удаление продукта
type DeleteProductRequest struct {
	ProductID int64 `json:"id"`
}

// DeleteProductResponse представляет собой ответ на запрос на удаление продукта
// @Description Ответ на запрос на удаление продукта
type DeleteProductResponse struct {
}
