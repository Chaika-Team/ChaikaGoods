package models

import (
	"database/sql"
	"encoding/json"
	"time"
)

// Product описывает товар.
type Product struct {
	ID          int64           `json:"id"`
	Name        sql.NullString  `json:"name"`
	Description sql.NullString  `json:"description"`
	Price       sql.NullFloat64 `json:"price"`
	ImageURL    sql.NullString  `json:"image_url"`
	SKU         sql.NullString  `json:"sku"`
}

// UnmarshalJSON реализует интерфейс Unmarshaler для типа Product.
func (p *Product) UnmarshalJSON(data []byte) error {
	type Alias Product
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(p),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	p.Name = sql.NullString{String: aux.Name.String, Valid: aux.Name.Valid}
	p.Description = sql.NullString{String: aux.Description.String, Valid: aux.Description.Valid}
	p.Price = sql.NullFloat64{Float64: aux.Price.Float64, Valid: aux.Price.Valid}
	p.ImageURL = sql.NullString{String: aux.ImageURL.String, Valid: aux.ImageURL.Valid}
	p.SKU = sql.NullString{String: aux.SKU.String, Valid: aux.SKU.Valid}
	return nil
}

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
