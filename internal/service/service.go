package service

import (
	"context"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Service описывает сервис для работы с продуктами и шаблонами.
type Service interface {
	// GetAllProducts возвращает список всех продуктов.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// GetProductByID возвращает продукт по его ID.
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	// SearchTemplates ищет шаблоны продуктов по их имени или ID с пагинацией.
	SearchTemplates(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Template, error)
	// AddTemplate добавляет новый шаблон продуктов в базу данных.
	AddTemplate(ctx context.Context, template *models.Template) (int64, error)
	// GetTemplateByID возвращает шаблон продуктов по его ID.
	GetTemplateByID(ctx context.Context, id int64) (models.Template, error)
	// CreateProduct добавляет новый продукт в базу данных.
	CreateProduct(ctx context.Context, p *models.Product) (int64, error)
	// UpdateProduct обновляет информацию о продукте в базе данных.
	UpdateProduct(ctx context.Context, p *models.Product) error
	// DeleteProduct удаляет продукт из базы данных.
	DeleteProduct(ctx context.Context, id int64) error
}

// GoodsService реализует интерфейс Service.
type GoodsService struct {
	repo models.GoodsRepository
	log  log.Logger
}

// NewService создает новый экземпляр Service.
func NewService(repo models.GoodsRepository, logger log.Logger) Service {
	return &GoodsService{
		repo: repo,
		log:  logger,
	}
}

// GetAllProducts возвращает список всех продуктов.
func (s *GoodsService) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	logger := log.With(s.log, "method", "GetAllProducts")
	products, err := s.repo.GetAllProducts(ctx)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	return products, nil
}

// GetProductByID возвращает продукт по его ID.
func (s *GoodsService) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	logger := log.With(s.log, "method", "GetProductByID")
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Product{}, err
	}
	return product, nil
}

// SearchTemplates ищет шаблоны продуктов по их имени или ID с пагинацией.
func (s *GoodsService) SearchTemplates(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Template, error) {
	logger := log.With(s.log, "method", "SearchTemplates")

	if searchString == "" {
		// Пустая строка поиска, возвращаем все шаблоны с пагинацией
		templates, err := s.repo.GetAllTemplates(ctx, limit, offset)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return nil, err
		}
		return templates, nil
	}

	// Поиск шаблонов по строке
	templates, err := s.repo.SearchTemplates(ctx, searchString, limit, offset)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}

	return templates, nil
}

// AddTemplate добавляет новый шаблон продуктов в базу данных.
func (s *GoodsService) AddTemplate(ctx context.Context, template *models.Template) (int64, error) {
	logger := log.With(s.log, "method", "AddTemplate")
	err := s.repo.CreateTemplate(ctx, template)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return template.ID, nil
}

// GetTemplateByID возвращает шаблон продуктов по его ID.
func (s *GoodsService) GetTemplateByID(ctx context.Context, id int64) (models.Template, error) {
	logger := log.With(s.log, "method", "GetTemplateByID")
	template, err := s.repo.GetTemplateByID(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Template{}, err
	}
	return template, nil
}

// CreateProduct добавляет новый продукт в базу данных.
func (s *GoodsService) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	logger := log.With(s.log, "method", "CreateProduct")
	productID, err := s.repo.CreateProduct(ctx, p)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return productID, nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (s *GoodsService) UpdateProduct(ctx context.Context, p *models.Product) error {
	logger := log.With(s.log, "method", "UpdateProduct")
	err := s.repo.UpdateProduct(ctx, p)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных.
func (s *GoodsService) DeleteProduct(ctx context.Context, id int64) error {
	logger := log.With(s.log, "method", "DeleteProduct")
	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	return nil
}
