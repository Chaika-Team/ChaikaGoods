package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestAddTemplate_Success() {
	newTemplate := &models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedID := int64(1)

	suite.mockRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(p *models.Template) bool {
		return p.TemplateName == newTemplate.TemplateName &&
			p.Description == newTemplate.Description &&
			len(p.Content) == len(newTemplate.Content)
	})).
		Run(func(args mock.Arguments) {
			Template := args.Get(1).(*models.Template)
			Template.ID = expectedID
			// Дополнительные проверки
			assert.NotEmpty(suite.T(), Template.TemplateName, "Template name should not be empty")
			assert.NotEmpty(suite.T(), Template.Description, "Template description should not be empty")
			assert.Len(suite.T(), Template.Content, 2, "Template should have two contents")
		}).
		Return(nil).
		Once()

	createdID, err := suite.svc.AddTemplate(context.Background(), newTemplate)

	assert.NoError(suite.T(), err, "Expected no error when adding Template")
	assert.Equal(suite.T(), expectedID, createdID, "Expected created ID to match the mocked ID")
	assert.Equal(suite.T(), expectedID, newTemplate.ID, "Expected Template ID to be set to the mocked ID")
}

func (suite *ServiceTestSuite) TestAddTemplate_AlreadyExists() {
	newTemplate := &models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.Conflict("Template with name New Template already exists", nil)

	suite.mockRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(p *models.Template) bool {
		return p.TemplateName == newTemplate.TemplateName &&
			p.Description == newTemplate.Description &&
			len(p.Content) == len(newTemplate.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddTemplate(context.Background(), newTemplate)

	assert.Error(suite.T(), err, "Expected error when Template already exists")
	assert.True(suite.T(), myerr.IsConflict(err), "Expected error to be of type Conflict")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}

func (suite *ServiceTestSuite) TestAddTemplate_ProductNotFound() {
	newTemplate := &models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.NotFound("Product with ID 1 not found", nil)

	suite.mockRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(p *models.Template) bool {
		return p.TemplateName == newTemplate.TemplateName &&
			p.Description == newTemplate.Description &&
			len(p.Content) == len(newTemplate.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddTemplate(context.Background(), newTemplate)

	assert.Error(suite.T(), err, "Expected error when Product is not found")
	assert.True(suite.T(), myerr.IsNotFound(err), "Expected error to be of type NotFound")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}

func (suite *ServiceTestSuite) TestAddTemplate_RepositoryError() {
	newTemplate := &models.Template{
		TemplateName: "New Template",
		Description:  "New Template Description",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("CreateTemplate", mock.Anything, mock.MatchedBy(func(p *models.Template) bool {
		return p.TemplateName == newTemplate.TemplateName &&
			p.Description == newTemplate.Description &&
			len(p.Content) == len(newTemplate.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddTemplate(context.Background(), newTemplate)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during AddTemplate")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}
