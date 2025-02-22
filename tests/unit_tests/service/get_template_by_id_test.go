package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetTemplateByID_Success() {
	TemplateID := int64(1)
	expectedTemplate := createTestTemplate(TemplateID, "Test Template")

	suite.mockRepo.On("GetTemplateByID", mock.Anything, TemplateID).
		Return(expectedTemplate, nil).
		Once()

	Template, err := suite.svc.GetTemplateByID(context.Background(), TemplateID)

	assert.NoError(suite.T(), err, "Expected no error when getting Template by ID")
	assert.Equal(suite.T(), expectedTemplate, Template, "Expected Template to match the mocked Template")
}

func (suite *ServiceTestSuite) TestGetTemplateByID_NotFound() {
	TemplateID := int64(2)
	expectedError := myerr.NotFound("Template not found", nil)

	suite.mockRepo.On("GetTemplateByID", mock.Anything, TemplateID).
		Return(models.Template{}, expectedError).
		Once()

	Template, err := suite.svc.GetTemplateByID(context.Background(), TemplateID)

	assert.Error(suite.T(), err, "Expected error when Template is not found")
	assert.True(suite.T(), myerr.IsNotFound(err), "Expected error to be of type NotFound")
	assert.Equal(suite.T(), models.Template{}, Template, "Expected Template to be empty")
}

func (suite *ServiceTestSuite) TestGetTemplateByID_RepositoryError() {
	TemplateID := int64(3)
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("GetTemplateByID", mock.Anything, TemplateID).
		Return(models.Template{}, expectedError).
		Once()

	Template, err := suite.svc.GetTemplateByID(context.Background(), TemplateID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during GetTemplateByID")
	assert.True(suite.T(), myerr.IsInternal(err), "Expected error to be of type Internal")
	assert.Equal(suite.T(), models.Template{}, Template, "Expected Template to be empty")
}
