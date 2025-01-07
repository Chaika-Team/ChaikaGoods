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

// SearchPackagesRequest представляет собой запрос на поиск пакетов
// @Description Запрос на поиск пакетов
type SearchPackagesRequest struct {
	Query  string `json:"query,omitempty"`
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

// SearchPackagesResponse представляет собой ответ на запрос на поиск пакетов
// @Description Ответ на запрос на поиск пакетов
type SearchPackagesResponse struct {
	Packages []PackageSchema `json:"Packages"`
}

// AddPackageRequest представляет собой запрос на добавление пакета
// @Description Запрос на добавление пакета
type AddPackageRequest struct {
	Package PackageSchema `json:"Package"` // Сведения о новом пакете
}

// AddPackageResponse представляет собой ответ на запрос на добавление пакета
// @Description Ответ на запрос на добавление пакета
type AddPackageResponse struct {
	PackageID int64 `json:"id"` // ID созданного пакета
}

// GetPackageByIDRequest представляет собой запрос на получение пакета по его ID
// @Description Запрос на получение пакета по его ID
type GetPackageByIDRequest struct {
	PackageID int64 `json:"id"`
}

// GetPackageByIDResponse представляет собой ответ на запрос на получение пакета по его ID
// @Description Ответ на запрос на получение пакета по его ID
type GetPackageByIDResponse struct {
	Package PackageSchema `json:"Package"`
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
