package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetProductByID_Success() {
	productID := int64(1)
	expectedProduct := createTestProduct(productID, "Test Product")

	suite.mockRepo.
		EXPECT().
		GetProductByID(mock.Anything, productID).
		Return(expectedProduct, nil).
		Once()

	product, err := suite.svc.GetProductByID(context.Background(), productID)

	assert.NoError(suite.T(), err, "Expected no error when getting product by ID")
	assert.Equal(suite.T(), expectedProduct, product, "Expected product to match the mocked product")
}

func (suite *ServiceTestSuite) TestGetProductByID_NotFound() {
	productID := int64(2)
	expectedError := errors.New("product not found")

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(models.Product{}, expectedError).
		Once()

	product, err := suite.svc.GetProductByID(context.Background(), productID)

	assert.Error(suite.T(), err, "Expected error when product is not found")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), models.Product{}, product, "Expected product to be empty")
}

func (suite *ServiceTestSuite) TestGetProductByID_RepositoryError() {
	productID := int64(3)
	expectedError := errors.New("database error")

	suite.mockRepo.On("GetProductByID", mock.Anything, productID).
		Return(models.Product{}, expectedError).
		Once()

	product, err := suite.svc.GetProductByID(context.Background(), productID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), models.Product{}, product, "Expected product to be empty")
}
