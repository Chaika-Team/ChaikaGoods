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

func TestListTemplates(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	defer pool.Close()

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблицы перед тестом
	_, err = pool.Exec(ctx, `DELETE FROM packagecontent`)
	require.NoError(t, err)
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)

	t.Run("Пустой список шаблонов", func(t *testing.T) {
		templates, err := repo.ListTemplates(ctx)
		assert.NoError(t, err)
		assert.Len(t, templates, 0)
	})

	t.Run("Успешное получение шаблонов", func(t *testing.T) {
		// Вставка двух шаблонов
		_, err := pool.Exec(ctx, `
			INSERT INTO package (packageid, packagename, description)
			VALUES 
				(1, 'Template 1', 'Description 1'),
				(2, 'Template 2', 'Description 2')
		`)
		require.NoError(t, err)

		templates, err := repo.ListTemplates(ctx)
		assert.NoError(t, err)
		assert.Len(t, templates, 2)

		assert.Equal(t, int64(1), templates[0].ID)
		assert.Equal(t, "Template 1", templates[0].TemplateName)
		assert.Equal(t, "Description 1", templates[0].Description)

		assert.Equal(t, int64(2), templates[1].ID)
		assert.Equal(t, "Template 2", templates[1].TemplateName)
		assert.Equal(t, "Description 2", templates[1].Description)
	})

	t.Run("Ошибка соединения с БД", func(t *testing.T) {
		pool.Close()
		_, err := repo.ListTemplates(ctx)
		assert.Error(t, err)
	})
}
