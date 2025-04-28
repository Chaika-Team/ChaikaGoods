package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetAllProducts_Success() {
	expectedProducts := []models.Product{
		createTestProduct(1, "Product 1"),
		createTestProduct(2, "Product 2"),
	}

	limit := int64(10)
	offset := int64(0)

	// Настраиваем мок так, чтобы он ждал вызов с контекстом, limit и offset
	suite.mockRepo.
		On("GetAllProducts", mock.Anything, limit, offset).
		Return(expectedProducts, nil).
		Once()

	// Вызываем именно новую сигнатуру
	products, err := suite.svc.GetAllProducts(context.Background(), limit, offset)

	assert.NoError(suite.T(), err, "Expected no error when getting all products")
	assert.Equal(suite.T(), expectedProducts, products, "Expected products to match the mocked products")

	// Проверяем, что мок-предсказания сработали
	suite.mockRepo.AssertExpectations(suite.T())
}

func (suite *ServiceTestSuite) TestGetAllProducts_RepositoryError() {
	// эмулируем внутреннюю ошибку
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	limit := int64(5)
	offset := int64(2)

	suite.mockRepo.
		On("GetAllProducts", mock.Anything, limit, offset).
		Return(nil, expectedError).
		Once()

	products, err := suite.svc.GetAllProducts(context.Background(), limit, offset)

	assert.Error(suite.T(), err, "Expected error when repository returns an error")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), products, "Expected products to be nil")

	suite.mockRepo.AssertExpectations(suite.T())
}
