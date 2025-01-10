package models

import (
	"context"
)

// ProductRepository defines methods for product-related database operations.
type ProductRepository interface {
	GetProductByID(ctx context.Context, id int64) (Product, error)
	GetAllProducts(ctx context.Context) ([]Product, error)
	CreateProduct(ctx context.Context, p *Product) (int64, error)
	UpdateProduct(ctx context.Context, p *Product) error
	DeleteProduct(ctx context.Context, id int64) error
}

// TemplateRepository defines methods for template-related database operations.
type TemplateRepository interface {
	GetTemplateByID(ctx context.Context, id int64) (Template, error)
	GetProductsByTemplateID(ctx context.Context, templateID int64) ([]TemplateContent, error)
	ListTemplates(ctx context.Context) ([]Template, error)
	CreateTemplate(ctx context.Context, template *Template) error
	DeleteTemplate(ctx context.Context, templateID int64) error
	SearchTemplates(ctx context.Context, searchString string, limit int64, offset int64) ([]Template, error)
	GetAllTemplates(ctx context.Context, limit int64, offset int64) ([]Template, error)
}

// GoodsRepository объединяет репозитории для продуктов и шаблонов.
type GoodsRepository interface {
	ProductRepository
	TemplateRepository
}
