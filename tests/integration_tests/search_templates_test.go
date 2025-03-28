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

func TestSearchTemplates(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	repo := postgresql.NewGoodsRepository(pool, nil)
	t.Cleanup(func() { pool.Close() })

	// Очистка
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)

	// Вставка данных
	_, err = pool.Exec(ctx, `
		INSERT INTO package (packageid, packagename, description) VALUES
		(501, 'Summer Promo', 'Hot discount templates'),
		(502, 'Winter Sale', 'Cold deals inside'),
		(503, 'Test Template', 'Something to test'),
		(504, 'Ignore This', 'Should not match')
	`)
	require.NoError(t, err)

	t.Run("Поиск по ключевому слову (название)", func(t *testing.T) {
		results, err := repo.SearchTemplates(ctx, "Promo", 10, 0)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Summer Promo", results[0].TemplateName)
	})

	t.Run("Поиск по описанию", func(t *testing.T) {
		results, err := repo.SearchTemplates(ctx, "deals", 10, 0)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Winter Sale", results[0].TemplateName)
	})

	t.Run("Поиск с оффсетом и лимитом", func(t *testing.T) {
		results, err := repo.SearchTemplates(ctx, "template", 1, 0)
		assert.NoError(t, err)
		assert.Len(t, results, 1)

		names := []string{"Test Template", "Ignore This", "Summer Promo", "Winter Sale"}
		assert.Contains(t, names, results[0].TemplateName)
	})

	t.Run("Пустой результат", func(t *testing.T) {
		results, err := repo.SearchTemplates(ctx, "nothing", 10, 0)
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("Ошибка соединения", func(t *testing.T) {
		brokenPool, err := pgxpool.New(ctx, dbURL)
		require.NoError(t, err)
		brokenRepo := postgresql.NewGoodsRepository(brokenPool, nil)
		brokenPool.Close()

		_, err = brokenRepo.SearchTemplates(ctx, "Promo", 10, 0)
		assert.Error(t, err)
	})
}
