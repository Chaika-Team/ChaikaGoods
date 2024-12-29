package postgresql

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/myerr"
	"ChaikaGoods/internal/repository"
	"ChaikaGoods/internal/utils"
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/go-kit/log"
	"github.com/jackc/pgx/v5"
)

type GoodsPGRepository struct {
	client Client
	log    log.Logger
}

// NewGoodsRepository создает новый экземпляр GoodsRepository.
func NewGoodsRepository(client Client, logger log.Logger) repository.GoodsRepository {
	return &GoodsPGRepository{
		client: client,
		log:    logger,
	}
}

// GetProductByID возвращает продукт по его ID.
func (r *GoodsPGRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product WHERE id = $1;`
	row := r.client.QueryRow(ctx, sql, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU)
	if errors.Is(err, pgx.ErrNoRows) {
		return p, myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Product with id %d not found", id), nil)
	} else if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get product by ID: %v", err))
		return p, err
	}

	return p, nil
}

// GetAllProducts возвращает список всех продуктов.
func (r *GoodsPGRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
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
		_ = r.log.Log("error", fmt.Sprintf("Failed to get all products: %v", err))
		return nil, err
	}
	return products, nil
}

// CreateProduct создание продукта в базе данных.
func (r *GoodsPGRepository) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	sql := `INSERT INTO public.product (name, description, price, imageurl, sku) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := r.client.QueryRow(ctx, sql, p.Name, p.Description, p.Price, p.ImageURL, p.SKU).Scan(&p.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create product: %v", err))
		return 0, err
	}
	return p.ID, nil

}

// UpdateProduct обновление продукта в базе данных.
func (r *GoodsPGRepository) UpdateProduct(ctx context.Context, p *models.Product) error {
	sql := `UPDATE public.product SET name = $1, description = $2, price = $3, imageurl = $4, sku = $5 WHERE id = $6;`
	_, err := r.client.Exec(ctx, sql, p.Name, p.Description, p.Price, p.ImageURL, p.SKU, p.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Product with id %d not found", p.ID), nil)
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to update product: %v", err))
		return err
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных по его ID.
func (r *GoodsPGRepository) DeleteProduct(ctx context.Context, id int64) error {
	sql := `DELETE FROM public.product WHERE id = $1;`
	_, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Product with id %d not found", id), nil)
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete product: %v", err))
		return err
	}
	return nil
}

// Package queries

// GetPackageByID получает полную информацию о пакете, включая его содержимое.
func (r *GoodsPGRepository) GetPackageByID(ctx context.Context, p *models.Package) ([]models.PackageContent, error) {
	sqlPackage := `SELECT packageid, packagename, description FROM public."package" WHERE packageid = $1;`
	err := r.client.QueryRow(ctx, sqlPackage, p.ID).Scan(&p.ID, &p.PackageName, &p.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Package with id %d not found", p.ID), nil)
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package by ID: %v", err))
		return nil, err
	}

	sqlContents := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sqlContents, p.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package contents: %v", err))
		return nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ID, &c.PackageID, &c.ProductID, &c.Quantity); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan package content: %v", err))
			return nil, err
		}
		contents = append(contents, c)
	}

	return contents, nil
}

// GetProductsByPackageID получает список продуктов в определенном пакете.
func (r *GoodsPGRepository) GetProductsByPackageID(ctx context.Context, p *models.Package) ([]models.PackageContent, error) {
	sql := `SELECT packagecontentid, packageid, productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sql, p.ID)
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
func (r *GoodsPGRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
	sql := `SELECT packageid, packagename, description FROM public.package;`
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
func (r *GoodsPGRepository) CreatePackage(ctx context.Context, pkg *models.Package) error {
	sql := `INSERT INTO public."package" (packagename, description) VALUES ($1, $2) RETURNING packageid;`
	err := r.client.QueryRow(ctx, sql, &pkg.PackageName, &pkg.Description).Scan(&pkg.ID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to create package: %v", err))
		return err
	}
	return nil
}

// AddProductToPackage добавляет продукты в пакет.
func (r *GoodsPGRepository) AddProductToPackage(ctx context.Context, packageID int64, products []models.PackageContent) error {
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
func (r *GoodsPGRepository) DeletePackage(ctx context.Context, packageID int64) error {
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
		// Check if not found
		if errors.Is(err, pgx.ErrNoRows) {
			return myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Package content with id %d not found", packageID), nil)
		}
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

func (r *GoodsPGRepository) SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error) {
	sql := `SELECT packageid, packagename, description FROM public."package" WHERE packagename LIKE $1 OR description LIKE $1 LIMIT $2 OFFSET $3;`
	rows, err := r.client.Query(ctx, sql, searchString, quantity, offset)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, myerr.NewAppError(myerr.ErrorTypeNotFound, fmt.Sprintf("Package with name %s not found", searchString), nil)
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to search package: %v", err))
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

// insertProduct добавляет новый продукт в базу данных, только для внутреннего использования.
func (r *GoodsPGRepository) insertProduct(ctx context.Context, data *map[string]interface{}) error {
	// Verify that all keys in the map correspond to the fields of the models.Product structure
	err := utils.VerifyMapFields[models.Product](*data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to verify product fields: %v", err))
		return err
	}
	// Generate sql query with all fields in the map
	sql := `INSERT INTO public.product (`
	values := `) VALUES (`
	var args []interface{}
	idx := 1
	for key, value := range *data {
		sql += key + ", "
		values += "$" + strconv.Itoa(idx) + ", "
		args = append(args, value) // Добавление значения в список аргументов
		idx++
	}
	sql = sql[:len(sql)-2] + values[:len(values)-2] + ");"
	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to insert product: %v", err)
	}
	return nil
}

// updateProduct обновляет информацию о продукте в базе данных, только для внутреннего использования.
func (r *GoodsPGRepository) updateProduct(ctx context.Context, data *map[string]interface{}) error {
	// Verify that all keys in the map correspond to the fields of the models.Product structure
	err := utils.VerifyMapFields[models.Product](*data)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to verify product fields: %v", err))
		return err
	}
	// Generate sql query with all fields in the map
	sql := `UPDATE public.product SET `
	var args []interface{}
	idx := 1
	for key, value := range *data {
		if key == "id" {
			continue
		}
		sql += key + " = $" + strconv.Itoa(idx) + ", "
		args = append(args, value) // Добавление значения в список аргументов
		idx++
	}
	if len(args) == 0 {
		return fmt.Errorf("no fields to update")
	}
	args = append(args, (*data)["id"]) // Добавление id в список аргументов
	sql = sql[:len(sql)-2] + " WHERE id = $" + strconv.Itoa(idx) + ";"
	_, err = r.client.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to update product: %v", err)
	}
	return nil
}

// deleteProduct удаляет продукт из базы данных по его ID, только для внутреннего использования.
func (r *GoodsPGRepository) deleteProduct(ctx context.Context, data *map[string]interface{}) error {
	// get id
	rawId, ok := (*data)["id"]
	if !ok {
		return myerr.NewAppError(myerr.ErrorTypeValidation, "id is required", nil)
	}

	var id = int64(rawId.(float64))

	sql := `DELETE FROM public.product WHERE id = $1;`
	_, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %v", err)
	}
	return nil
}
