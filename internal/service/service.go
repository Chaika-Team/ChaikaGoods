package service

import (
	"context"
	"fmt"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	repo "github.com/Chaika-Team/ChaikaGoods/internal/repository"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

// Service описывает сервис для работы с продуктами и пакетами.
type Service interface {
	// GetAllProducts возвращает список всех продуктов.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// GetProductByID возвращает продукт по его ID.
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	// SearchPackages ищет пакеты продуктов по их имени или ID с пагинацией.
	SearchPackages(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Package, error)
	// AddPackage добавляет новый пакет продуктов в базу данных.
	AddPackage(ctx context.Context, pkg *models.Package) (int64, error)
	// GetPackageByID возвращает пакет продуктов по его ID.
	GetPackageByID(ctx context.Context, id int64) (models.Package, error)
	// CreateProduct добавляет новый продукт в базу данных.
	CreateProduct(ctx context.Context, p *models.Product) (int64, error)
	// UpdateProduct обновляет информацию о продукте в базе данных.
	UpdateProduct(ctx context.Context, p *models.Product) error
	// DeleteProduct удаляет продукт из базы данных.
	DeleteProduct(ctx context.Context, id int64) error
}

// GoodsService реализует интерфейс Service.
type GoodsService struct {
	repo repo.GoodsRepository
	log  log.Logger
}

// NewService создает новый экземпляр Service.
func NewService(repo repo.GoodsRepository, logger log.Logger) Service {
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
		appErr := myerr.ToAppError(logger, err, "Failed to retrieve all products")
		_ = level.Error(logger).Log("err", appErr)
		return nil, appErr
	}
	return products, nil
}

// GetProductByID возвращает продукт по его ID.
func (s *GoodsService) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	logger := log.With(s.log, "method", "GetProductByID")
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, fmt.Sprintf("Failed to retrieve product with ID %d", id))
		_ = level.Error(logger).Log("err", appErr)
		return models.Product{}, appErr
	}
	return product, nil
}

// SearchPackages ищет пакеты продуктов по их имени или ID с пагинацией.
func (s *GoodsService) SearchPackages(ctx context.Context, searchString string, limit int64, offset int64) ([]models.Package, error) {
	logger := log.With(s.log, "method", "SearchPackages")

	if searchString == "" {
		// Пустая строка поиска, возвращаем все пакеты с пагинацией
		packages, err := s.repo.GetAllPackages(ctx, limit, offset)
		if err != nil {
			appErr := myerr.ToAppError(logger, err, "Failed to retrieve all packages")
			_ = level.Error(logger).Log("err", appErr)
			return nil, appErr
		}
		return packages, nil
	}

	// Поиск пакетов по строке
	packages, err := s.repo.SearchPackages(ctx, searchString, limit, offset)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, fmt.Sprintf("Failed to search packages with query '%s'", searchString))
		_ = level.Error(logger).Log("err", appErr)
		return nil, appErr
	}

	return packages, nil
}

// AddPackage добавляет новый пакет продуктов в базу данных.
func (s *GoodsService) AddPackage(ctx context.Context, pkg *models.Package) (int64, error) {
	logger := log.With(s.log, "method", "AddPackage")
	err := s.repo.CreatePackage(ctx, pkg)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, "Failed to create package")
		_ = level.Error(logger).Log("err", appErr)
		return 0, appErr
	}
	return pkg.ID, nil
}

// GetPackageByID возвращает пакет продуктов по его ID.
func (s *GoodsService) GetPackageByID(ctx context.Context, id int64) (models.Package, error) {
	logger := log.With(s.log, "method", "GetPackageByID")
	pkg, err := s.repo.GetPackageByID(ctx, id)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, fmt.Sprintf("Failed to retrieve package with ID %d", id))
		_ = level.Error(logger).Log("err", appErr)
		return models.Package{}, appErr
	}
	return pkg, nil
}

// CreateProduct добавляет новый продукт в базу данных.
func (s *GoodsService) CreateProduct(ctx context.Context, p *models.Product) (int64, error) {
	logger := log.With(s.log, "method", "CreateProduct")
	productID, err := s.repo.CreateProduct(ctx, p)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, "Failed to create product")
		_ = level.Error(logger).Log("err", appErr)
		return 0, appErr
	}
	return productID, nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (s *GoodsService) UpdateProduct(ctx context.Context, p *models.Product) error {
	logger := log.With(s.log, "method", "UpdateProduct")
	err := s.repo.UpdateProduct(ctx, p)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, fmt.Sprintf("Failed to update product with ID %d", p.ID))
		_ = level.Error(logger).Log("err", appErr)
		return appErr
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных.
func (s *GoodsService) DeleteProduct(ctx context.Context, id int64) error {
	logger := log.With(s.log, "method", "DeleteProduct")
	err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		appErr := myerr.ToAppError(logger, err, fmt.Sprintf("Failed to delete product with ID %d", id))
		_ = level.Error(logger).Log("err", appErr)
		return appErr
	}
	return nil
}
