package service

import (
	"ChaikaGoods/internal/models"
	"ChaikaGoods/internal/myerr"
	repo "ChaikaGoods/internal/repository"
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
	// GetProductByID возвращает продукт по его ID.
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	// GetCurrentVersion возвращает текущую версию базы данных продуктов.
	GetCurrentVersion(ctx context.Context) (models.Version, error)
	// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
	GetDelta(ctx context.Context, version int) ([]models.Change, error)
	// SearchPacket ищет пакеты продуктов по их имени или ID.
	SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error)
	// AddPacket добавляет новый пакет продуктов в базу данных.
	AddPacket(ctx context.Context, packet *models.Package, packageContent []models.PackageContent) (int64, error)
	// AddProduct добавляет новый продукт в базу данных (добавляет запрос).
	AddProduct(ctx context.Context, p *map[string]interface{}) (changeID int64, err error)
	// UpdateProduct обновляет информацию о продукте в базе данных (добавляет запрос).
	UpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error)
	// DeleteProduct удаляет продукт из базы данных (добавляет запрос).
	DeleteProduct(ctx context.Context, id int64) (changeID int64, err error)
}
type Service struct {
	repo repo.GoodsRepository
	log  log.Logger
}

// NewService создает новый экземпляр GoodsService.
func NewService(repo repo.GoodsRepository, logger log.Logger) GoodsService {
	return &Service{
		repo: repo,
		log:  logger,
	}
}

// GetAllProducts возвращает список всех продуктов.
func (s *Service) GetAllProducts(ctx context.Context) ([]models.Product, error) {
	logger := log.With(s.log, "method", "GetAllProducts")
	products, err := s.repo.GetAllProducts(ctx)

	if myerr.ToAppError(logger, err, "Error to get all products") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	return products, nil
}

func (s *Service) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	logger := log.With(s.log, "method", "GetProductByID")
	product, err := s.repo.GetProductByID(ctx, id)
	if myerr.ToAppError(logger, err, "Error to get product by ID") != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Product{}, err
	}
	return product, nil
}

// GetCurrentVersion возвращает текущую версию базы данных продуктов.
func (s *Service) GetCurrentVersion(ctx context.Context) (models.Version, error) {
	logger := log.With(s.log, "method", "GetCurrentVersion")
	version, err := s.repo.GetCurrentActualVersion(ctx)
	if myerr.ToAppError(logger, err, "Error to get current version") != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Version{}, err
	}
	return version, nil
}

// GetDelta возвращает изменения в базе данных продуктов по сравнению с версией на устройстве
func (s *Service) GetDelta(ctx context.Context, version int) ([]models.Change, error) {
	logger := log.With(s.log, "method", "GetDelta")
	// get actual version
	actualVersion, err := s.repo.GetCurrentActualVersion(ctx)
	if myerr.ToAppError(logger, err, "Error to get current version") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	// get all versions between the device version and the actual version
	versions, err := s.repo.GetVersionsBetween(ctx, version, actualVersion.VersionID)
	if myerr.ToAppError(logger, err, "Error to get versions between") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	// get all changes between the device version and the actual version
	changes := make([]models.Change, 0)
	for _, v := range versions {
		c, err := s.repo.GetAllChanges(ctx, v)
		if myerr.ToAppError(logger, err, "Error to get all changes") != nil {
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
	packages, err := s.repo.SearchPacket(ctx, searchString, quantity, offset)
	if myerr.ToAppError(logger, err, "Error to search packet") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	return packages, nil

}

// AddPacket добавляет новый пакет продуктов в базу данных.
func (s *Service) AddPacket(ctx context.Context, packet *models.Package, packageContent []models.PackageContent) (int64, error) {
	logger := log.With(s.log, "method", "AddPacket")
	err := s.repo.CreatePackage(ctx, packet)
	if myerr.ToAppError(logger, err, "Error to create package") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	// add products to the package
	err = s.repo.AddProductToPackage(ctx, packet.ID, packageContent)
	if myerr.ToAppError(logger, err, "Error to add product to package") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return packet.ID, nil
}

// AddProduct добавляет новый продукт в базу данных.
func (s *Service) AddProduct(ctx context.Context, p *map[string]interface{}) (changeID int64, err error) {
	logger := log.With(s.log, "method", "AddProduct")
	changeID, err = s.repo.AddQueryToCreateProduct(ctx, p)
	if myerr.ToAppError(logger, err, "Error to add product") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (s *Service) UpdateProduct(ctx context.Context, data *map[string]interface{}) (changeID int64, err error) {
	logger := log.With(s.log, "method", "UpdateProduct")
	changeID, err = s.repo.AddQueryToUpdateProduct(ctx, data)
	if myerr.ToAppError(logger, err, "Error to update product") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}

// DeleteProduct удаляет продукт из базы данных.
func (s *Service) DeleteProduct(ctx context.Context, id int64) (changeID int64, err error) {
	logger := log.With(s.log, "method", "DeleteProduct")
	changeID, err = s.repo.AddQueryToDeleteProduct(ctx, id)
	if myerr.ToAppError(logger, err, "Error to delete product") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return changeID, nil
}
