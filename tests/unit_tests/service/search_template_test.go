package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestSearchTemplates_WithQuery_Success() {
	searchString := "Test"
	limit := int64(10)
	offset := int64(0)
	expectedTemplates := []models.Template{
		createTestTemplate(1, "Test Template 1"),
	}

	suite.mockRepo.On("SearchTemplates", mock.Anything, searchString, limit, offset).
		Return(expectedTemplates, nil).
		Once()

	templates, err := suite.svc.SearchTemplates(context.Background(), searchString, limit, offset)

	assert.NoError(suite.T(), err, "Expected no error when searching Templates with query")
	assert.Equal(suite.T(), expectedTemplates, templates, "Expected templates to match the mocked templates")
}

func (suite *ServiceTestSuite) TestSearchTemplates_EmptyQuery_Success() {
	searchString := ""
	limit := int64(5)
	offset := int64(0)
	expectedTemplates := []models.Template{
		createTestTemplate(1, "Template 1"),
		createTestTemplate(2, "Template 2"),
	}

	suite.mockRepo.On("GetAllTemplates", mock.Anything, limit, offset).
		Return(expectedTemplates, nil).
		Once()

	templates, err := suite.svc.SearchTemplates(context.Background(), searchString, limit, offset)

	assert.NoError(suite.T(), err, "Expected no error when searching Templates with empty query")
	assert.Equal(suite.T(), expectedTemplates, templates, "Expected templates to match the mocked templates")
}

func (suite *ServiceTestSuite) TestSearchTemplates_WithQuery_RepositoryError() {
	searchString := "NonExistent"
	limit := int64(10)
	offset := int64(0)
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("SearchTemplates", mock.Anything, searchString, limit, offset).
		Return(nil, expectedError).
		Once()

	templates, err := suite.svc.SearchTemplates(context.Background(), searchString, limit, offset)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during search")
	assert.True(suite.T(), myerr.IsInternal(err), "Expected error to be of type Internal")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), templates, "Expected templates to be nil on error")
}

func (suite *ServiceTestSuite) TestSearchTemplates_EmptyQuery_RepositoryError() {
	searchString := ""
	limit := int64(5)
	offset := int64(0)
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("GetAllTemplates", mock.Anything, limit, offset).
		Return(nil, expectedError).
		Once()

	templates, err := suite.svc.SearchTemplates(context.Background(), searchString, limit, offset)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during GetAllTemplates")
	assert.True(suite.T(), myerr.IsInternal(err), "Expected error to be of type Internal")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), templates, "Expected templates to be nil on error")
}
