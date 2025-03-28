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

func TestDeleteTemplate(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	repo := postgresql.NewGoodsRepository(pool, nil)

	t.Cleanup(func() {
		pool.Close()
	})

	t.Log("Очистка данных")
	_, err = pool.Exec(ctx, `DELETE FROM packagecontent`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM product`)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `
		INSERT INTO product (id, name, description, price, imageurl, sku)
		VALUES (1001, 'Test Product', 'Desc', 10.0, 'url', 'SKU1001')
	`)
	require.NoError(t, err)

	t.Run("Успешное удаление шаблона", func(t *testing.T) {
		t.Log("Создание шаблона перед удалением")

		template := &models.Template{
			TemplateName: "To Delete",
			Description:  "Will be deleted",
			Content: []models.TemplateContent{
				{ProductID: 1001, Quantity: 2},
			},
		}

		err := repo.CreateTemplate(ctx, template)
		require.NoError(t, err)
		t.Logf("Создан шаблон с ID: %d", template.ID)

		t.Log("Удаление шаблона")
		err = repo.DeleteTemplate(ctx, template.ID)
		assert.NoError(t, err)

		t.Log("Проверка удаления из package")
		var count int
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM package WHERE packageid = $1`, template.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)

		t.Log("Проверка удаления из packagecontent")
		err = pool.QueryRow(ctx, `SELECT COUNT(*) FROM packagecontent WHERE packageid = $1`, template.ID).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 0, count)
	})

	t.Run("Шаблон не найден", func(t *testing.T) {
		t.Log("Удаление несуществующего шаблона")
		err := repo.DeleteTemplate(ctx, 999999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Ошибка соединения с БД", func(t *testing.T) {
		t.Log("Имитируем ошибку соединения")
		brokenPool, err := pgxpool.New(ctx, dbURL)
		require.NoError(t, err)
		brokenRepo := postgresql.NewGoodsRepository(brokenPool, nil)
		brokenPool.Close()

		err = brokenRepo.DeleteTemplate(ctx, 1)
		assert.Error(t, err)
	})
}
