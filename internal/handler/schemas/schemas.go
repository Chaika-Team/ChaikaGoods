package schemas

import (
	_ "ChaikaGoods/docs"
	"ChaikaGoods/internal/models"
	"time"
)

// GetAllProductsRequest представляет собой запрос на получение всех продуктов
// @Description Запрос на получение всех продуктов
type GetAllProductsRequest struct {
	// Фильтры, пагинация и другие параметры могут быть добавлены здесь
}

// GetAllProductsResponse представляет собой ответ на запрос на получение всех продуктов
// @Description Ответ на запрос на получение всех продуктов
type GetAllProductsResponse struct {
	Products []models.Product `json:"products"`
	Err      string           `json:"err,omitempty"`
}

// GetProductByIDRequest представляет собой запрос на получение продукта по его ID
// @Description Запрос на получение продукта по его ID
type GetProductByIDRequest struct {
	ProductID int64 `json:"product_id"`
}

// GetProductByIDResponse представляет собой ответ на запрос на получение продукта по его ID
// @Description Ответ на запрос на получение продукта по его ID
type GetProductByIDResponse struct {
	Product models.Product `json:"product"`
	Err     string         `json:"err,omitempty"`
}

// GetCurrentVersionRequest представляет собой запрос на получение текущей версии
// @Description Запрос на получение текущей версии
type GetCurrentVersionRequest struct{}

// GetCurrentVersionResponse представляет собой ответ на запрос на получение текущей версии
// @Description Ответ на запрос на получение текущей версии
type GetCurrentVersionResponse struct {
	Version models.Version `json:"version"`
	Err     string         `json:"err,omitempty"`
}

// GetDeltaRequest представляет собой запрос на получение изменений в базе данных
// @Description Запрос на получение изменений в базе данных
type GetDeltaRequest struct {
	VersionID int `json:"version_id"`
}

// GetDeltaResponse представляет собой ответ на запрос на получение изменений в базе данных
// @Description Ответ на запрос на получение изменений в базе данных
type GetDeltaResponse struct {
	Changes []ChangeSchema `json:"changes"` // Change представляет собой структуру изменения
	Err     string         `json:"err,omitempty"`
}

type ChangeSchema struct {
	ID              int64     `json:"id"`
	VersionID       int       `json:"version_id"`
	OperationType   int       `json:"operation_type"`
	NewValue        string    `json:"new_value"` // NewValue as a JSON string
	ChangeTimestamp time.Time `json:"change_timestamp"`
	Considered      bool      `json:"considered"`
}

// SearchPacketRequest представляет собой запрос на поиск пакетов
// @Description Запрос на поиск пакетов
type SearchPacketRequest struct {
	Query  string `json:"query"` // Поисковый запрос для фильтрации пакетов
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

// SearchPacketResponse представляет собой ответ на запрос на поиск пакетов
// @Description Ответ на запрос на поиск пакетов
type SearchPacketResponse struct {
	Packets []models.Package `json:"packets"`
	Err     string           `json:"err,omitempty"`
}

// AddPacketRequest представляет собой запрос на добавление пакета
// @Description Запрос на добавление пакета
type AddPacketRequest struct {
	Packet        models.Package          `json:"packet"`         // Сведения о новом пакете
	PacketContent []models.PackageContent `json:"packet_content"` // Содержимое пакета
}

// AddPacketResponse представляет собой ответ на запрос на добавление пакета
// @Description Ответ на запрос на добавление пакета
type AddPacketResponse struct {
	PacketID int64  `json:"packet_id"` // ID созданного пакета
	Err      string `json:"err,omitempty"`
}

// AddProductRequest представляет собой запрос на добавление продукта
// @Description Запрос на добавление продукта
type AddProductRequest struct {
	ProductData map[string]interface{} `json:"product_data"`
}

// AddProductResponse представляет собой ответ на запрос на добавление продукта
// @Description Ответ на запрос на добавление продукта
type AddProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}

// UpdateProductRequest представляет собой запрос на обновление продукта
// @Description Запрос на обновление продукта
type UpdateProductRequest struct {
	ProductData map[string]interface{} `json:"product_data"`
}

// UpdateProductResponse представляет собой ответ на запрос на обновление продукта
// @Description Ответ на запрос на обновление продукта
type UpdateProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}

// DeleteProductRequest представляет собой запрос на удаление продукта
// @Description Запрос на удаление продукта
type DeleteProductRequest struct {
	ProductID int64 `json:"product_id"`
}

// DeleteProductResponse представляет собой ответ на запрос на удаление продукта
// @Description Ответ на запрос на удаление продукта
type DeleteProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}
