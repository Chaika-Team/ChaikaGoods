//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
)

func TestCreateTemplateProductNotFound(t *testing.T) {
	keyspace := "create_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template := models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	id, err := svc.AddTemplate(ctx, &template)
	if err == nil {
		t.Fatalf("Expected err, got id: %v", id)
	}
}
