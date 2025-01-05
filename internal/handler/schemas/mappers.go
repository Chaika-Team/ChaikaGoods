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

// PackageContentMapper реализует интерфейс Mapper для PackageContent.
type PackageContentMapper struct{}

func NewPackageContentMapper() *PackageContentMapper {
	return &PackageContentMapper{}
}

func (pcm *PackageContentMapper) ToSchema(content models.PackageContent) PackageContentSchema {
	return PackageContentSchema{
		ProductID: content.ProductID,
		Quantity:  content.Quantity,
	}
}

func (pcm *PackageContentMapper) ToModel(contentSchema PackageContentSchema) models.PackageContent {
	return models.PackageContent{
		ProductID: contentSchema.ProductID,
		Quantity:  contentSchema.Quantity,
	}
}

// PackageMapper реализует интерфейс Mapper для Package.
type PackageMapper struct {
	ContentMapper Mapper[models.PackageContent, PackageContentSchema]
	ProductMapper Mapper[models.Product, ProductSchema]
}

func NewPackageMapper(contentMapper Mapper[models.PackageContent, PackageContentSchema], productMapper Mapper[models.Product, ProductSchema]) *PackageMapper {
	return &PackageMapper{
		ContentMapper: contentMapper,
		ProductMapper: productMapper,
	}
}

func (pm *PackageMapper) ToSchema(pkg models.Package) PackageSchema {
	contentSchemas := make([]PackageContentSchema, len(pkg.Content))
	for i, content := range pkg.Content {
		contentSchemas[i] = pm.ContentMapper.ToSchema(content)
	}
	return PackageSchema{
		ID:          pkg.ID,
		PackageName: pkg.PackageName,
		Description: pkg.Description,
		Content:     contentSchemas,
	}
}

func (pm *PackageMapper) ToModel(pkgSchema PackageSchema) models.Package {
	contentModels := make([]models.PackageContent, len(pkgSchema.Content))
	for i, contentSchema := range pkgSchema.Content {
		contentModels[i] = pm.ContentMapper.ToModel(contentSchema)
	}
	return models.Package{
		ID:          pkgSchema.ID,
		PackageName: pkgSchema.PackageName,
		Description: pkgSchema.Description,
		Content:     contentModels,
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

// PackagesMapper реализует методы для работы с коллекциями пакетов.
type PackagesMapper struct {
	PackageMapper Mapper[models.Package, PackageSchema]
}

func NewPackagesMapper(pm Mapper[models.Package, PackageSchema]) *PackagesMapper {
	return &PackagesMapper{
		PackageMapper: pm,
	}
}

func (pm *PackagesMapper) ToSchemas(packages []models.Package) []PackageSchema {
	schemasList := make([]PackageSchema, len(packages))
	for i, pkg := range packages {
		schemasList[i] = pm.PackageMapper.ToSchema(pkg)
	}
	return schemasList
}

func (pm *PackagesMapper) ToModels(packageSchemas []PackageSchema) []models.Package {
	modelsList := make([]models.Package, len(packageSchemas))
	for i, pkgSchema := range packageSchemas {
		modelsList[i] = pm.PackageMapper.ToModel(pkgSchema)
	}
	return modelsList
}
