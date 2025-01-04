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

type PackageSchema struct {
	ID          int64                  `json:"id"`
	PackageName string                 `json:"package_name"`
	Description string                 `json:"description"`
	Content     []PackageContentSchema `json:"content"`
}

type PackageContentSchema struct {
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
	Err      string          `json:"err,omitempty"`
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
	Err     string        `json:"err,omitempty"`
}

// SearchPacketRequest представляет собой запрос на поиск пакетов
// @Description Запрос на поиск пакетов
type SearchPacketRequest struct {
	Query  string `json:"query,omitempty"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

// SearchPacketResponse представляет собой ответ на запрос на поиск пакетов
// @Description Ответ на запрос на поиск пакетов
type SearchPacketResponse struct {
	Packets []PackageSchema `json:"packets"`
	Err     string          `json:"err,omitempty"`
}

// AddPacketRequest представляет собой запрос на добавление пакета
// @Description Запрос на добавление пакета
type AddPacketRequest struct {
	Packet PackageSchema `json:"packet"` // Сведения о новом пакете
}

// AddPacketResponse представляет собой ответ на запрос на добавление пакета
// @Description Ответ на запрос на добавление пакета
type AddPacketResponse struct {
	PacketID int64  `json:"id"` // ID созданного пакета
	Err      string `json:"err,omitempty"`
}

// GetPacketByIDRequest представляет собой запрос на получение пакета по его ID
// @Description Запрос на получение пакета по его ID
type GetPacketByIDRequest struct {
	PacketID int64 `json:"id"`
}

// GetPacketByIDResponse представляет собой ответ на запрос на получение пакета по его ID
// @Description Ответ на запрос на получение пакета по его ID
type GetPacketByIDResponse struct {
	Packet PackageSchema `json:"packet"`
	Err    string        `json:"err,omitempty"`
}

// CreateProductRequest представляет собой запрос на добавление продукта
// @Description Запрос на добавление продукта
type CreateProductRequest struct {
	Product ProductSchema `json:"product"`
}

// CreateProductResponse представляет собой ответ на запрос на добавление продукта
// @Description Ответ на запрос на добавление продукта
type CreateProductResponse struct {
	ProductID int64  `json:"id"`
	Err       string `json:"err,omitempty"`
}

// UpdateProductRequest представляет собой запрос на обновление продукта
// @Description Запрос на обновление продукта
type UpdateProductRequest struct {
	Product ProductSchema `json:"product"`
}

// UpdateProductResponse представляет собой ответ на запрос на обновление продукта
// @Description Ответ на запрос на обновление продукта
type UpdateProductResponse struct {
	Err string `json:"err,omitempty"`
}

// DeleteProductRequest представляет собой запрос на удаление продукта
// @Description Запрос на удаление продукта
type DeleteProductRequest struct {
	ProductID int64 `json:"id"`
}

// DeleteProductResponse представляет собой ответ на запрос на удаление продукта
// @Description Ответ на запрос на удаление продукта
type DeleteProductResponse struct {
	Err string `json:"err,omitempty"`
}
