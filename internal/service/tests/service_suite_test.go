package tests

import (
	"testing"

	"github.com/Chaika-Team/ChaikaGoods/internal/repository/mocks"
	"github.com/Chaika-Team/ChaikaGoods/internal/service"

	"github.com/go-kit/log"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite
	mockRepo *mocks.MockGoodsRepository
	svc      service.Service // Ваш сервисный слой
}

func (suite *ServiceTestSuite) SetupTest() {
	// Инициализация мока
	suite.mockRepo = mocks.NewMockGoodsRepository(suite.T())
	// Инициализация логгера (можно использовать заглушку или реальный)
	logger := log.NewNopLogger()
	// Инициализация сервисного слоя с мок-репозиторием
	suite.svc = service.NewService(suite.mockRepo, logger)
}

func (suite *ServiceTestSuite) TearDownTest() {
	// Проверка, что все ожидания моков были выполнены
	suite.mockRepo.AssertExpectations(suite.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}
