//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateTemplate(t *testing.T) {
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

	// Добавим продукты, которые будут использоваться в шаблоне
	_, err = pool.Exec(ctx, `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES 
			(1001, 'Product A', 'Desc A', 10.0, 'url1', 'SKU_A'),
			(1002, 'Product B', 'Desc B', 20.0, 'url2', 'SKU_B')
	`)
	require.NoError(t, err)

	t.Run("Успешное создание шаблона", func(t *testing.T) {
		template := &models.Template{
			TemplateName: "Template A",
			Description:  "Test template A",
			Content: []models.TemplateContent{
				{ProductID: 1001, Quantity: 2},
				{ProductID: 1002, Quantity: 3},
			},
		}

		err := repo.CreateTemplate(ctx, template)
		assert.NoError(t, err)
		assert.True(t, template.ID > 0)

		// Проверим, что шаблон и его содержимое добавлены
		var count int
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM package WHERE packagename = $1`, template.TemplateName).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM packagecontent WHERE packageid = $1`, template.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 2, count)
	})

	t.Run("Ошибка: продукт не существует", func(t *testing.T) {
		template := &models.Template{
			TemplateName: "Template B",
			Description:  "Bad content",
			Content: []models.TemplateContent{
				{ProductID: 9999, Quantity: 1}, // такого продукта нет
			},
		}
		err := repo.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "violates foreign key constraint")
	})

	t.Run("Ошибка соединения с БД", func(t *testing.T) {
		pool.Close()

		template := &models.Template{
			TemplateName: "Disconnected",
			Description:  "Should fail",
			Content: []models.TemplateContent{
				{ProductID: 1001, Quantity: 1},
			},
		}
		err := repo.CreateTemplate(ctx, template)
		assert.Error(t, err)
	})
}
