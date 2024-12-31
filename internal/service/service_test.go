package service

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/repository/mocks"
	"context"
	"errors"
	"testing"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductByID(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	productID := int64(1)
	expectedProduct := models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       99.99,
		ImageURL:    "http://example.com/image.png",
		SKU:         "SKU12345",
	}

	// Настраиваем ожидания мок репозитория
	mockRepo.On("GetProductByID", mock.Anything, productID).Return(expectedProduct, nil)

	// Вызываем метод сервиса
	product, err := svc.GetProductByID(context.Background(), productID)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestGetProductByID_NotFound(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	productID := int64(2)
	expectedError := errors.New("product not found")

	// Настраиваем ожидания мок репозитория
	mockRepo.On("GetProductByID", mock.Anything, productID).Return(models.Product{}, expectedError)

	// Вызываем метод сервиса
	product, err := svc.GetProductByID(context.Background(), productID)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Equal(t, models.Product{}, product)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestGetAllProducts(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	expectedProducts := []models.Product{
		{
			ID:          1,
			Name:        "Product 1",
			Description: "Description 1",
			Price:       49.99,
			ImageURL:    "http://example.com/image1.png",
			SKU:         "SKU001",
		},
		{
			ID:          2,
			Name:        "Product 2",
			Description: "Description 2",
			Price:       59.99,
			ImageURL:    "http://example.com/image2.png",
			SKU:         "SKU002",
		},
	}

	// Настраиваем ожидания мок репозитория
	mockRepo.On("GetAllProducts", mock.Anything).Return(expectedProducts, nil)

	// Вызываем метод сервиса
	products, err := svc.GetAllProducts(context.Background())

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	newProduct := &models.Product{
		Name:        "New Product",
		Description: "New Description",
		Price:       79.99,
		ImageURL:    "http://example.com/newimage.png",
		SKU:         "SKUNEW",
	}
	expectedID := int64(3)

	// Настраиваем ожидания мок репозитория с установкой ID
	mockRepo.On("CreateProduct", mock.Anything, mock.AnythingOfType("*models.Product")).
		Run(func(args mock.Arguments) {
			p := args.Get(1).(*models.Product)
			p.ID = expectedID
		}).
		Return(expectedID, nil)

	// Вызываем метод сервиса
	createdID, err := svc.CreateProduct(context.Background(), newProduct)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedID, createdID)
	assert.Equal(t, expectedID, newProduct.ID)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestUpdateProduct(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	updatedProduct := models.Product{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       89.99,
		ImageURL:    "http://example.com/updatedimage.png",
		SKU:         "SKU001U",
	}

	// Настраиваем ожидания мок репозитория
	mockRepo.On("UpdateProduct", mock.Anything, &updatedProduct).Return(nil)

	// Вызываем метод сервиса
	err := svc.UpdateProduct(context.Background(), &updatedProduct)

	// Проверяем результаты
	assert.NoError(t, err)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestDeleteProduct(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	productID := int64(1)

	// Настраиваем ожидания мок репозитория
	mockRepo.On("DeleteProduct", mock.Anything, productID).Return(nil)

	// Вызываем метод сервиса
	err := svc.DeleteProduct(context.Background(), productID)

	// Проверяем результаты
	assert.NoError(t, err)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestSearchPacket_WithQuery(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	searchString := "Test"
	limit := int64(10)
	offset := int64(0)
	expectedPackages := []models.Package{
		{
			ID:          1,
			PackageName: "Test Package 1",
			Description: "Description 1",
			Content: []models.PackageContent{
				{ProductID: 1, Quantity: 2},
				{ProductID: 2, Quantity: 3},
			},
		},
	}

	// Настраиваем ожидания мок репозитория
	mockRepo.On("SearchPacket", mock.Anything, searchString, limit, offset).Return(expectedPackages, nil)

	// Вызываем метод сервиса
	packages, err := svc.SearchPacket(context.Background(), searchString, limit, offset)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedPackages, packages)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestSearchPacket_EmptyQuery(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	searchString := ""
	limit := int64(5)
	offset := int64(0)
	expectedPackages := []models.Package{
		{
			ID:          1,
			PackageName: "Package 1",
			Description: "Description 1",
			Content: []models.PackageContent{
				{ProductID: 1, Quantity: 1},
			},
		},
		{
			ID:          2,
			PackageName: "Package 2",
			Description: "Description 2",
			Content: []models.PackageContent{
				{ProductID: 2, Quantity: 2},
			},
		},
	}

	// Настраиваем ожидания мок репозитория
	mockRepo.On("GetAllPackages", mock.Anything, limit, offset).Return(expectedPackages, nil)

	// Вызываем метод сервиса
	packages, err := svc.SearchPacket(context.Background(), searchString, limit, offset)

	// Проверяем результаты
	assert.NoError(t, err)
	assert.Equal(t, expectedPackages, packages)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}

func TestSearchPacket_Error(t *testing.T) {
	// Создаем мок репозитория
	mockRepo := new(mocks.MockGoodsRepository)
	logger := log.NewNopLogger()
	svc := NewService(mockRepo, logger)

	// Определяем входные данные и ожидаемые результаты
	searchString := "NonExistent"
	limit := int64(10)
	offset := int64(0)
	expectedError := errors.New("database error")

	// Настраиваем ожидания мок репозитория
	mockRepo.On("SearchPacket", mock.Anything, searchString, limit, offset).
		Return([]models.Package(nil), expectedError)

	// Вызываем метод сервиса
	packages, err := svc.SearchPacket(context.Background(), searchString, limit, offset)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Nil(t, packages)

	// Проверяем, что ожидания были выполнены
	mockRepo.AssertExpectations(t)
}
