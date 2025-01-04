package tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestCreateProduct_Success() {
	newProduct := &models.Product{
		Name:        "New Product",
		Description: "New Description",
		Price:       79.99,
		ImageURL:    "http://example.com/newimage.png",
		SKU:         "SKUNEW",
	}
	expectedID := int64(3)

	suite.mockRepo.On("CreateProduct", mock.Anything, newProduct).
		Run(func(args mock.Arguments) {
			product := args.Get(1).(*models.Product)
			product.ID = expectedID
			// Дополнительные проверки
			assert.NotEmpty(suite.T(), product.Name, "Product name should not be empty")
			assert.GreaterOrEqual(suite.T(), product.Price, float64(0), "Product price should be non-negative")
		}).
		Return(expectedID, nil).
		Once()

	createdID, err := suite.svc.CreateProduct(context.Background(), newProduct)

	assert.NoError(suite.T(), err, "Expected no error when creating product")
	assert.Equal(suite.T(), expectedID, createdID, "Expected created ID to match the mocked ID")
	assert.Equal(suite.T(), expectedID, newProduct.ID, "Expected product ID to be set to the mocked ID")
}

func (suite *ServiceTestSuite) TestCreateProduct_RepositoryError() {
	newProduct := &models.Product{
		Name:        "New Product",
		Description: "New Description",
		Price:       79.99,
		ImageURL:    "http://example.com/newimage.png",
		SKU:         "SKUNEW",
	}
	expectedError := errors.New("database error")

	suite.mockRepo.On("CreateProduct", mock.Anything, newProduct).
		Return(int64(0), expectedError).
		Once()

	createdID, err := suite.svc.CreateProduct(context.Background(), newProduct)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
	assert.Equal(suite.T(), int64(0), newProduct.ID, "Expected product ID to remain 0 on error")
}
