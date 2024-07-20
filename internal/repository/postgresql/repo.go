package postgresql

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/myerr"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GoodsRepository struct {
	client *pgxpool.Pool
	log    log.Logger
}

// NewGoodsRepository создает новый экземпляр GoodsRepository.
func NewGoodsRepository(client *pgxpool.Pool, logger log.Logger) *GoodsRepository {
	return &GoodsRepository{
		client: client,
		log:    logger,
	}
}

// GetProductByID возвращает продукт по его ID.
func (r *GoodsRepository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product WHERE id = $1;`
	row := r.client.QueryRow(ctx, sql, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, &myerr.NotFound{ID: fmt.Sprintf("%d", id)}
	} else if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get product by ID: %v", err))
		return nil, err
	}

	return &p, nil
}

// GetAllProducts возвращает список всех продуктов.
func (r *GoodsRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product;`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get all products: %v", err))
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan product: %v", err))
			continue
		}
		products = append(products, p)
	}
	if err = rows.Err(); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed during rows iteration: %v", err))
		return nil, err
	}
	return products, nil
}

// CreateProduct добавляет новый продукт в базу данных.
func (r *GoodsRepository) CreateProduct(ctx context.Context, product *models.Product) error {
	sql := `INSERT INTO public.product (name, description, price, imageurl, sku) VALUES ($1, $2, $3, $4, $5) ON CONFLICT (sku) DO NOTHING RETURNING id;`
	err := r.client.QueryRow(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).Scan(&product.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		return &myerr.DuplicateError{Str: product.SKU}
	} else if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create product: %v", err))
		return err
	}
	return nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (r *GoodsRepository) UpdateProduct(ctx context.Context, product *models.Product) error {
	sql := `UPDATE public.product SET name = $1, description = $2, price = $3, imageurl = $4, sku = $5 WHERE id = $6;`
	commandTag, err := r.client.Exec(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to update product: %v", err))
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return &myerr.NotFound{ID: fmt.Sprintf("%d", product.ID)}
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных по его ID.
func (r *GoodsRepository) DeleteProduct(ctx context.Context, id int64) error {
	sql := `DELETE FROM public.product WHERE id = $1;`
	commandTag, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete product: %v", err))
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return &myerr.NotFound{ID: fmt.Sprintf("%d", id)}
	}
	return nil
}

// Package queries

// GetPackageByID получает полную информацию о пакете, включая его содержимое.
func (r *GoodsRepository) GetPackageByID(ctx context.Context, packageID int64) (*models.Package, []models.PackageContent, error) {
	sqlPackage := `SELECT packageid, packagename, description FROM public."package" WHERE packageid = $1;`
	pkg := models.Package{}
	err := r.client.QueryRow(ctx, sqlPackage, packageID).Scan(&pkg.ID, &pkg.PackageName, &pkg.Description)
	if err != nil {
		err := r.log.Log("error", fmt.Sprintf("Failed to get package by ID: %v", err))
		if err != nil {
			return nil, nil, err
		}
		return nil, nil, err
	}

	sqlContents := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sqlContents, packageID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package contents: %v", err))
		return nil, nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ID, &c.PackageID, &c.ProductID, &c.Quantity); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package content: %v", err))
			continue
		}
		contents = append(contents, c)
	}

	return &pkg, contents, nil
}

// GetProductsByPackageID получает список продуктов в определенном пакете.
func (r *GoodsRepository) GetProductsByPackageID(ctx context.Context, packageID int64) ([]models.PackageContent, error) {
	sql := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sql, packageID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get products by package ID: %v", err))
		return nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ID, &c.PackageID, &c.ProductID, &c.Quantity); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package content: %v", err))
			continue
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// ListPackages возвращает список всех пакетов.
func (r *GoodsRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
	sql := `SELECT packageid, packagename, description FROM public."package";`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to list packages: %v", err))
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package: %v", err))
			continue
		}
		packages = append(packages, p)
	}

	return packages, nil
}

// CreatePackage добавляет новый пустой пакет в базу данных.
func (r *GoodsRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	sql := `INSERT INTO public."package" (packagename, description) VALUES ($1, $2) RETURNING packageid;`
	err := r.client.QueryRow(ctx, sql, pkg.PackageName, pkg.Description).Scan(&pkg.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create package: %v", err))
		return err
	}
	return nil
}

// AddProductToPackage добавляет продукты в пакет.
func (r *GoodsRepository) AddProductToPackage(ctx context.Context, packageID int64, products []models.PackageContent) error {
	sql := `INSERT INTO public.packagecontent (packageid, productid, quantity) VALUES ($1, $2, $3);`
	batch := &pgx.Batch{}

	for _, product := range products {
		batch.Queue(sql, packageID, product.ProductID, product.Quantity)
	}

	br := r.client.SendBatch(ctx, batch)
	defer func(br pgx.BatchResults) {
		err := br.Close()
		if err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to close batch results: %v", err))
		}
	}(br)
	_, err := br.Exec()
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add products to package: %v", err))
		return err
	}
	return nil
}

// DeletePackage удаляет пакет и его содержимое.
func (r *GoodsRepository) DeletePackage(ctx context.Context, packageID int64) error {
	sqlDeleteContents := `DELETE FROM public.packagecontent WHERE packageid = $1;`
	sqlDeletePackage := `DELETE FROM public."package" WHERE packageid = $1;`

	tx, err := r.client.Begin(ctx)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to begin transaction: %v", err))
		return err
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to rollback transaction: %v", err))
		}
	}(tx, ctx)

	if _, err := tx.Exec(ctx, sqlDeleteContents, packageID); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete package contents: %v", err))
		return err
	}
	if _, err := tx.Exec(ctx, sqlDeletePackage, packageID); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete package: %v", err))
		return err
	}

	if err := tx.Commit(ctx); err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to commit transaction: %v", err))
		return err
	}
	return nil
}
