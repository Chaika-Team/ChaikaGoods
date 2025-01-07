package postgresql

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/go-kit/log"
	"github.com/jackc/pgx/v5"
)

const (
	msgFailedToScanPackage = "Failed to scan package"
	fmtProductNotFound     = "Product with ID %d not found"
)

// GoodsPGRepository implements the GoodsRepository interface using PostgreSQL.
type GoodsPGRepository struct {
	client Client
	logger log.Logger
}

// NewGoodsRepository creates a new instance of GoodsRepository.
func NewGoodsRepository(client Client, logger log.Logger) models.GoodsRepository {
	return &GoodsPGRepository{
		client: client,
		logger: logger,
	}
}

// ---------- ProductRepository Implementation ----------

// GetProductByID returns a product by its ID.
func (r *GoodsPGRepository) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	const sql = `SELECT id, name, description, price, imageurl, sku FROM public.product WHERE id = $1;`
	row := r.client.QueryRow(ctx, sql, id)

	var p models.Product
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return p, myerr.NotFound(fmt.Sprintf(fmtProductNotFound, id), nil)
		}
		_ = r.logger.Log("error", "Failed to get product by ID", "id", id, "err", err)
		return p, err
	}

	return p, nil
}

// GetAllProducts returns a list of all products.
func (r *GoodsPGRepository) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	const sql = `SELECT id, name, description, price, imageurl, sku FROM public.product;`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.logger.Log("error", "Failed to retrieve all products", "err", err)
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU); err != nil {
			_ = r.logger.Log("error", "Failed to scan product", "err", err)
			continue // Пропускаем некорректную строку, но продолжаем обработку остальных.
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Error during rows iteration", "err", err)
		return nil, err
	}

	return products, nil
}

// CreateProduct creates a new product in the database.
func (r *GoodsPGRepository) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	const sql = `INSERT INTO public.product (name, description, price, imageurl, sku) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	if err := r.client.QueryRow(ctx, sql, p.Name, p.Description, p.Price, p.ImageURL, p.SKU).Scan(&p.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, myerr.Conflict(fmt.Sprintf("Product with SKU %s already exists", p.SKU), err)
		}
		_ = r.logger.Log("error", "Failed to create product", "err", err)
		return 0, err
	}
	return p.ID, nil
}

// UpdateProduct updates an existing product in the database.
func (r *GoodsPGRepository) UpdateProduct(ctx context.Context, p *models.Product) error {
	const sql = `UPDATE public.product SET name = $1, description = $2, price = $3, imageurl = $4, sku = $5 WHERE id = $6;`
	ct, err := r.client.Exec(ctx, sql, p.Name, p.Description, p.Price, p.ImageURL, p.SKU, p.ID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return myerr.Conflict(fmt.Sprintf("Updated data conflicts with existing product with SKU %s", p.SKU), err)
		}
		_ = r.logger.Log("error", "Failed to update product", "id", p.ID, "err", err)
		return err
	}
	if ct.RowsAffected() == 0 {
		return myerr.NotFound(fmt.Sprintf(fmtProductNotFound, p.ID), nil)
	}
	return nil
}

// DeleteProduct deletes a product from the database by its ID.
func (r *GoodsPGRepository) DeleteProduct(ctx context.Context, id int64) error {
	const sql = `DELETE FROM public.product WHERE id = $1;`
	ct, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return myerr.NotFound(fmt.Sprintf(fmtProductNotFound, id), nil)
		}
		_ = r.logger.Log("error", "Failed to delete product", "id", id, "err", err)
		return err
	}
	if ct.RowsAffected() == 0 {
		return myerr.NotFound(fmt.Sprintf(fmtProductNotFound, id), nil)
	}
	return nil
}

// ---------- PackageRepository Implementation ----------

// GetPackageByID retrieves package details along with its contents.
func (r *GoodsPGRepository) GetPackageByID(ctx context.Context, id int64) (models.Package, error) {
	var pkg models.Package
	const sqlPackage = `SELECT packageid, packagename, description FROM public."package" WHERE packageid = $1;`
	if err := r.client.QueryRow(ctx, sqlPackage, id).Scan(&pkg.ID, &pkg.PackageName, &pkg.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pkg, myerr.NotFound(fmt.Sprintf("Package with ID %d not found", id), nil)
		}
		_ = r.logger.Log("error", "Failed to get package by ID", "id", id, "err", err)
		return pkg, err
	}

	// Get package contents
	const sqlContents = `SELECT productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sqlContents, id)
	if err != nil {
		_ = r.logger.Log("error", "Failed to get package contents", "package_id", id, "err", err)
		return pkg, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ProductID, &c.Quantity); err != nil {
			_ = r.logger.Log("error", "Failed to scan package content", "err", err)
			return pkg, err
		}
		pkg.Content = append(pkg.Content, c)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting package contents", "err", err)
		return pkg, err
	}

	return pkg, nil
}

// GetProductsByPackageID retrieves all products within a specific package.
func (r *GoodsPGRepository) GetProductsByPackageID(ctx context.Context, packageID int64) ([]models.PackageContent, error) {
	const sql = `SELECT productid, quantity FROM public.packagecontent WHERE packageid = $1;`
	rows, err := r.client.Query(ctx, sql, packageID)
	if err != nil {
		_ = r.logger.Log("error", "Failed to get products by package ID", "package_id", packageID, "err", err)
		return nil, err
	}
	defer rows.Close()

	var contents []models.PackageContent
	for rows.Next() {
		var c models.PackageContent
		if err := rows.Scan(&c.ProductID, &c.Quantity); err != nil {
			_ = r.logger.Log("error", "Failed to scan package content", "err", err)
			return nil, err
		}
		contents = append(contents, c)
	}

	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting products by package ID", "err", err)
		return nil, err
	}

	return contents, nil
}

// ListPackages returns a list of all packages.
func (r *GoodsPGRepository) ListPackages(ctx context.Context) ([]models.Package, error) {
	const sql = `SELECT packageid, packagename, description FROM public."package";`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.logger.Log("error", "Failed to list packages", "err", err)
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanPackage, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		packages = append(packages, p)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while listing packages", "err", err)
		return nil, err
	}

	return packages, nil
}

// CreatePackage adds a new package to the database along with its contents.
func (r *GoodsPGRepository) CreatePackage(ctx context.Context, pkg *models.Package) (err error) {
	const sqlInsertPackage = `INSERT INTO public."package" (packagename, description) VALUES ($1, $2) RETURNING packageid;`

	tx, err := r.client.Begin(ctx)
	if err != nil {
		_ = r.logger.Log("error", "Failed to begin transaction", "err", err)
		return err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback(ctx)
			panic(p)
		} else if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			err = tx.Commit(ctx)
			if err != nil {
				_ = r.logger.Log("error", "Failed to commit transaction", "err", err)
			}
		}
	}()

	// Insert package
	if err = tx.QueryRow(ctx, sqlInsertPackage, pkg.PackageName, pkg.Description).Scan(&pkg.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return myerr.Conflict(fmt.Sprintf("Package with name %s already exists", pkg.PackageName), err)
		}
		_ = r.logger.Log("error", "Failed to insert package", "pkg_name", pkg.PackageName, "err", err)
		return err
	}

	// Insert package contents
	for _, content := range pkg.Content {
		if err = r.createProductToPackage(ctx, tx, pkg.ID, content); err != nil {
			_ = r.logger.Log("error", "Failed to add product to package", "pkg_id", pkg.ID, "product_id", content.ProductID, "err", err)
			return err
		}
	}

	return nil
}

// createProductToPackage adds a single package content entry.
func (r *GoodsPGRepository) createProductToPackage(ctx context.Context, tx pgx.Tx, packageID int64, content models.PackageContent) error {
	// Insert package content
	const sqlInsertContent = `INSERT INTO public.packagecontent (packageid, productid, quantity) VALUES ($1, $2, $3);`
	if _, err := tx.Exec(ctx, sqlInsertContent, packageID, content.ProductID, content.Quantity); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return myerr.NotFound(fmt.Sprintf(fmtProductNotFound, content.ProductID), err)
		} else if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return myerr.Conflict(fmt.Sprintf("Product with ID %d already exists in package", content.ProductID), err)
		}
		_ = r.logger.Log("error", "Failed to insert package content", "package_id", packageID, "product_id", content.ProductID, "err", err)
		return err
	}
	return nil
}

// DeletePackage deletes a package and its contents from the database by package ID.
func (r *GoodsPGRepository) DeletePackage(ctx context.Context, packageID int64) error {
	const (
		sqlDeleteContents = `DELETE FROM public.packagecontent WHERE packageid = $1;`
		sqlDeletePackage  = `DELETE FROM public."package" WHERE packageid = $1;`
	)

	tx, err := r.client.Begin(ctx)
	if err != nil {
		_ = r.logger.Log("error", "Failed to begin transaction", "err", err)
		return err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// Delete package contents
	if _, err = tx.Exec(ctx, sqlDeleteContents, packageID); err != nil {
		_ = r.logger.Log("error", "Failed to delete package contents", "package_id", packageID, "err", err)
		return err
	}

	// Delete package
	ct, err := tx.Exec(ctx, sqlDeletePackage, packageID)
	if err != nil {
		_ = r.logger.Log("error", "Failed to delete package", "package_id", packageID, "err", err)
		return err
	}
	if ct.RowsAffected() == 0 {
		return myerr.NotFound(fmt.Sprintf("Package with id %d not found", packageID), nil)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		_ = r.logger.Log("error", "Failed to commit transaction", "err", err)
		return err
	}

	return nil
}

// SearchPackages searches for packages by name or description with pagination.
func (r *GoodsPGRepository) SearchPackages(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Package, error) {
	searchPattern := "%" + searchString + "%"
	const sql = `SELECT packageid, packagename, description FROM public."package" 
	        WHERE packagename ILIKE $1 OR description ILIKE $1 
	        LIMIT $2 OFFSET $3;`

	rows, err := r.client.Query(ctx, sql, searchPattern, limit, offset)
	if err != nil {
		_ = r.logger.Log("error", "Failed to search packages", "search", searchString, "err", err)
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanPackage, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		packages = append(packages, p)
	}

	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while searching packages", "err", err)
		return nil, err
	}

	return packages, nil
}

// GetAllPackages returns all packages with pagination.
func (r *GoodsPGRepository) GetAllPackages(ctx context.Context, limit int64, offset int64) ([]models.Package, error) {
	const sql = `SELECT packageid, packagename, description FROM public."package" LIMIT $1 OFFSET $2;`
	rows, err := r.client.Query(ctx, sql, limit, offset)
	if err != nil {
		_ = r.logger.Log("error", "Failed to retrieve all packages", "err", err)
		return nil, err
	}
	defer rows.Close()

	var packages []models.Package
	for rows.Next() {
		var p models.Package
		if err := rows.Scan(&p.ID, &p.PackageName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanPackage, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		packages = append(packages, p)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting all packages", "err", err)
		return nil, err
	}

	return packages, nil
}
