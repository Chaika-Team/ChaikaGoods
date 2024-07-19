package models

// Product описывает товар.
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	SKU         string  `json:"sku"`
}

// Package описывает пакет товаров.
type Package struct {
	ID          int64  `json:"id"`
	PackageName string `json:"package_name"`
	Description string `json:"description"`
}

// PackageContent описывает содержимое пакета.
type PackageContent struct {
	ID        int64 `json:"id"`
	PackageID int64 `json:"package_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}
