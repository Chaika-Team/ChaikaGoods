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

// Package описывает пакет товаров.
type Package struct {
	ID          int64            `json:"id"`
	PackageName string           `json:"package_name"`
	Description string           `json:"description"`
	Content     []PackageContent `json:"content"`
}

// PackageContent описывает одно содержимое пакета.
type PackageContent struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
