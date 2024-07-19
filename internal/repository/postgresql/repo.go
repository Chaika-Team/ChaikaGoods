package postgresql

import (
	"ChaikaGoods/internal/models"
	"context"
	"fmt"
	"github.com/go-kit/kit/log"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	client *pgxpool.Pool
	log    log.Logger
}

// NewFolderRepository создает новый экземпляр Repository.
func NewFolderRepository(client *pgxpool.Pool, logger log.Logger) *Repository {
	return &Repository{
		client: client,
		log:    logger,
	}
}

// GetProductByID возвращает продукт по его ID.
func (r *Repository) GetProductByID(ctx context.Context, id int64) (*models.Product, error) {
	sql := `SELECT id, name, description, price, imageurl, sku FROM public.product WHERE id = $1;`
	row := r.client.QueryRow(ctx, sql, id)

	var p models.Product
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.SKU)
	if err != nil {
		r.log.Log("error", fmt.Sprintf("Failed to get product by ID: %v", err))
		return nil, err
	}

	return &p, nil
}

// CreateProduct добавляет новый продукт в базу данных.
func (r *Repository) CreateProduct(ctx context.Context, product *models.Product) error {
	sql := `INSERT INTO public.product (name, description, price, imageurl, sku) VALUES ($1, $2, $3, $4, $5) RETURNING id;`
	err := r.client.QueryRow(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).Scan(&product.ID)
	if err != nil {
		r.log.Log("error", fmt.Sprintf("Failed to create product: %v", err))
		return err
	}
	return nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (r *Repository) UpdateProduct(ctx context.Context, product *models.Product) error {
	sql := `UPDATE public.product SET name = $1, description = $2, price = $3, imageurl = $4, sku = $5 WHERE id = $6;`
	commandTag, err := r.client.Exec(ctx, sql, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID)
	if err != nil {
		r.log.Log("error", fmt.Sprintf("Failed to update product: %v", err))
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных по его ID.
func (r *Repository) DeleteProduct(ctx context.Context, id int64) error {
	sql := `DELETE FROM public.product WHERE id = $1;`
	commandTag, err := r.client.Exec(ctx, sql, id)
	if err != nil {
		r.log.Log("error", fmt.Sprintf("Failed to delete product: %v", err))
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return fmt.Errorf("no rows affected")
	}
	return nil
}

// Дополнительные методы для работы с пакетами и содержимым пакетов могут быть реализованы аналогично
