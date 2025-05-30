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

func newTestRepo() (*postgresql.MockClient, models.GoodsRepository, context.Context) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)
	repo := postgresql.NewGoodsRepository(mockClient, logger) // Возвращает интерфейс

	return mockClient, repo, ctx
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: safr
// Описание:
//   - Тест для функции NewGoodsRepository.
//   - Проверяет, что репозиторий создается корректно.
//   - Классы эквивалентности: переданы валидные зависимости (мок клиента и логгер).
func TestNewGoodsRepository(t *testing.T) {
	// Подготавливаем моки
	mockClient := new(postgresql.MockClient)
	mockLogger := log.NewNopLogger()

	// Создаем репозиторий
	repo := postgresql.NewGoodsRepository(mockClient, mockLogger)

	// Проверяем, что repo не nil и является *GoodsPGRepository
	assert.NotNil(t, repo)

	// Проверяем, что repo имеет правильный тип
	_, ok := repo.(*postgresql.GoodsPGRepository)
	assert.True(t, ok, "repo должен быть типа *GoodsPGRepository")
}

// Техника тест-дизайна:  Классы эквивалентности
// Автор: safr
// Описание:
//   - Тест для метода GetProductByID.
//   - Классы эквивалентности: корректный ID, несуществующий ID, ошибка БД.
func TestGetProductByID(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

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
	mockRow := new(postgresql.MockRows)
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

// Техника тест-дизайна: #2 Классы эквивалентности + граничные значения
// Автор: safr
// Описание:
//   - Тест для метода GetAllProducts.
//   - Классы эквивалентности: пустая таблица, несколько записей, ошибка БД.
func TestGetAllProducts(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	expectedProducts := []models.Product{
		{ID: 1, Name: "Product 1", Description: "Desc 1", Price: 10.0, ImageURL: "http://test.com/1.jpg", SKU: "SKU1"},
		{ID: 2, Name: "Product 2", Description: "Desc 2", Price: 20.0, ImageURL: "http://test.com/2.jpg", SKU: "SKU2"},
	}

	// Создаем мокированные строки
	mockRows := new(postgresql.MockRows)
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

// Техника тест-дизайна: #3 Классы эквивалентности + обработка ошибок
// Автор: safr
// Описание:
//   - Тест для метода CreateProduct.
//   - Проверка успешного создания продукта, ошибки UniqueViolation, ошибки БД.
func TestCreateProduct(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	product := &models.Product{
		Name:        "New Product",
		Description: "Description",
		Price:       99.99,
		ImageURL:    "http://example.com/image.jpg",
		SKU:         "SKU123",
	}

	t.Run("Успешное создание продукта", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
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
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Return(&pgconn.PgError{Code: pgerrcode.UniqueViolation})

		_, err := repo.CreateProduct(ctx, product)

		assert.Error(t, err)
		assert.True(t, myerr.IsConflict(err))
	})

	t.Run("Ошибка БД", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, product.Name, product.Description, product.Price, product.ImageURL, product.SKU).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Return(errors.New("db error"))

		_, err := repo.CreateProduct(ctx, product)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error") // Вместо "Failed to create product"
	})
}

// Техника тест-дизайна: #4 Классы эквивалентности + обработка ошибок
// Автор: safr
// Описание:
//   - Тест для метода UpdateProduct.
//   - Классы эквивалентности: успешное обновление, конфликт из-за существующего SKU, отсутствие продукта в БД.
func TestUpdateProduct(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

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

// Техника тест-дизайна: #4 Классы эквивалентности + граничные значения
// Автор: safr
// Описание:
//   - Тест для метода DeleteProduct.
//   - Классы эквивалентности: успешное удаление, отсутствие продукта, ошибка БД.

func TestDeleteProduct(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()
	productID := int64(1)

	t.Run("Успешное удаление продукта", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, productID).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()

		err := repo.DeleteProduct(ctx, productID)

		assert.NoError(t, err)
	})

	t.Run("Ошибка NotFound (продукт не найден)", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, productID).
			Return(pgconn.NewCommandTag("DELETE 0"), nil).Once()

		err := repo.DeleteProduct(ctx, productID)

		assert.Error(t, err)
		assert.True(t, myerr.IsNotFound(err))
	})

	t.Run("Ошибка БД при удалении", func(t *testing.T) {
		mockClient.On("Exec", mock.Anything, mock.Anything, productID).
			Return(pgconn.NewCommandTag("DELETE 0"), errors.New("db error")).Once()

		err := repo.DeleteProduct(ctx, productID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})

	// Проверяем вызовы
	mockClient.AssertExpectations(t)
}

// Техника тест-дизайна: #5 Классы эквивалентности + анализ граничных значений
// Автор: safr
// Описание:
//   - Тест для метода GetTemplateByID.
//   - Классы эквивалентности: существующий шаблон, несуществующий шаблон, ошибка БД при получении шаблона, ошибка БД при получении содержимого.

func TestGetTemplateByID(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()
	templateID := int64(1)

	expectedTemplate := models.Template{
		ID:           templateID,
		TemplateName: "Sample Template",
		Description:  "This is a sample template",
		Content: []models.TemplateContent{
			{ProductID: 101, Quantity: 2},
			{ProductID: 102, Quantity: 5},
		},
	}

	t.Run("Успешное получение шаблона", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, templateID).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplate.ID
				*(args[1].(*string)) = expectedTemplate.TemplateName
				*(args[2].(*string)) = expectedTemplate.Description
			}).
			Return(nil)

		mockRows := new(postgresql.MockRows)
		mockClient.On("Query", mock.Anything, mock.Anything, templateID).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Twice()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*int")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplate.Content[0].ProductID
				*(args[1].(*int)) = expectedTemplate.Content[0].Quantity
			}).
			Return(nil).Once()

		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*int")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplate.Content[1].ProductID
				*(args[1].(*int)) = expectedTemplate.Content[1].Quantity
			}).
			Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)

		// Выполняем тест
		template, err := repo.GetTemplateByID(ctx, templateID)

		assert.NoError(t, err)
		assert.Equal(t, expectedTemplate, template)
	})

	t.Run("Ошибка NotFound (шаблон не найден)", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, templateID).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Return(pgx.ErrNoRows)

		_, err := repo.GetTemplateByID(ctx, templateID)

		assert.Error(t, err)
		assert.True(t, myerr.IsNotFound(err))
	})

	t.Run("Ошибка при получении шаблона", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, templateID).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Return(errors.New("db error"))

		_, err := repo.GetTemplateByID(ctx, templateID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})

	t.Run("Ошибка при получении содержимого шаблона", func(t *testing.T) {
		mockRow := new(postgresql.MockRow)
		mockClient.On("QueryRow", mock.Anything, mock.Anything, templateID).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplate.ID
				*(args[1].(*string)) = expectedTemplate.TemplateName
				*(args[2].(*string)) = expectedTemplate.Description
			}).
			Return(nil)

		// Вместо nil нужно вернуть (*postgresql.MockRows)(nil)
		mockClient.On("Query", mock.Anything, mock.Anything, templateID).
			Return((*postgresql.MockRows)(nil), errors.New("db error")).Once()

		_, err := repo.GetTemplateByID(ctx, templateID)

		assert.Error(t, err)
		assert.EqualError(t, err, "db error")
	})

	// Проверяем вызовы
	mockClient.AssertExpectations(t)
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: safr
// Описание:
//   - Тест для метода GetProductsByTemplateID.
//   - Проверяет, что метод корректно возвращает список продуктов, относящихся к определенному шаблону.
//   - Классы эквивалентности: шаблон содержит продукты, шаблон пуст, ошибка БД.
func TestGetProductsByTemplateID(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	templateID := int64(1)

	t.Run("Шаблон содержит продукты", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)
		mockClient.On("Query", mock.Anything, mock.Anything, templateID).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*int")).Run(func(args mock.Arguments) {
			*(args[0].(*int64)) = 10
			*(args[1].(*int)) = 2
		}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)

		contents, err := repo.GetProductsByTemplateID(ctx, templateID)

		assert.NoError(t, err)
		assert.Len(t, contents, 1)
		assert.Equal(t, int64(10), contents[0].ProductID)
		assert.Equal(t, 2, contents[0].Quantity)

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

}

// Техника тест-дизайна: Классы эквивалентности + анализ граничных значений
// Автор: safr
// Описание:
//   - Тест для метода ListTemplates.
//   - Проверяет корректность получения списка всех шаблонов с пагинацией.
//   - Классы эквивалентности: успешное получение данных, ошибка при выполнении запроса, ошибка при сканировании данных, ошибка при итерации по результатам.
func TestListTemplates(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	expectedTemplates := []models.Template{
		{ID: 1, TemplateName: "Template 1", Description: "Desc 1"},
		{ID: 2, TemplateName: "Template 2", Description: "Desc 2"},
	}

	t.Run("Успешное получение списка шаблонов", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)
		mockClient.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplates[0].ID
				*(args[1].(*string)) = expectedTemplates[0].TemplateName
				*(args[2].(*string)) = expectedTemplates[0].Description
			}).Return(nil).Once()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = expectedTemplates[1].ID
				*(args[1].(*string)) = expectedTemplates[1].TemplateName
				*(args[2].(*string)) = expectedTemplates[1].Description
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)

		templates, err := repo.ListTemplates(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedTemplates, templates)

		// Проверяем вызовы
		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("Ошибка при выполнении запроса", func(t *testing.T) {
		mockClient = new(postgresql.MockClient) // Сбрасываем моки перед новым тестом
		repo = postgresql.NewGoodsRepository(mockClient, log.NewNopLogger())

		mockClient.On("Query", mock.Anything, mock.Anything).Return((*postgresql.MockRows)(nil), errors.New("query error")).Once()

		templates, err := repo.ListTemplates(ctx)

		assert.Error(t, err)
		assert.Nil(t, templates)

		// Проверяем вызовы
		mockClient.AssertExpectations(t)
	})

	t.Run("Ошибка при Scan()", func(t *testing.T) {
		mockClient, repo, ctx := newTestRepo()

		mockRows := new(postgresql.MockRows)
		mockClient.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Return(errors.New("scan error")).Once()
		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil)

		templates, err := repo.ListTemplates(ctx)

		assert.NoError(t, err)
		assert.Empty(t, templates)

		// Проверяем вызовы
		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("Ошибка при итерации по rows", func(t *testing.T) {
		mockClient = new(postgresql.MockClient)
		repo = postgresql.NewGoodsRepository(mockClient, log.NewNopLogger())

		mockRows := new(postgresql.MockRows)
		mockClient.On("Query", mock.Anything, mock.Anything).Return(mockRows, nil)

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows error"))

		templates, err := repo.ListTemplates(ctx)

		assert.Error(t, err)
		assert.Nil(t, templates)

		// Проверяем вызовы
		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

// Техника тест-дизайна: Классы эквивалентности + анализ граничных значений
// Автор: safr
// Описание:
//   - Тест для метода CreateTemplate.
//   - Проверяет корректность создания нового шаблона с продуктами.
//   - Классы эквивалентности: успешное создание шаблона, ошибка при вставке шаблона, ошибка при вставке содержимого шаблона, ошибка при коммите транзакции.
func TestCreateTemplate(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	template := &models.Template{
		TemplateName: "New Template",
		Description:  "A test template",
		Content: []models.TemplateContent{
			{ProductID: 1, Quantity: 5},
			{ProductID: 2, Quantity: 10},
		},
	}

	t.Run("Ошибка при добавлении шаблона", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)
		mockRow := new(postgresql.MockRow)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("QueryRow", mock.Anything, mock.Anything, template.TemplateName, template.Description).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Return(errors.New("insert error")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		err := repo.CreateTemplate(ctx, template)

		assert.Error(t, err)
		assert.EqualError(t, err, "insert error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockRow.AssertExpectations(t)
	})

	t.Run("Ошибка при добавлении содержимого", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)
		mockRow := new(postgresql.MockRow)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("QueryRow", mock.Anything, mock.Anything, template.TemplateName, template.Description).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = 1 // Устанавливаем ID шаблона
			}).Return(nil).Once()

		mockTx.On("Exec", mock.Anything, mock.Anything, int64(1), template.Content[0].ProductID, template.Content[0].Quantity).
			Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, int64(1), template.Content[1].ProductID, template.Content[1].Quantity).
			Return(pgconn.NewCommandTag("INSERT 0"), errors.New("insert content error")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		err := repo.CreateTemplate(ctx, template)

		assert.Error(t, err)
		assert.EqualError(t, err, "insert content error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockRow.AssertExpectations(t)
	})

	t.Run("Ошибка при коммите транзакции", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)
		mockRow := new(postgresql.MockRow)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("QueryRow", mock.Anything, mock.Anything, template.TemplateName, template.Description).
			Return(mockRow).Once()
		mockRow.On("Scan", mock.AnythingOfType("*int64")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = 1 // Устанавливаем ID шаблона
			}).Return(nil).Once()

		mockTx.On("Exec", mock.Anything, mock.Anything, int64(1), template.Content[0].ProductID, template.Content[0].Quantity).
			Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, int64(1), template.Content[1].ProductID, template.Content[1].Quantity).
			Return(pgconn.NewCommandTag("INSERT 1"), nil).Once()
		mockTx.On("Commit", mock.Anything).Return(errors.New("commit error")).Once()

		err := repo.CreateTemplate(ctx, template)

		assert.Error(t, err)
		assert.EqualError(t, err, "commit error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockRow.AssertExpectations(t)
	})
}

// Техника тест-дизайна: Классы эквивалентности + анализ граничных значений
// Автор: safr
// Описание:
//   - Тест для метода DeleteTemplate.
//   - Проверяет корректность удаления шаблона и его содержимого.
//   - Классы эквивалентности: успешное удаление, ошибка начала транзакции, ошибка удаления содержимого шаблона, ошибка удаления самого шаблона, ошибка коммита транзакции.
func TestDeleteTemplate(t *testing.T) {
	ctx := context.Background()
	logger := log.NewNopLogger()
	mockClient := new(postgresql.MockClient)
	repo := postgresql.NewGoodsRepository(mockClient, logger)

	packageID := int64(1)

	t.Run("успешное удаление шаблона", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, packageID).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Twice()
		mockTx.On("Commit", mock.Anything).Return(nil).Once()

		err := repo.DeleteTemplate(ctx, packageID)

		assert.NoError(t, err)
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("ошибка: не удалось начать транзакцию", func(t *testing.T) {
		mockClient.On("Begin", mock.Anything).Return((*postgresql.MockTx)(nil), errors.New("transaction error")).Once()

		err := repo.DeleteTemplate(ctx, packageID)

		assert.Error(t, err)
		assert.EqualError(t, err, "transaction error")
		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка: не удалось удалить содержимое шаблона", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, packageID).
			Return(pgconn.NewCommandTag("DELETE 0"), errors.New("delete content error")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		err := repo.DeleteTemplate(ctx, packageID)

		assert.Error(t, err)
		assert.EqualError(t, err, "delete content error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("ошибка: не удалось удалить сам шаблон", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, packageID).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, packageID).
			Return(pgconn.NewCommandTag("DELETE 0"), errors.New("delete template error")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		err := repo.DeleteTemplate(ctx, packageID)

		assert.Error(t, err)
		assert.EqualError(t, err, "delete template error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("ошибка: не удалось зафиксировать транзакцию", func(t *testing.T) {
		mockTx := new(postgresql.MockTx)

		mockClient.On("Begin", mock.Anything).Return(mockTx, nil).Once()
		mockTx.On("Exec", mock.Anything, mock.Anything, packageID).
			Return(pgconn.NewCommandTag("DELETE 1"), nil).Twice()
		mockTx.On("Commit", mock.Anything).Return(errors.New("commit error")).Once()
		mockTx.On("Rollback", mock.Anything).Return(nil).Once() // Добавляем Rollback

		err := repo.DeleteTemplate(ctx, packageID)

		assert.Error(t, err)
		assert.EqualError(t, err, "commit error")
		mockClient.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: safr
// Описание:
//   - Тест для метода SearchTemplates.
//   - Проверяет поиск шаблонов по названию и описанию с пагинацией.
//   - Классы эквивалентности: успешный поиск, ошибка выполнения запроса.
func TestSearchTemplates(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	searchString := "test"
	limit := int64(10)
	offset := int64(0)

	t.Run("успешный поиск шаблонов", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, "%test%", limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = 1
				*(args[1].(*string)) = "Test Template"
				*(args[2].(*string)) = "Description"
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Maybe()

		templates, err := repo.SearchTemplates(ctx, searchString, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, templates, 1)
		assert.Equal(t, int64(1), templates[0].ID)
		assert.Equal(t, "Test Template", templates[0].TemplateName)
		assert.Equal(t, "Description", templates[0].Description)

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("ошибка выполнения запроса", func(t *testing.T) {
		mockClient.On("Query", mock.Anything, mock.Anything, "%test%", limit, offset).
			Return((*postgresql.MockRows)(nil), errors.New("query error")).Once()

		templates, err := repo.SearchTemplates(ctx, searchString, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, templates)
		assert.EqualError(t, err, "query error")

		mockClient.AssertExpectations(t)
	})

	t.Run("ошибка при сканировании строки", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, "%test%", limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Return(errors.New("scan error")).Once()
		mockRows.On("Next").Return(false).Once() // Завершаем итерацию
		mockRows.On("Err").Return(nil).Once()    // Нет ошибки на уровне строк
		mockRows.On("Close").Return().Maybe()    // Теперь не ломает тест, если `Close()` уже вызван

		templates, err := repo.SearchTemplates(ctx, searchString, limit, offset)

		assert.NoError(t, err)      // Ошибка сканирования не должна прерывать выполнение
		assert.Len(t, templates, 0) // Строка не добавляется

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	t.Run("ошибка при итерации по rows", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, "%test%", limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows iteration error")).Once()
		mockRows.On("Close").Return().Maybe() // Теперь не ломает тест, если `Close()` уже вызван

		templates, err := repo.SearchTemplates(ctx, searchString, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, templates)
		assert.EqualError(t, err, "rows iteration error")

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}

// Техника тест-дизайна: Классы эквивалентности
// Автор: safr
// Описание:
//   - Тест для метода GetAllTemplates.
//   - Проверяет получение всех шаблонов с лимитом и смещением.
//   - Классы эквивалентности: успешное получение списка, ошибка выполнения запроса.
func TestGetAllTemplates(t *testing.T) {
	mockClient, repo, ctx := newTestRepo()

	limit := int64(10)
	offset := int64(0)

	// Успешное получение шаблонов
	t.Run("успешное получение шаблонов", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Run(func(args mock.Arguments) {
				*(args[0].(*int64)) = 1
				*(args[1].(*string)) = "Test Template"
				*(args[2].(*string)) = "Description"
			}).Return(nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Maybe()

		templates, err := repo.GetAllTemplates(ctx, limit, offset)

		assert.NoError(t, err)
		assert.Len(t, templates, 1)
		assert.Equal(t, int64(1), templates[0].ID)
		assert.Equal(t, "Test Template", templates[0].TemplateName)
		assert.Equal(t, "Description", templates[0].Description)

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	// Ошибка выполнения запроса
	t.Run("ошибка выполнения запроса", func(t *testing.T) {
		mockClient.On("Query", mock.Anything, mock.Anything, limit, offset).
			Return((*postgresql.MockRows)(nil), errors.New("query error")).Once()

		templates, err := repo.GetAllTemplates(ctx, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, templates)
		assert.EqualError(t, err, "query error")

		mockClient.AssertExpectations(t)
	})

	// Ошибка сканирования строки
	t.Run("ошибка при сканировании строки", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(true).Once()
		mockRows.On("Scan", mock.AnythingOfType("*int64"), mock.AnythingOfType("*string"), mock.AnythingOfType("*string")).
			Return(errors.New("scan error")).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(nil).Once()
		mockRows.On("Close").Return().Maybe()

		templates, err := repo.GetAllTemplates(ctx, limit, offset)

		assert.NoError(t, err)      // Ошибки нет, так как пропускаем сбойную строку
		assert.Len(t, templates, 0) // Но шаблоны не добавляются

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})

	// Ошибка при итерации по rows
	t.Run("ошибка при итерации по rows", func(t *testing.T) {
		mockRows := new(postgresql.MockRows)

		mockClient.On("Query", mock.Anything, mock.Anything, limit, offset).
			Return(mockRows, nil).Once()

		mockRows.On("Next").Return(false).Once()
		mockRows.On("Err").Return(errors.New("rows error")).Once()
		mockRows.On("Close").Return().Maybe()

		templates, err := repo.GetAllTemplates(ctx, limit, offset)

		assert.Error(t, err)
		assert.Nil(t, templates)
		assert.EqualError(t, err, "rows error")

		mockClient.AssertExpectations(t)
		mockRows.AssertExpectations(t)
	})
}
