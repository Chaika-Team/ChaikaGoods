package schemas

import "ChaikaGoods/internal/models"

type GetAllProductsRequest struct {
	// Фильтры, пагинация и другие параметры могут быть добавлены здесь
}

type GetAllProductsResponse struct {
	Products []models.Product `json:"products"`
	Err      string           `json:"err,omitempty"` // пустое значение означает отсутствие ошибки
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
	Product models.Product `json:"product"`
}

type AddProductResponse struct {
	ProductID int64  `json:"product_id"`
	Err       string `json:"err,omitempty"`
}

type UpdateProductRequest struct {
	ProductID int64          `json:"product_id"`
	Product   models.Product `json:"product"`
}

type UpdateProductResponse struct {
	Err string `json:"err,omitempty"`
}

type DeleteProductRequest struct {
	ProductID int64 `json:"product_id"`
}

type DeleteProductResponse struct {
	Err string `json:"err,omitempty"`
}
