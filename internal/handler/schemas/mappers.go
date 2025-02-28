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

func (tcm *TemplateContentMapper) ToSchema(content models.TemplateContent) TemplateContentSchema {
	return TemplateContentSchema{
		ProductID: content.ProductID,
		Quantity:  content.Quantity,
	}
}

func (tcm *TemplateContentMapper) ToModel(contentSchema TemplateContentSchema) models.TemplateContent {
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

func (tm *TemplateMapper) ToSchema(template models.Template) TemplateSchema {
	contentSchemas := make([]TemplateContentSchema, len(template.Content))
	for i, content := range template.Content {
		contentSchemas[i] = tm.ContentMapper.ToSchema(content)
	}
	return TemplateSchema{
		ID:           template.ID,
		TemplateName: template.TemplateName,
		Description:  template.Description,
		Content:      contentSchemas,
	}
}

func (tm *TemplateMapper) ToModel(templateSchema TemplateSchema) models.Template {
	contentModels := make([]models.TemplateContent, len(templateSchema.Content))
	for i, contentSchema := range templateSchema.Content {
		contentModels[i] = tm.ContentMapper.ToModel(contentSchema)
	}
	return models.Template{
		ID:           templateSchema.ID,
		TemplateName: templateSchema.TemplateName,
		Description:  templateSchema.Description,
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

// TemplatesMapper реализует методы для работы с коллекциями шаблонов.
type TemplatesMapper struct {
	TemplateMapper Mapper[models.Template, TemplateSchema]
}

func NewTemplatesMapper(tm Mapper[models.Template, TemplateSchema]) *TemplatesMapper {
	return &TemplatesMapper{
		TemplateMapper: tm,
	}
}

func (tm *TemplatesMapper) ToSchemas(templates []models.Template) []TemplateSchema {
	schemasList := make([]TemplateSchema, len(templates))
	for i, template := range templates {
		schemasList[i] = tm.TemplateMapper.ToSchema(template)
	}
	return schemasList
}

func (tm *TemplatesMapper) ToModels(templateSchemas []TemplateSchema) []models.Template {
	modelsList := make([]models.Template, len(templateSchemas))
	for i, templateSchema := range templateSchemas {
		modelsList[i] = tm.TemplateMapper.ToModel(templateSchema)
	}
	return modelsList
}
