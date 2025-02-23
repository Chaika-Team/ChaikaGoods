package postgresql_test

import (
	"context"
	"errors"
	"testing"

	"github.com/go-kit/log"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/Chaika-Team/ChaikaGoods/internal/repository/postgresql"
)

// MockRow для мокирования ответа QueryRow
type MockRow struct {
	mock.Mock
}

func (m *MockRow) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func TestGetProductByID(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)

	repo := postgresql.NewGoodsRepository(mockClient, logger)

	productID := int64(1)
	expectedProduct := models.Product{
		ID:          productID,
		Name:        "Test Product",
		Description: "Test Description",
		Price:       100.0,
		ImageURL:    "http://test.com/image.jpg",
		SKU:         "SKU123",
	}

	// Создаём мок для QueryRow
	mockRow := new(MockRow)
	mockClient.On("QueryRow", mock.Anything, mock.Anything, productID).Return(mockRow)
	mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*float64"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			*(args[0].(*int64)) = expectedProduct.ID
			*(args[1].(*string)) = expectedProduct.Name
			*(args[2].(*string)) = expectedProduct.Description
			*(args[3].(*float64)) = expectedProduct.Price
			*(args[4].(*string)) = expectedProduct.ImageURL
			*(args[5].(*string)) = expectedProduct.SKU
		}).
		Return(nil)

	// Выполняем тест
	product, err := repo.GetProductByID(ctx, productID)

	assert.NoError(t, err)
	assert.Equal(t, expectedProduct, product)

	// Проверяем вызовы
	mockClient.AssertExpectations(t)
	mockRow.AssertExpectations(t)
}

// Тест 2 --------------------------------------------------------------------------------

// MockRows для мокирования pgx.Rows
type MockRows struct {
	mock.Mock
}

func (m *MockRows) Next() bool {
	args := m.Called()
	return args.Bool(0)
}

func (m *MockRows) Scan(dest ...interface{}) error {
	args := m.Called(dest...)
	return args.Error(0)
}

func (m *MockRows) Err() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockRows) CommandTag() pgconn.CommandTag {
	return pgconn.CommandTag{}
}

func (m *MockRows) Conn() *pgx.Conn {
	return nil
}

func (m *MockRows) Close() {}
func (m *MockRows) FieldDescriptions() []pgconn.FieldDescription {
	return nil
}
func (m *MockRows) RawValues() [][]byte {
	return nil
}

// Исправление: добавляем метод Values()
func (m *MockRows) Values() ([]interface{}, error) {
	args := m.Called()
	return args.Get(0).([]interface{}), args.Error(1)
}

func TestGetAllProducts(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)
	repo := postgresql.NewGoodsRepository(mockClient, logger)

	expectedProducts := []models.Product{
		{ID: 1, Name: "Product 1", Description: "Desc 1", Price: 10.0, ImageURL: "http://test.com/1.jpg", SKU: "SKU1"},
		{ID: 2, Name: "Product 2", Description: "Desc 2", Price: 20.0, ImageURL: "http://test.com/2.jpg", SKU: "SKU2"},
	}

	// Создаем мокированные строки
	mockRows := new(MockRows)
	mockClient.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

	// Ожидаем вызовы Next()
	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*float64"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			*(args[0].(*int64)) = expectedProducts[0].ID
			*(args[1].(*string)) = expectedProducts[0].Name
			*(args[2].(*string)) = expectedProducts[0].Description
			*(args[3].(*float64)) = expectedProducts[0].Price
			*(args[4].(*string)) = expectedProducts[0].ImageURL
			*(args[5].(*string)) = expectedProducts[0].SKU
		}).
		Return(nil).Once()

	mockRows.On("Next").Return(true).Once()
	mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*float64"),
		mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
		Run(func(args mock.Arguments) {
			*(args[0].(*int64)) = expectedProducts[1].ID
			*(args[1].(*string)) = expectedProducts[1].Name
			*(args[2].(*string)) = expectedProducts[1].Description
			*(args[3].(*float64)) = expectedProducts[1].Price
			*(args[4].(*string)) = expectedProducts[1].ImageURL
			*(args[5].(*string)) = expectedProducts[1].SKU
		}).
		Return(nil).Once()

	mockRows.On("Next").Return(false).Once()
	mockRows.On("Err").Return(nil)

	// Выполняем тест
	products, err := repo.GetAllProducts(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expectedProducts, products)

	// Проверяем вызовы
	mockClient.AssertExpectations(t)
	mockRows.AssertExpectations(t)
}

// Тест 3 --------------------------------------------------------------------------------

func TestCreateProduct(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)
	repo := postgresql.NewGoodsRepository(mockClient, logger)

	product := &models.Product{
		Name:        "New Product",
		Description: "Description",
		Price:       99.99,
		ImageURL:    "http://example.com/image.jpg",
		SKU:         "SKU123",
	}

	t.Run("Успешное создание продукта", func(t *testing.T) {
		mockRow := new(MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = 1
			}).
			Return(nil)

		id, err := repo.CreateProduct(ctx, product)

		assert.NoError(t, err)
		assert.Equal(t, int64(1), id)
	})

	t.Run("Ошибка UniqueViolation (SKU уже существует)", func(t *testing.T) {
		mockRow := new(MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

		_, err := repo.CreateProduct(ctx, product)

		assert.Error(t, err)
		assert.True(t, myerr.IsConflict(err))
	})

	t.Run("Ошибка БД", func(t *testing.T) {
		mockRow := new(MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Return(errors.New("db error"))

		_, err := repo.CreateProduct(ctx, product)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error") // Вместо "Failed to create product"
	})
}

// Тест 4 --------------------------------------------------------------------------------

// MockResult для мокирования CommandTag (результата Exec)
type MockResult struct {
	mock.Mock
}

func (m *MockResult) RowsAffected() int64 {
	args := m.Called()
	return args.Get(0).(int64)
}

func TestUpdateProduct(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)
	repo := postgresql.NewGoodsRepository(mockClient, logger)

	product := &models.Product{
		ID:          1,
		Name:        "Updated Product",
		Description: "Updated Description",
		Price:       150.0,
		ImageURL:    "http://example.com/updated.jpg",
		SKU:         "SKU123",
	}

	t.Run("Успешное обновление продукта", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID).
			Return(pgconn.NewCommandTag("UPDATE 1"), nil).Once()

		err := repo.UpdateProduct(ctx, product)

		assert.NoError(t, err)
	})

	t.Run("Ошибка UniqueViolation (SKU уже существует)", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID).
			Return(pgconn.NewCommandTag("UPDATE 0"), &pgconn.PgError{Code: pgerrcode.UniqueViolation}).Once()

		err := repo.UpdateProduct(ctx, product)

		assert.Error(t, err)
		assert.True(t, myerr.IsConflict(err))
	})

	t.Run("Ошибка NotFound (Продукт не найден)", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU, product.ID).
			Return(pgconn.NewCommandTag("UPDATE 0"), nil).Once()

		err := repo.UpdateProduct(ctx, product)

		assert.Error(t, err)
		assert.True(t, myerr.IsNotFound(err))
	})

	// Проверяем вызовы
	mockClient.AssertExpectations(t)
}
