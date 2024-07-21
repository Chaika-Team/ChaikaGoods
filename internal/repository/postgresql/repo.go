package postgresql

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/myerr"
	"context"
	"errors"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v5"
)

type GoodsRepository struct {
	client Client
	log    log.Logger
}

// NewGoodsRepository создает новый экземпляр GoodsRepository.
func NewGoodsRepository(client Client, logger log.Logger) *GoodsRepository {
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

// AddQueryToCreateProduct добавляет запрос на создание продукта в базе данных.
func (r *GoodsRepository) AddQueryToCreateProduct(ctx context.Context, product models.Product) error {
	// добавляем новое изменение в базу
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ( $1, $2) RETURNING change_id;` //version_id подставляется автоматически
	var changeID int
	err := r.client.QueryRow(ctx, sql, models.OperationTypeInsert, product).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return err

	}
	return nil
}

// AddQueryToUpdateProduct добавление запроса на обновления продукта в базе данных.
func (r *GoodsRepository) AddQueryToUpdateProduct(ctx context.Context, product models.Product) error {
	// добавляем новое изменение в базу
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ( $1, $2) RETURNING change_id;` //version_id подставляется автоматически
	//TODO: в new_value ранится вся инфа о продукте, а не только измененные поля. Надо улучшить, чтобы хранились только измененные поля
	var changeID int
	err := r.client.QueryRow(ctx, sql, models.OperationTypeUpdate, product).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return err
	}
	return nil
}

// AddQueryToDeleteProduct удаляет продукт из базы данных по его ID.
func (r *GoodsRepository) AddQueryToDeleteProduct(ctx context.Context, id int64) error {
	sql := `INSERT INTO public.changes(operation, new_value) VALUES ($1, $2) RETURNING change_id;` //version_id подставляется автоматически
	var changeID int
	err := r.client.QueryRow(ctx, sql, models.OperationTypeDelete, id).Scan(&changeID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to add change: %v", err))
		return err
	}
	return nil
}

// ApplyChanges применяет изменения в базе данных продуктов.
func (r *GoodsRepository) ApplyChanges(ctx context.Context, version int) error {
	// берём все изменения, которые не были применены
	sql := `SELECT change_id, new_value, operation FROM public.changes WHERE version_id = $1 AND considered = FALSE;`
	rows, err := r.client.Query(ctx, sql, version)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get changes: %v", err))
		return err
	}
	defer rows.Close()
	// применяем изменения
	for rows.Next() {
		var changeID int
		var product models.Product
		var operation models.OperationType
		if err := rows.Scan(&changeID, &product, operation); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan change: %v", err))
			continue
		}
		// применяем изменение в зависимости от операции
		switch operation := operation; operation {
		case models.OperationTypeInsert:
			sql = `INSERT INTO public.product (name, description, price, imageurl, sku) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
			err = r.client.QueryRow(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).Scan(&product.ID)
			if err != nil {
				_ = r.log.Log("error", fmt.Sprintf("Failed to insert product: %v", err))
				continue
			}
		case models.OperationTypeUpdate:
			sql = `UPDATE public.product SET name = $1, description = $2, price = $3, imageurl = $4, sku = $5 WHERE id = $6;`
			_, err = r.client.Exec(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID)
			if err != nil {
				_ = r.log.Log("error", fmt.Sprintf("Failed to update product: %v", err))
				continue
			}
		case models.OperationTypeDelete:
			sql = `DELETE FROM public.product WHERE id = $1;`
			_, err = r.client.Exec(ctx, sql, product.ID)
			if err != nil {
				_ = r.log.Log("error", fmt.Sprintf("Failed to delete product: %v", err))
				continue
			}
		default:
			_ = r.log.Log("error", fmt.Sprintf("Unknown operation: %v", operation))
		}

		// помечаем изменение как примененное
		sql = `UPDATE public.changes SET considered = TRUE WHERE change_id = $1;`
		_, err = r.client.Exec(ctx, sql, changeID)
		if err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to update change: %v", err))
			continue
		}
	}
	return nil
}

// GetAllChanges возвращает все изменения в базе данных продуктов за конкретную версию.
func (r *GoodsRepository) GetAllChanges(ctx context.Context, version models.Version) ([]models.Change, error) {
	sql := `SELECT change_id, operation, new_value, change_timestamp, considered FROM public.changes WHERE version_id = $1;`
	rows, err := r.client.Query(ctx, sql, version.VersionID)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get all changes: %v", err))
		return nil, err
	}
	defer rows.Close()
	var changes []models.Change
	for rows.Next() {
		var c models.Change
		if err := rows.Scan(&c.ID, &c.OperationType, &c.NewValue, &c.ChangeTimestamp, &c.Considered); err != nil {
			_ = r.log.Log("error", fmt.Sprintf("Failed to scan change: %v", err))
			continue
		}
		c.VersionID = version.VersionID
		changes = append(changes, c)
	}
	return changes, nil
}

// GetCurrentDevVersion возвращает текущую версию базы данных продуктов к которой привязываются новые изменения.
func (r *GoodsRepository) GetCurrentDevVersion(ctx context.Context) (models.Version, error) {
	sql := `SELECT version_id, creation_date, is_dev, applied FROM public.version WHERE is_dev = TRUE ORDER BY creation_date DESC LIMIT 1;`
	var v models.Version
	err := r.client.QueryRow(ctx, sql).Scan(&v.VersionID, &v.CreationDate, &v.IsDev, &v.Applied)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to get current dev version: %v", err))
		return models.Version{}, err
	}
	return v, nil
}

// DeleteChange удаляет изменение из базы данных по его ID.
func (r *GoodsRepository) DeleteChange(ctx context.Context, id int64) error {
	sql := `DELETE FROM public.changes WHERE change_id = $1;`
	_, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		_ = r.log.Log("error", fmt.Sprintf("Failed to delete change: %v", err))
		return err
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, &myerr.NotFound{ID: fmt.Sprintf("%d", packageID)}
		}
		_ = r.log.Log("error", fmt.Sprintf("Failed to get package by ID: %v", err))
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

	if len(contents) == 0 {
		return nil, &myerr.NotFound{ID: fmt.Sprintf("%d", packageID)}
	}

	return contents, nil
}

// ListPackages возвращает список всех пакетов.
func (r *GoodsRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
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
