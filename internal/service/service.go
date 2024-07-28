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
4. Все изменения продуктов происходят через CQRS и создание версий базы продуктов.
*/

// GoodsService описывает сервис для работы с продуктами.
type GoodsService interface {
	// GetAllProducts возвращает список всех продуктов.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// GetCurrentVersion возвращает текущую версию базы данных продуктов.
	GetCurrentVersion(ctx context.Context) (models.Version, error)
	// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
	GetDelta(ctx context.Context, version int) ([]models.Product, error)
	// SearchPacket ищет пакеты продуктов по их имени или ID.
	SearchPacket(ctx context.Context, searchString string) ([]models.Product, error)
	// AddPacket добавляет новый пакет продуктов в базу данных.
	AddPacket(ctx context.Context, name string, description string, packageContent []models.PackageContent) error
	// AddProduct добавляет новый продукт в базу данных.
	AddProduct(ctx context.Context, data *map[string]interface{}) error
	// UpdateProduct обновляет информацию о продукте в базе данных.
	UpdateProduct(ctx context.Context, data *map[string]interface{}) error
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

func (s *Service) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	logger := log.With(s.log, "method", "GetProductByID")
	if product, err := s.repo.GetProductByID(ctx, id); err != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Product{}, err
	} else {
		return product, nil
	}
}

// GetCurrentVersion возвращает текущую версию базы данных продуктов.
func (s *Service) GetCurrentVersion(ctx context.Context) (models.Version, error) {
	logger := log.With(s.log, "method", "GetCurrentVersion")
	if version, err := s.repo.GetCurrentActualVersion(ctx); err != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Version{}, err
	} else {
		return version, nil
	}
}

// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
func (s *Service) GetDelta(ctx context.Context, version int) ([]models.Change, error) {
	logger := log.With(s.log, "method", "GetDelta")
	// get actual version
	actualVersion, err := s.repo.GetCurrentActualVersion(ctx)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err

	}
	// get all versions between the device version and the actual version
	versions, err := s.repo.GetVersionsBetween(ctx, version, actualVersion.VersionID)
	if err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	// get all changes between the device version and the actual version
	changes := make([]models.Change, 0)
	for _, v := range versions {
		c, err := s.repo.GetAllChanges(ctx, v)
		if err != nil {
			_ = level.Error(logger).Log("err", err)
			return nil, err
		}
		changes = append(changes, c...)
	}
	return changes, nil
}

// SearchPacket ищет пакеты продуктов по их имени или ID.
func (s *Service) SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error) {
	logger := log.With(s.log, "method", "SearchPacket")
	if packages, err := s.repo.SearchPacket(ctx, searchString, quantity, offset); err != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	} else {
		return packages, nil
	}
}

// AddPacket добавляет новый пакет продуктов в базу данных.
func (s *Service) AddPacket(ctx context.Context, packet *models.Package, packageContent []models.PackageContent) (int64, error) {
	logger := log.With(s.log, "method", "AddPacket")
	if err := s.repo.CreatePackage(ctx, packet); err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	// add products to the package
	if err := s.repo.AddProductToPackage(ctx, packet.ID, packageContent); err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return packet.ID, nil
}

// AddProduct добавляет новый продукт в базу данных.
func (s *Service) AddProduct(ctx context.Context, p *map[string]interface{}) (changeID int64, err error) {
	logger := log.With(s.log, "method", "AddProduct")
	if changeID, err = s.repo.AddQueryToCreateProduct(ctx, p); err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (s *Service) UpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error) {
	logger := log.With(s.log, "method", "UpdateProduct")
	if changeID, err = s.repo.AddQueryToUpdateProduct(ctx, data); err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}

// DeleteProduct удаляет продукт из базы данных.
func (s *Service) DeleteProduct(ctx context.Context, id int64) (changeID int64, err error) {
	logger := log.With(s.log, "method", "DeleteProduct")
	if changeID, err = s.repo.AddQueryToDeleteProduct(ctx, id); err != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}
