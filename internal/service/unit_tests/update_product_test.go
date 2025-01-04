// internal/service/update_product_test.go
package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestUpdateProduct_Success() {
	updatedProduct := &models.Product{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       89.99,
		ImageURL:    "http://example.com/updatedimage.png",
		SKU:         "SKU001U",
	}

	suite.mockRepo.On("UpdateProduct", mock.Anything, updatedProduct).
		Return(nil).
		Once()

	err := suite.svc.UpdateProduct(context.Background(), updatedProduct)

	assert.NoError(suite.T(), err, "Expected no error when updating product")
}

func (suite *ServiceTestSuite) TestUpdateProduct_RepositoryError() {
	updatedProduct := &models.Product{
		ID:          2,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       89.99,
		ImageURL:    "http://example.com/updatedimage.png",
		SKU:         "SKU002U",
	}
	expectedError := errors.New("database error")

	suite.mockRepo.On("UpdateProduct", mock.Anything, updatedProduct).
		Return(expectedError).
		Once()

	err := suite.svc.UpdateProduct(context.Background(), updatedProduct)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
}
