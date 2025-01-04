package tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestSearchPacket_WithQuery_Success() {
	searchString := "Test"
	limit := int64(10)
	offset := int64(0)
	expectedPackages := []models.Package{
		createTestPackage(1, "Test Package 1"),
	}

	suite.mockRepo.On("SearchPacket", mock.Anything, searchString, limit, offset).
		Return(expectedPackages, nil).
		Once()

	packages, err := suite.svc.SearchPacket(context.Background(), searchString, limit, offset)

	assert.NoError(suite.T(), err, "Expected no error when searching packets with query")
	assert.Equal(suite.T(), expectedPackages, packages, "Expected packages to match the mocked packages")
}

func (suite *ServiceTestSuite) TestSearchPacket_EmptyQuery_Success() {
	searchString := ""
	limit := int64(5)
	offset := int64(0)
	expectedPackages := []models.Package{
		createTestPackage(1, "Package 1"),
		createTestPackage(2, "Package 2"),
	}

	suite.mockRepo.On("GetAllPackages", mock.Anything, limit, offset).
		Return(expectedPackages, nil).
		Once()

	packages, err := suite.svc.SearchPacket(context.Background(), searchString, limit, offset)

	assert.NoError(suite.T(), err, "Expected no error when searching packets with empty query")
	assert.Equal(suite.T(), expectedPackages, packages, "Expected packages to match the mocked packages")
}

func (suite *ServiceTestSuite) TestSearchPacket_WithQuery_RepositoryError() {
	searchString := "NonExistent"
	limit := int64(10)
	offset := int64(0)
	expectedError := errors.New("database error")

	suite.mockRepo.On("SearchPacket", mock.Anything, searchString, limit, offset).
		Return(nil, expectedError).
		Once()

	packages, err := suite.svc.SearchPacket(context.Background(), searchString, limit, offset)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during search")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), packages, "Expected packages to be nil on error")
}

func (suite *ServiceTestSuite) TestSearchPacket_EmptyQuery_RepositoryError() {
	searchString := ""
	limit := int64(5)
	offset := int64(0)
	expectedError := errors.New("database error")

	suite.mockRepo.On("GetAllPackages", mock.Anything, limit, offset).
		Return(nil, expectedError).
		Once()

	packages, err := suite.svc.SearchPacket(context.Background(), searchString, limit, offset)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during GetAllPackages")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Nil(suite.T(), packages, "Expected packages to be nil on error")
}
