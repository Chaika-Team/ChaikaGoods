package service

import (
	"ChaikaGoods/internal/models"
	repo "ChaikaGoods/internal/repository/postgresql"
	"context"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

/* Бизнес задачи для микросервиса продуктов продаваемых в вагонах поезда:
Со стороны проводника поезда:
1. Проводник запрашивает синхронизацию своей базы данных с центральной базой данных продуктов.
2. Проводник запрашивает текущую версию базы продуктов, чтобы проверить актуальность своей базы.
3. Проводник ищет пакеты продуктов по их имени или ID.
4. Проводник добавляет новый пакет продуктов в базу данных.
5. Проводник обновляет информацию о пакете продуктов в базе данных.
Со стороны администратора системы:
1. Администратор может добавить новый продукт в базу данных.
2. Администратор может обновить информацию о продукте в базе данных.
3. Администратор может удалить продукт из базы данных.
*/

// GoodsService описывает сервис для работы с продуктами.
type GoodsService interface {
	// GetAllProducts возвращает список всех продуктов.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// GetCurrentVersion возвращает текущую версию базы данных продуктов.
	GetCurrentVersion(ctx context.Context) (int, error)
	// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
	GetDelta(ctx context.Context, version int) ([]models.Product, error)
	// SearchPacket ищет пакеты продуктов по их имени или ID.
	SearchPacket(ctx context.Context, searchString string) ([]models.Product, error)
	// AddPacket добавляет новый пакет продуктов в базу данных.
	AddPacket(ctx context.Context, name string, description string, packageContent []models.PackageContent) error
	// AddProduct добавляет новый продукт в базу данных.
	AddProduct(ctx context.Context, name string) error
	// UpdateProduct обновляет информацию о продукте в базе данных.
	UpdateProduct(ctx context.Context, product *models.Product) error
	// DeleteProduct удаляет продукт из базы данных.
	DeleteProduct(ctx context.Context, id int64) error
}
type Service struct {
	repo repo.GoodsRepository
	log  log.Logger
}

// NewGoodsService создает новый экземпляр GoodsService.
func NewGoodsService(repo repo.GoodsRepository, logger log.Logger) *Service {
	return &Service{
		repo: repo,
		log:  logger,
	}
}

// GetAllProducts возвращает список всех продуктов.
func (s *Service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	logger := log.With(s.log, "method", "GetAllProducts")
	if products, err := s.repo.GetAllProducts(ctx); err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	} else {
		return products, nil
	}
}

// GetCurrentVersion возвращает текущую версию базы данных продуктов.
func (s *Service) GetCurrentVersion(ctx context.Context) (int, error) {
	//TODO: implement
	return 0, nil
}

// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
func (s *Service) GetDelta(ctx context.Context, version int) ([]models.Product, error) {
	//TODO: implement
	return nil, nil
}

// SearchPacket ищет пакеты продуктов по их имени или ID.
func (s *Service) SearchPacket(ctx context.Context, searchString string) ([]models.Product, error) {
	//TODO: implement
	return nil, nil
}

// AddPacket добавляет новый пакет продуктов в базу данных.
func (s *Service) AddPacket(ctx context.Context, name string, description string, packageContent []models.PackageContent) error {
	logger := log.With(s.log, "method", "AddPacket")
	packet := models.Package{
		PackageName: name,
		Description: description,
	}
	if err := s.repo.CreatePackage(ctx, &packet); err != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	// add products to the package
	if err := s.repo.AddProductToPackage(ctx, packet.ID, packageContent); err != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	return nil
}
