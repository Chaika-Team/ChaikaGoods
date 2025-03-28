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

func TestCreateTemplateProductNotFound(t *testing.T) {
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

	t.Run("Ошибка: продукты не существуют", func(t *testing.T) {
		template := &models.Template{
			TemplateName: "Bad Template",
			Description:  "Should fail",
			Content: []models.TemplateContent{
				{ProductID: 9999, Quantity: 1},
				{ProductID: 8888, Quantity: 2},
			},
		}

		err := repo.CreateTemplate(ctx, template)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "violates foreign key constraint")
	})
}

// package integration

// import (
// 	"context"
// 	"testing"

// 	"github.com/Chaika-Team/ChaikaGoods/internal/models"
// )

// func TestCreateTemplateProductNotFound(t *testing.T) {
// 	keyspace := "create_template_test"

// 	svc := setupService(t, keyspace)
// 	ctx := context.Background()

// 	template := models.Template{
// 		TemplateName: "New Template",
// 		Description:  "New Template Description",
// 		Content: []models.TemplateContent{
// 			{ProductID: 1, Quantity: 2},
// 			{ProductID: 2, Quantity: 3},
// 		},
// 	}
// 	id, err := svc.AddTemplate(ctx, &template)
// 	if err == nil {
// 		t.Fatalf("Expected err, got id: %v", id)
// 	}
// }
