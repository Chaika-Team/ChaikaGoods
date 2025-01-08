package schemas

import (
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

// Mapper интерфейс для преобразования между моделью и схемой.
type Mapper[TModel any, TSchema any] interface {
	ToSchema(model TModel) TSchema
	ToModel(schema TSchema) TModel
}

// ProductMapper реализует интерфейс Mapper для Product.
type ProductMapper struct{}

func NewProductMapper() *ProductMapper {
	return &ProductMapper{}
}

func (pm *ProductMapper) ToSchema(product models.Product) ProductSchema {
	return ProductSchema{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		ImageURL:    product.ImageURL,
	}
}

// ToModel преобразует schemas.ProductSchema в models.Product.
func (pm *ProductMapper) ToModel(productSchema ProductSchema) models.Product {
	return models.Product{
		ID:          productSchema.ID,
		Name:        productSchema.Name,
		Description: productSchema.Description,
		Price:       productSchema.Price,
		ImageURL:    productSchema.ImageURL,
	}
}

// TemplateContentMapper реализует интерфейс Mapper для TemplateContent.
type TemplateContentMapper struct{}

func NewTemplateContentMapper() *TemplateContentMapper {
	return &TemplateContentMapper{}
}

func (pcm *TemplateContentMapper) ToSchema(content models.TemplateContent) TemplateContentSchema {
	return TemplateContentSchema{
		ProductID: content.ProductID,
		Quantity:  content.Quantity,
	}
}

func (pcm *TemplateContentMapper) ToModel(contentSchema TemplateContentSchema) models.TemplateContent {
	return models.TemplateContent{
		ProductID: contentSchema.ProductID,
		Quantity:  contentSchema.Quantity,
	}
}

// TemplateMapper реализует интерфейс Mapper для Template.
type TemplateMapper struct {
	ContentMapper Mapper[models.TemplateContent, TemplateContentSchema]
	ProductMapper Mapper[models.Product, ProductSchema]
}

func NewTemplateMapper(contentMapper Mapper[models.TemplateContent, TemplateContentSchema], productMapper Mapper[models.Product, ProductSchema]) *TemplateMapper {
	return &TemplateMapper{
		ContentMapper: contentMapper,
		ProductMapper: productMapper,
	}
}

func (pm *TemplateMapper) ToSchema(pkg models.Template) TemplateSchema {
	contentSchemas := make([]TemplateContentSchema, len(pkg.Content))
	for i, content := range pkg.Content {
		contentSchemas[i] = pm.ContentMapper.ToSchema(content)
	}
	return TemplateSchema{
		ID:           pkg.ID,
		TemplateName: pkg.TemplateName,
		Description:  pkg.Description,
		Content:      contentSchemas,
	}
}

func (pm *TemplateMapper) ToModel(pkgSchema TemplateSchema) models.Template {
	contentModels := make([]models.TemplateContent, len(pkgSchema.Content))
	for i, contentSchema := range pkgSchema.Content {
		contentModels[i] = pm.ContentMapper.ToModel(contentSchema)
	}
	return models.Template{
		ID:           pkgSchema.ID,
		TemplateName: pkgSchema.TemplateName,
		Description:  pkgSchema.Description,
		Content:      contentModels,
	}
}

// ProductsMapper реализует методы для работы с коллекциями продуктов.
type ProductsMapper struct {
	ProductMapper Mapper[models.Product, ProductSchema]
}

func NewProductsMapper(pm Mapper[models.Product, ProductSchema]) *ProductsMapper {
	return &ProductsMapper{
		ProductMapper: pm,
	}
}

func (pm *ProductsMapper) ToSchemas(products []models.Product) []ProductSchema {
	schemasList := make([]ProductSchema, len(products))
	for i, product := range products {
		schemasList[i] = pm.ProductMapper.ToSchema(product)
	}
	return schemasList
}

func (pm *ProductsMapper) ToModels(productSchemas []ProductSchema) []models.Product {
	modelsList := make([]models.Product, len(productSchemas))
	for i, productSchema := range productSchemas {
		modelsList[i] = pm.ProductMapper.ToModel(productSchema)
	}
	return modelsList
}

// TemplatesMapper реализует методы для работы с коллекциями пакетов.
type TemplatesMapper struct {
	TemplateMapper Mapper[models.Template, TemplateSchema]
}

func NewTemplatesMapper(pm Mapper[models.Template, TemplateSchema]) *TemplatesMapper {
	return &TemplatesMapper{
		TemplateMapper: pm,
	}
}

func (pm *TemplatesMapper) ToSchemas(templates []models.Template) []TemplateSchema {
	schemasList := make([]TemplateSchema, len(templates))
	for i, pkg := range templates {
		schemasList[i] = pm.TemplateMapper.ToSchema(pkg)
	}
	return schemasList
}

func (pm *TemplatesMapper) ToModels(templateSchemas []TemplateSchema) []models.Template {
	modelsList := make([]models.Template, len(templateSchemas))
	for i, pkgSchema := range templateSchemas {
		modelsList[i] = pm.TemplateMapper.ToModel(pkgSchema)
	}
	return modelsList
}
