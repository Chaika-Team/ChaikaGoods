//go:build integration
// +build integration

package integration

import (
	"context"
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestSearchTemplateNonExisting(t *testing.T) {
	keyspace := "search_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template, err := svc.SearchTemplates(ctx, "Default", 1, 1)
	if err != nil {
		t.Fatalf("Expected empty template, got err: %v", err)
	}
	assert.Nil(t, template)
}

func TestSearchTemplateExistingVariousString(t *testing.T) {
	keyspace := "search_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template := models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content:      []models.TemplateContent{},
	}

	id, err := svc.AddTemplate(ctx, &template)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	assert.Equal(t, id, int64(1))

	// Empty string case
	templateResponse, err := svc.SearchTemplates(ctx, "", 1, 0)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	if templateResponse == nil {
		t.Fatalf("Expected non empty template, got nil: %v", templateResponse)
	}
	assert.Equal(t, templateResponse[0].TemplateName, template.TemplateName)
	assert.Equal(t, templateResponse[0].Description, template.Description)
	assert.Nil(t, templateResponse[0].Content)

	// Existing string case
	templateResponse, err = svc.SearchTemplates(ctx, "Template", 1, 0)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	if templateResponse == nil {
		t.Fatalf("Expected non empty template, got nil: %v", templateResponse)
	}
	assert.Equal(t, templateResponse[0].TemplateName, template.TemplateName)
	assert.Equal(t, templateResponse[0].Description, template.Description)
	assert.Nil(t, templateResponse[0].Content)

	// Non existing string case
	templateResponse, err = svc.SearchTemplates(ctx, "NonExists", 1, 0)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	if templateResponse != nil {
		t.Fatalf("Expected nil template, got: %v", templateResponse)
	}
	assert.Nil(t, templateResponse)
}

func TestSearchTemplateNotFull(t *testing.T) {
	keyspace := "search_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template := models.Template{
		TemplateName: "New Template",
		Content:      []models.TemplateContent{},
	}

	id, err := svc.AddTemplate(ctx, &template)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	assert.Equal(t, id, int64(1))

	templateResponse, err := svc.SearchTemplates(ctx, "New Template", 1, 0)
	if err != nil {
		t.Fatalf("Expected successful operation, got err: %v", err)
	}
	if templateResponse == nil {
		t.Fatalf("Expected non empty template, got nil: %v", templateResponse)
	}
	assert.Equal(t, templateResponse[0].TemplateName, template.TemplateName)
	assert.Equal(t, templateResponse[0].Description, "")
	assert.Nil(t, templateResponse[0].Content)
}

func TestSearchTemplateInvalidRequests(t *testing.T) {
	keyspace := "search_template_test"

	svc := setupService(t, keyspace)
	ctx := context.Background()

	template, err := svc.SearchTemplates(ctx, "Context", -10, 10)
	if err == nil {
		t.Fatalf("Expected error, got template: %v", template)
	}
	assert.Nil(t, template)
	assert.Error(t, err)

	template, err = svc.SearchTemplates(ctx, "Context", 10, -10)
	if err == nil {
		t.Fatalf("Expected error, got template: %v", template)
	}
	assert.Nil(t, template)
	assert.Error(t, err)
}
