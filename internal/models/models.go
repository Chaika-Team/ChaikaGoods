package models

import (
	"encoding/json"
	"time"
)

// Product описывает товар.
type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"imageurl"`
	SKU         string  `json:"sku"`
}

// Add custom methods to MarshalJSON and UnmarshalJSON for the models.Product type.

// Package описывает пакет товаров.
type Package struct {
	ID          int64  `json:"id"`
	PackageName string `json:"package_name"`
	Description string `json:"description"`
}

// PackageContent описывает одно содержимое пакета.
type PackageContent struct {
	ID        int64 `json:"id"`
	PackageID int64 `json:"package_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

// Version описывает версию базы продуктов.
type Version struct {
	VersionID    int       `json:"version_id"`
	CreationDate time.Time `json:"creation_date"`
	IsDev        bool      `json:"is_dev"`
	Applied      bool      `json:"applied"`
}

// Change описывает изменение в базе продуктов в модели CQRS.
type Change struct {
	ID              int64           `json:"id"`
	VersionID       int             `json:"version_id"`
	OperationType   OperationType   `json:"operation_type"`
	NewValue        json.RawMessage `json:"new_value"`
	ChangeTimestamp time.Time       `json:"change_timestamp"`
	Considered      bool            `json:"considered"`
}
