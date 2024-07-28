package schemas

import "ChaikaGoods/internal/models"

type GetAllProductsRequest struct {
	// Фильтры, пагинация и другие параметры могут быть добавлены здесь
}

type GetAllProductsResponse struct {
	Products []models.Product `json:"products"`
	Err      string           `json:"err,omitempty"` // пустое значение означает отсутствие ошибки
}

type GetProductByIDRequest struct {
	ProductID int64 `json:"product_id"`
}

type GetProductByIDResponse struct {
	Product models.Product `json:"product"`
	Err     string         `json:"err,omitempty"`
}

type GetCurrentVersionRequest struct{}

type GetCurrentVersionResponse struct {
	Version models.Version `json:"version"`
	Err     string         `json:"err,omitempty"`
}

type GetDeltaRequest struct {
	VersionID int `json:"version_id"`
}

type GetDeltaResponse struct {
	Changes []models.Change `json:"changes"` // Change представляет собой структуру изменения
	Err     string          `json:"err,omitempty"`
}

type SearchPacketRequest struct {
	Query  string `json:"query"` // Поисковый запрос для фильтрации пакетов
	Limit  int64  `json:"limit"`
	Offset int64  `json:"offset"`
}

type SearchPacketResponse struct {
	Packets []models.Package `json:"packets"`
	Err     string           `json:"err,omitempty"`
}

type AddPacketRequest struct {
	Packet        models.Package          `json:"packet"`         // Сведения о новом пакете
	PacketContent []models.PackageContent `json:"packet_content"` // Содержимое пакета
}

type AddPacketResponse struct {
	PacketID int64  `json:"packet_id"` // ID созданного пакета
	Err      string `json:"err,omitempty"`
}

type AddProductRequest struct {
	ProductData map[string]interface{} `json:"product_data"`
}

type AddProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}

type UpdateProductRequest struct {
	ProductData map[string]interface{} `json:"product_data"`
}

type UpdateProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}

type DeleteProductRequest struct {
	ProductID int64 `json:"product_id"`
}

type DeleteProductResponse struct {
	ChangeID int64  `json:"change_id"`
	Err      string `json:"err,omitempty"`
}
