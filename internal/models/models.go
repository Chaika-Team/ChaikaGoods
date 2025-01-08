package models

// Product описывает товар.
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageurl"`
	SKU         string  `json:"sku"`
}

// Template описывает пакет товаров.
type Template struct {
	ID           int64             `json:"id"`
	TemplateName string            `json:"template_name"`
	Description  string            `json:"description"`
	Content      []TemplateContent `json:"content"`
}

// TemplateContent описывает одно содержимое пакета.
type TemplateContent struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
