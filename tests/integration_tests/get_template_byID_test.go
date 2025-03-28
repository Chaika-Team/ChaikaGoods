//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetTemplateByID(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблиц
	_, err = pool.Exec(ctx, `DELETE FROM packagecontent`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM product`)
	require.NoError(t, err)

	// Вставка необходимых продуктов
	_, err = pool.Exec(ctx, `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES 
			(101, 'Product 101', 'Desc 101', 10.0, 'url1', 'SKU101'),
			(102, 'Product 102', 'Desc 102', 20.0, 'url2', 'SKU102')
		ON CONFLICT (id) DO NOTHING;
	`)
	require.NoError(t, err)

	// Вставка шаблона и содержимого
	templateID := int64(100)
	_, err = pool.Exec(ctx, `
		INSERT INTO package (packageid, packagename, description)
		VALUES ($1, $2, $3)
	`, templateID, "Sample Template", "Test template")
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO packagecontent (packageid, productid, quantity)
		VALUES 
			($1, 101, 2),
			($1, 102, 5)
	`, templateID)
	require.NoError(t, err)

	t.Run("Успешное получение шаблона", func(t *testing.T) {
		template, err := repo.GetTemplateByID(ctx, templateID)
		assert.NoError(t, err)
		assert.Equal(t, templateID, template.ID)
		assert.Equal(t, "Sample Template", template.TemplateName)
		assert.Equal(t, "Test template", template.Description)
		assert.Len(t, template.Content, 2)
		assert.Equal(t, int64(101), template.Content[0].ProductID)
		assert.Equal(t, 2, template.Content[0].Quantity)
		assert.Equal(t, int64(102), template.Content[1].ProductID)
		assert.Equal(t, 5, template.Content[1].Quantity)
	})

	t.Run("Шаблон не найден", func(t *testing.T) {
		_, err := repo.GetTemplateByID(ctx, 999999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Ошибка соединения с БД", func(t *testing.T) {
		pool.Close()
		_, err := repo.GetTemplateByID(ctx, templateID)
		assert.Error(t, err)
	})
}
