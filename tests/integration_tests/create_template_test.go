//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, id, int64(0))
}

func TestCreateTemplateEmpty(t *testing.T) {
	keyspace := "create_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template := models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content:      []models.TemplateContent{},
	}
	id, err := svc.AddTemplate(ctx, &template)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}
	assert.Equal(t, id, int64(1))
}

func TestCreateTemplateWithProducts(t *testing.T) {
	keyspace := "create_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	product := models.Product{
		Name:        "Milk",
		Price:       99.99,
		Description: "Basic 2.5% milk",
	}

	id, err := svc.CreateProduct(ctx, &product)
	if err != nil {
		t.Fatalf("Cannot create product, got err: %v", err)
	}

	template := models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content:      []models.TemplateContent{},
	}

	id, err = svc.AddTemplate(ctx, &template)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}
	assert.Equal(t, id, int64(1))
}

func TestCreateTemplateExisting(t *testing.T) {
	keyspace := "create_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template := models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content:      []models.TemplateContent{},
	}
	id, err := svc.AddTemplate(ctx, &template)
	if err != nil {
		t.Fatalf("Expected id, got err: %v", err)
	}

	templateCopy := template
	id, err = svc.AddTemplate(ctx, &templateCopy)
	if err == nil {
		t.Fatalf("Expected error, got id: %v", id)
	}
}
