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
	msgFailedToScanTemplate = "Failed to scan template"
	fmtProductNotFound      = "Product with ID %d not found"
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

// ---------- TemplateRepository Implementation ----------

// GetTemplateByID retrieves template details along with its contents.
func (r *GoodsPGRepository) GetTemplateByID(ctx context.Context, id int64) (models.Template, error) {
	var template models.Template
	const sqlTemplate = `SELECT templateid, templatename, description FROM public."template" WHERE templateid = $1;`
	if err := r.client.QueryRow(ctx, sqlTemplate, id).Scan(&template.ID, &template.TemplateName, &template.Description); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return template, myerr.NotFound(fmt.Sprintf("Template with ID %d not found", id), nil)
		}
		_ = r.logger.Log("error", "Failed to get template by ID", "id", id, "err", err)
		return template, err
	}

	// Get template contents
	const sqlContents = `SELECT productid, quantity FROM public.templatecontent WHERE templateid = $1;`
	rows, err := r.client.Query(ctx, sqlContents, id)
	if err != nil {
		_ = r.logger.Log("error", "Failed to get template contents", "template_id", id, "err", err)
		return template, err
	}
	defer rows.Close()

	for rows.Next() {
		var c models.TemplateContent
		if err := rows.Scan(&c.ProductID, &c.Quantity); err != nil {
			_ = r.logger.Log("error", "Failed to scan template content", "err", err)
			return template, err
		}
		template.Content = append(template.Content, c)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting template contents", "err", err)
		return template, err
	}

	return template, nil
}

// GetProductsByTemplateID retrieves all products within a specific template.
func (r *GoodsPGRepository) GetProductsByTemplateID(ctx context.Context, templateID int64) ([]models.TemplateContent, error) {
	const sql = `SELECT productid, quantity FROM public.templatecontent WHERE templateid = $1;`
	rows, err := r.client.Query(ctx, sql, templateID)
	if err != nil {
		_ = r.logger.Log("error", "Failed to get products by template ID", "template_id", templateID, "err", err)
		return nil, err
	}
	defer rows.Close()

	var contents []models.TemplateContent
	for rows.Next() {
		var c models.TemplateContent
		if err := rows.Scan(&c.ProductID, &c.Quantity); err != nil {
			_ = r.logger.Log("error", "Failed to scan template content", "err", err)
			return nil, err
		}
		contents = append(contents, c)
	}

	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting products by template ID", "err", err)
		return nil, err
	}

	return contents, nil
}

// ListTemplates returns a list of all templates.
func (r *GoodsPGRepository) ListTemplates(ctx context.Context) ([]models.Template, error) {
	const sql = `SELECT templateid, templatename, description FROM public."template";`
	rows, err := r.client.Query(ctx, sql)
	if err != nil {
		_ = r.logger.Log("error", "Failed to list templates", "err", err)
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var p models.Template
		if err := rows.Scan(&p.ID, &p.TemplateName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanTemplate, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		templates = append(templates, p)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while listing templates", "err", err)
		return nil, err
	}

	return templates, nil
}

// CreateTemplate adds a new template to the database along with its contents.
func (r *GoodsPGRepository) CreateTemplate(ctx context.Context, template *models.Template) (err error) {
	const sqlInsertTemplate = `INSERT INTO public."template" (templatename, description) VALUES ($1, $2) RETURNING templateid;`

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

	// Insert template
	if err = tx.QueryRow(ctx, sqlInsertTemplate, template.TemplateName, template.Description).Scan(&template.ID); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return myerr.Conflict(fmt.Sprintf("Template with name %s already exists", template.TemplateName), err)
		}
		_ = r.logger.Log("error", "Failed to insert template", "template_name", template.TemplateName, "err", err)
		return err
	}

	// Insert template contents
	for _, content := range template.Content {
		if err = r.createProductToTemplate(ctx, tx, template.ID, content); err != nil {
			_ = r.logger.Log("error", "Failed to add product to template", "template_id", template.ID, "product_id", content.ProductID, "err", err)
			return err
		}
	}

	return nil
}

// createProductToTemplate adds a single template content entry.
func (r *GoodsPGRepository) createProductToTemplate(ctx context.Context, tx pgx.Tx, templateID int64, content models.TemplateContent) error {
	// Insert template content
	const sqlInsertContent = `INSERT INTO public.templatecontent (templateid, productid, quantity) VALUES ($1, $2, $3);`
	if _, err := tx.Exec(ctx, sqlInsertContent, templateID, content.ProductID, content.Quantity); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.ForeignKeyViolation {
			return myerr.NotFound(fmt.Sprintf(fmtProductNotFound, content.ProductID), err)
		} else if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return myerr.Conflict(fmt.Sprintf("Product with ID %d already exists in template", content.ProductID), err)
		}
		_ = r.logger.Log("error", "Failed to insert template content", "template_id", templateID, "product_id", content.ProductID, "err", err)
		return err
	}
	return nil
}

// DeleteTemplate deletes a template and its contents from the database by template ID.
func (r *GoodsPGRepository) DeleteTemplate(ctx context.Context, templateID int64) error {
	const (
		sqlDeleteContents = `DELETE FROM public.templatecontent WHERE templateid = $1;`
		sqlDeleteTemplate = `DELETE FROM public."template" WHERE templateid = $1;`
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

	// Delete template contents
	if _, err = tx.Exec(ctx, sqlDeleteContents, templateID); err != nil {
		_ = r.logger.Log("error", "Failed to delete template contents", "template_id", templateID, "err", err)
		return err
	}

	// Delete template
	ct, err := tx.Exec(ctx, sqlDeleteTemplate, templateID)
	if err != nil {
		_ = r.logger.Log("error", "Failed to delete template", "template_id", templateID, "err", err)
		return err
	}
	if ct.RowsAffected() == 0 {
		return myerr.NotFound(fmt.Sprintf("Template with id %d not found", templateID), nil)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		_ = r.logger.Log("error", "Failed to commit transaction", "err", err)
		return err
	}

	return nil
}

// SearchTemplates searches for templates by name or description with pagination.
func (r *GoodsPGRepository) SearchTemplates(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Template, error) {
	searchPattern := "%" + searchString + "%"
	const sql = `SELECT templateid, templatename, description FROM public."template"
	        WHERE templatename ILIKE $1 OR description ILIKE $1 
	        LIMIT $2 OFFSET $3;`

	rows, err := r.client.Query(ctx, sql, searchPattern, limit, offset)
	if err != nil {
		_ = r.logger.Log("error", "Failed to search templates", "search", searchString, "err", err)
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var p models.Template
		if err := rows.Scan(&p.ID, &p.TemplateName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanTemplate, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		templates = append(templates, p)
	}

	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while searching templates", "err", err)
		return nil, err
	}

	return templates, nil
}

// GetAllTemplates returns all templates with pagination.
func (r *GoodsPGRepository) GetAllTemplates(ctx context.Context, limit int64, offset int64) ([]models.Template, error) {
	const sql = `SELECT templateid, templatename, description FROM public."template" LIMIT $1 OFFSET $2;`
	rows, err := r.client.Query(ctx, sql, limit, offset)
	if err != nil {
		_ = r.logger.Log("error", "Failed to retrieve all templates", "err", err)
		return nil, err
	}
	defer rows.Close()

	var templates []models.Template
	for rows.Next() {
		var p models.Template
		if err := rows.Scan(&p.ID, &p.TemplateName, &p.Description); err != nil {
			_ = r.logger.Log("error", msgFailedToScanTemplate, "err", err)
			continue // Можно решить, нужно ли пропускать или возвращать ошибку
		}
		templates = append(templates, p)
	}
	if err := rows.Err(); err != nil {
		_ = r.logger.Log("error", "Rows iteration error while getting all templates", "err", err)
		return nil, err
	}

	return templates, nil
}
