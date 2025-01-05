package unit_tests

import (
	"context"
	"errors"

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

func (suite *ServiceTestSuite) TestDeleteProduct_RepositoryError() {
	productID := int64(2)
	expectedError := errors.New("database error")

	suite.mockRepo.On("DeleteProduct", mock.Anything, mock.MatchedBy(func(id int64) bool {
		return id == productID
	})).
		Return(expectedError).
		Once()

	err := suite.svc.DeleteProduct(context.Background(), productID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
}
