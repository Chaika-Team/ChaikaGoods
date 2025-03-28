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

func TestGetAllTemplates(t *testing.T) {
	ctx := context.Background()

	dbURL := "postgres://" + os.Getenv("DB_USER") + ":" + os.Getenv("DB_PASS") +
		"@" + os.Getenv("DB_HOST") + ":" + os.Getenv("DB_PORT") + "/" + os.Getenv("DB_NAME")

	pool, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)
	t.Cleanup(func() { pool.Close() })

	repo := postgresql.NewGoodsRepository(pool, nil)

	// Очистка таблицы
	_, err = pool.Exec(ctx, `DELETE FROM package`)
	require.NoError(t, err)

	// Вставка 3 шаблонов
	_, err = pool.Exec(ctx, `
		INSERT INTO package (packageid, packagename, description) VALUES
		(601, 'Template A', 'Desc A'),
		(602, 'Template B', 'Desc B'),
		(603, 'Template C', 'Desc C')
	`)
	require.NoError(t, err)

	t.Run("Получение всех шаблонов (limit 2 offset 0)", func(t *testing.T) {
		templates, err := repo.GetAllTemplates(ctx, 2, 0)
		assert.NoError(t, err)
		assert.Len(t, templates, 2)
	})

	t.Run("Получение с offset (limit 2 offset 2)", func(t *testing.T) {
		templates, err := repo.GetAllTemplates(ctx, 2, 2)
		assert.NoError(t, err)
		assert.Len(t, templates, 1)
		assert.Equal(t, "Template C", templates[0].TemplateName)
	})

	t.Run("Получение с offset за пределами (limit 2 offset 10)", func(t *testing.T) {
		templates, err := repo.GetAllTemplates(ctx, 2, 10)
		assert.NoError(t, err)
		assert.Len(t, templates, 0)
	})

	t.Run("Ошибка соединения", func(t *testing.T) {
		brokenPool, err := pgxpool.New(ctx, dbURL)
		require.NoError(t, err)
		brokenRepo := postgresql.NewGoodsRepository(brokenPool, nil)
		brokenPool.Close()

		_, err = brokenRepo.GetAllTemplates(ctx, 1, 0)
		assert.Error(t, err)
	})
}
