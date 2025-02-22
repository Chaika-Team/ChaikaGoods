package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestDeleteProduct_Success() {
	productID := int64(1)

	suite.mockRepo.On("DeleteProduct", mock.Anything, mock.MatchedBy(func(id int64) bool {
		return id == productID
	})).
		Return(nil).
		Once()

	err := suite.svc.DeleteProduct(context.Background(), productID)

	assert.NoError(suite.T(), err, "Expected no error when deleting product")
}

func (suite *ServiceTestSuite) TestDeleteProduct_NotFound() {
	productID := int64(2)
	expectedError := myerr.NotFound("Product not found", nil)

	suite.mockRepo.On("DeleteProduct", mock.Anything, mock.MatchedBy(func(id int64) bool {
		return id == productID
	})).
		Return(expectedError).
		Once()

	err := suite.svc.DeleteProduct(context.Background(), productID)

	assert.Error(suite.T(), err, "Expected error when product is not found")
	assert.True(suite.T(), myerr.IsNotFound(err), "Expected error to be of type NotFound")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
}

func (suite *ServiceTestSuite) TestDeleteProduct_RepositoryError() {
	productID := int64(2)
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("DeleteProduct", mock.Anything, mock.MatchedBy(func(id int64) bool {
		return id == productID
	})).
		Return(expectedError).
		Once()

	err := suite.svc.DeleteProduct(context.Background(), productID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.True(suite.T(), myerr.IsInternal(err), "Expected error to be of type Internal")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
}
