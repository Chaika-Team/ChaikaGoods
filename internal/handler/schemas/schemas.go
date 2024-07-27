package schemas

type ProductSchema struct {
	ID          int64   `json:"id,omitempty"`
	Name        string  `json:"name,omitempty"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price,omitempty"`
	ImageURL    string  `json:"image_url,omitempty"`
	SKU         string  `json:"sku,omitempty"`
}

type ChangeSchema struct {
	ProductID int64       `json:"product_id"`
	Operation string      `json:"operation"` // "insert", "update", "delete"
	Details   interface{} `json:"details"`
}

type VersionSchema struct {
	VersionID int  `json:"version_id"`
	IsDev     bool `json:"is_dev"`
}

type PacketSchema struct {
	PacketID    int64           `json:"packet_id"`
	PacketName  string          `json:"packet_name"`
	Description string          `json:"description"`
	Products    []ProductSchema `json:"products"`
}

type GetAllProductsRequest struct {
	// Фильтры, пагинация и другие параметры могут быть добавлены здесь
}

type GetAllProductsResponse struct {
	Products []ProductSchema `json:"products"`
	Err      string          `json:"err,omitempty"` // пустое значение означает отсутствие ошибки
}

type GetCurrentVersionRequest struct{}

type GetCurrentVersionResponse struct {
	Version VersionSchema `json:"version"`
	Err     string        `json:"err,omitempty"`
}

type GetDeltaRequest struct {
	VersionID int `json:"version_id"`
}

type GetDeltaResponse struct {
	Changes []ChangeSchema `json:"changes"` // Change представляет собой структуру изменения
	Err     string         `json:"err,omitempty"`
}

type SearchPacketRequest struct {
	Query string `json:"query"` // Поисковый запрос для фильтрации пакетов
}

type SearchPacketResponse struct {
	Packets []PacketSchema `json:"packets"`
	Err     string         `json:"err,omitempty"`
}

type AddPacketRequest struct {
	Packet PacketSchema `json:"packet"` // Сведения о новом пакете
}

type AddPacketResponse struct {
	PacketID int64  `json:"packet_id"` // ID созданного пакета
	Err      string `json:"err,omitempty"`
}

type AddProductRequest struct {
	Product ProductSchema `json:"product"`
}

type AddProductResponse struct {
	ProductID int64  `json:"product_id"`
	Err       string `json:"err,omitempty"`
}

type UpdateProductRequest struct {
	ProductID int64         `json:"product_id"`
	Product   ProductSchema `json:"product"`
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
