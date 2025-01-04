package service

import (
	"context"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	repo "github.com/Chaika-Team/ChaikaGoods/internal/repository"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

/* Бизнес задачи для микросервиса продуктов продаваемых в вагонах поезда:
Со стороны проводника поезда:
1. Проводник запрашивает список продуктов для своей базы данных с центральной базой данных продуктов.
3. Проводник ищет пакеты продуктов по их имени или ID.
4. Проводник добавляет новый пакет продуктов в базу данных.
5. Проводник обновляет информацию о пакете продуктов в базе данных.
Со стороны администратора системы:
1. Администратор может добавить новый продукт в базу данных.
2. Администратор может обновить информацию о продукте в базе данных.
3. Администратор может удалить продукт из базы данных.
4. Все изменения продуктов происходят через CQRS и создание версий базы продуктов.
*/

// Service описывает сервис для работы с продуктами.
type Service interface {
	// GetAllProducts возвращает список всех продуктов.
	GetAllProducts(ctx context.Context) ([]models.Product, error)
	// GetProductByID возвращает продукт по его ID.
	GetProductByID(ctx context.Context, id int64) (models.Product, error)
	// SearchPacket ищет пакеты продуктов по их имени или ID.
	SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error)
	// AddPacket добавляет новый пакет продуктов в базу данных.
	AddPacket(ctx context.Context, packet *models.Package) (int64, error)
	// GetPacketByID возвращает пакет продуктов по его ID.
	GetPacketByID(ctx context.Context, id int64) (models.Package, error)
	// CreateProduct добавляет новый продукт в базу данных.
	CreateProduct(ctx context.Context, p *models.Product) (productId int64, err error)
	// UpdateProduct обновляет информацию о продукте в базе данных.
	UpdateProduct(ctx context.Context, p *models.Product) (err error)
	// DeleteProduct удаляет продукт из базы данных.
	DeleteProduct(ctx context.Context, id int64) (err error)
}
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

	if myerr.ToAppError(logger, err, "Error to get all products") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}
	return products, nil
}

func (s *GoodsService) GetProductByID(ctx context.Context, id int64) (models.Product, error) {
	logger := log.With(s.log, "method", "GetProductByID")
	product, err := s.repo.GetProductByID(ctx, id)
	if myerr.ToAppError(logger, err, "Error to get product by ID") != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Product{}, err
	}
	return product, nil
}

// SearchPacket ищет пакеты продуктов по их имени или ID.
// TODO переделать через поисковые системы типа ElasticSearch
func (s *GoodsService) SearchPacket(ctx context.Context, searchString string, quantity int64, offset int64) ([]models.Package, error) {
	logger := log.With(s.log, "method", "SearchPacket")

	if searchString == "" {
		// Пустая строка поиска, возвращаем все пакеты с пагинацией
		packages, err := s.repo.GetAllPackages(ctx, quantity, offset)
		if myerr.ToAppError(logger, err, "Error to retrieve all packages") != nil {
			_ = level.Error(logger).Log("err", err)
			return nil, err
		}
		return packages, nil
	}

	// Поиск пакетов по строке
	packages, err := s.repo.SearchPacket(ctx, searchString, quantity, offset)
	if myerr.ToAppError(logger, err, "Error to search packet") != nil {
		_ = level.Error(logger).Log("err", err)
		return nil, err
	}

	return packages, nil
}

// AddPacket добавляет новый пакет продуктов в базу данных.
func (s *GoodsService) AddPacket(ctx context.Context, packet *models.Package) (int64, error) {
	logger := log.With(s.log, "method", "AddPacket")
	err := s.repo.CreatePackage(ctx, packet)
	if myerr.ToAppError(logger, err, "Error to create package") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return packet.ID, nil
}

// CreateProduct добавляет новый продукт в базу данных.
func (s *GoodsService) CreateProduct(ctx context.Context, p *models.Product) (productId int64, err error) {
	logger := log.With(s.log, "method", "CreateProduct")
	productId, err = s.repo.CreateProduct(ctx, p)
	if myerr.ToAppError(logger, err, "Error to create product") != nil {
		_ = level.Error(logger).Log("err", err)
		return 0, err
	}
	return productId, nil
}

func (s *GoodsService) GetPacketByID(ctx context.Context, id int64) (models.Package, error) {
	logger := log.With(s.log, "method", "GetPacketByID")
	packet := models.Package{ID: id}
	err := s.repo.GetPackageByID(ctx, &packet)
	if myerr.ToAppError(logger, err, "Error to get packet by ID") != nil {
		_ = level.Error(logger).Log("err", err)
		return models.Package{}, err
	}
	return packet, nil
}

// UpdateProduct обновляет информацию о продукте в базе данных.
func (s *GoodsService) UpdateProduct(ctx context.Context, p *models.Product) (err error) {
	logger := log.With(s.log, "method", "UpdateProduct")
	err = s.repo.UpdateProduct(ctx, p)
	if myerr.ToAppError(logger, err, "Error to update product") != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	return nil
}

// DeleteProduct удаляет продукт из базы данных.
func (s *GoodsService) DeleteProduct(ctx context.Context, id int64) (err error) {
	logger := log.With(s.log, "method", "DeleteProduct")
	err = s.repo.DeleteProduct(ctx, id)
	if myerr.ToAppError(logger, err, "Error to delete product") != nil {
		_ = level.Error(logger).Log("err", err)
		return err
	}
	return nil
}
