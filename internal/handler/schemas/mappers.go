package schemas

import (
	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

// ProductMapper содержит методы для преобразования между models.Product и schemas.ProductSchema.
type ProductMapper struct{}

// ToSchema преобразует models.Product в schemas.ProductSchema.
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
		// SKU отсутствует в schemas.ProductSchema, его можно задать отдельно при необходимости
	}
}

// PackageMapper содержит методы для преобразования между models.Package и schemas.PackageSchema.
type PackageMapper struct {
	ContentMapper *PackageContentMapper
	ProductMapper *ProductMapper
}

// ToSchema преобразует models.Package в schemas.PackageSchema.
func (pm *PackageMapper) ToSchema(pkg models.Package) PackageSchema {
	contentSchemas := pm.ContentMapper.ToSchemas(pkg.Content)
	return PackageSchema{
		ID:          pkg.ID,
		PackageName: pkg.PackageName,
		Description: pkg.Description,
		Content:     contentSchemas,
	}
}

// ToModel преобразует schemas.PackageSchema в models.Package.
func (pm *PackageMapper) ToModel(pkgSchema PackageSchema) models.Package {
	contentModels := pm.ContentMapper.ToModels(pkgSchema.Content)
	return models.Package{
		ID:          pkgSchema.ID,
		PackageName: pkgSchema.PackageName,
		Description: pkgSchema.Description,
		Content:     contentModels,
	}
}

// PackageContentMapper содержит методы для преобразования между models.PackageContent и schemas.PackageContentSchema.
type PackageContentMapper struct{}

// ToSchema преобразует models.PackageContent в schemas.PackageContentSchema.
func (pcm *PackageContentMapper) ToSchema(content models.PackageContent) PackageContentSchema {
	return PackageContentSchema{
		ProductID: content.ProductID,
		Quantity:  content.Quantity,
	}
}

// ToModel преобразует schemas.PackageContentSchema в models.PackageContent.
func (pcm *PackageContentMapper) ToModel(contentSchema PackageContentSchema) models.PackageContent {
	return models.PackageContent{
		ProductID: contentSchema.ProductID,
		Quantity:  contentSchema.Quantity,
	}
}

// ToSchemas преобразует слайс models.PackageContent в слайс schemas.PackageContentSchema.
func (pcm *PackageContentMapper) ToSchemas(contents []models.PackageContent) []PackageContentSchema {
	schemasList := make([]PackageContentSchema, len(contents))
	for i, content := range contents {
		schemasList[i] = pcm.ToSchema(content)
	}
	return schemasList
}

// ToModels преобразует слайс schemas.PackageContentSchema в слайс models.PackageContent.
func (pcm *PackageContentMapper) ToModels(contentSchemas []PackageContentSchema) []models.PackageContent {
	modelsList := make([]models.PackageContent, len(contentSchemas))
	for i, contentSchema := range contentSchemas {
		modelsList[i] = pcm.ToModel(contentSchema)
	}
	return modelsList
}

// Мапперы для списков продуктов

// ProductsMapper содержит методы для преобразования списков продуктов.
type ProductsMapper struct {
	ProductMapper *ProductMapper
}

// ToSchemas преобразует слайс models.Product в слайс schemas.ProductSchema.
func (pm *ProductsMapper) ToSchemas(products []models.Product) []ProductSchema {
	schemasList := make([]ProductSchema, len(products))
	for i, product := range products {
		schemasList[i] = pm.ProductMapper.ToSchema(product)
	}
	return schemasList
}

// ToModels преобразует слайс schemas.ProductSchema в слайс models.Product.
func (pm *ProductsMapper) ToModels(productSchemas []ProductSchema) []models.Product {
	modelsList := make([]models.Product, len(productSchemas))
	for i, productSchema := range productSchemas {
		modelsList[i] = pm.ProductMapper.ToModel(productSchema)
	}
	return modelsList
}

// Мапперы для списков пакетов

// PackagesMapper содержит методы для преобразования списков пакетов.
type PackagesMapper struct {
	PackageMapper *PackageMapper
}

// ToSchemas преобразует слайс models.Package в слайс schemas.PackageSchema.
func (pm *PackagesMapper) ToSchemas(packages []models.Package) []PackageSchema {
	schemasList := make([]PackageSchema, len(packages))
	for i, pkg := range packages {
		schemasList[i] = pm.PackageMapper.ToSchema(pkg)
	}
	return schemasList
}

// ToModels преобразует слайс schemas.PackageSchema в слайс models.Package.
func (pm *PackagesMapper) ToModels(packageSchemas []PackageSchema) []models.Package {
	modelsList := make([]models.Package, len(packageSchemas))
	for i, pkgSchema := range packageSchemas {
		modelsList[i] = pm.PackageMapper.ToModel(pkgSchema)
	}
	return modelsList
}
