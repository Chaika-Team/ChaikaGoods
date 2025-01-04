package tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetAllProducts_Success() {
	expectedProducts := []models.Product{
		createTestProduct(1, "Product 1"),
		createTestProduct(2, "Product 2"),
	}

	suite.mockRepo.On("GetAllProducts", mock.Anything).
		Return(expectedProducts, nil).
		Once()

	products, err := suite.svc.GetAllProducts(context.Background())

	assert.NoError(suite.T(), err, "Expected no error when getting all products")
	assert.Equal(suite.T(), expectedProducts, products, "Expected products to match the mocked products")
}

func (suite *ServiceTestSuite) TestGetAllProducts_RepositoryError() {
	expectedError := errors.New("database error")

	suite.mockRepo.On("GetAllProducts", mock.Anything).
		Return(nil, expectedError).
		Once()

	products, err := suite.svc.GetAllProducts(context.Background())

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), products, "Expected products to be nil")
}
