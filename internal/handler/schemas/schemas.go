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

// SearchTemplatesRequest представляет собой запрос на поиск шаблонов
// @Description Запрос на поиск шаблонов
type SearchTemplatesRequest struct {
	Query  string `json:"query,omitempty"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

// SearchTemplatesResponse представляет собой ответ на запрос на поиск шаблонов
// @Description Ответ на запрос на поиск шаблонов
type SearchTemplatesResponse struct {
	Templates []TemplateSchema `json:"templates"`
}

// AddTemplateRequest представляет собой запрос на добавление шаблона
// @Description Запрос на добавление шаблона
type AddTemplateRequest struct {
	Template TemplateSchema `json:"template"` // Сведения о новом шаблоне
}

// AddTemplateResponse представляет собой ответ на запрос на добавление шаблона
// @Description Ответ на запрос на добавление шаблона
type AddTemplateResponse struct {
	TemplateID int64 `json:"id"` // ID созданного шаблона
}

// GetTemplateByIDRequest представляет собой запрос на получение шаблона по его ID
// @Description Запрос на получение шаблона по его ID
type GetTemplateByIDRequest struct {
	TemplateID int64 `json:"id"`
}

// GetTemplateByIDResponse представляет собой ответ на запрос на получение шаблона по его ID
// @Description Ответ на запрос на получение шаблона по его ID
type GetTemplateByIDResponse struct {
	Template TemplateSchema `json:"template"`
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
