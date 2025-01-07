package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetPackageByID_Success() {
	PackageID := int64(1)
	expectedPackage := createTestPackage(PackageID, "Test Package")

	suite.mockRepo.On("GetPackageByID", mock.Anything, PackageID).
		Return(expectedPackage, nil).
		Once()

	Package, err := suite.svc.GetPackageByID(context.Background(), PackageID)

	assert.NoError(suite.T(), err, "Expected no error when getting Package by ID")
	assert.Equal(suite.T(), expectedPackage, Package, "Expected Package to match the mocked Package")
}

func (suite *ServiceTestSuite) TestGetPackageByID_NotFound() {
	PackageID := int64(2)
	expectedError := myerr.NotFound("Package not found", nil)

	suite.mockRepo.On("GetPackageByID", mock.Anything, PackageID).
		Return(models.Package{}, expectedError).
		Once()

	Package, err := suite.svc.GetPackageByID(context.Background(), PackageID)

	assert.Error(suite.T(), err, "Expected error when Package is not found")
	assert.True(suite.T(), myerr.IsNotFound(err), "Expected error to be of type NotFound")
	assert.Equal(suite.T(), models.Package{}, Package, "Expected Package to be empty")
}

func (suite *ServiceTestSuite) TestGetPackageByID_RepositoryError() {
	PackageID := int64(3)
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("GetPackageByID", mock.Anything, PackageID).
		Return(models.Package{}, expectedError).
		Once()

	Package, err := suite.svc.GetPackageByID(context.Background(), PackageID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during GetPackageByID")
	assert.True(suite.T(), myerr.IsInternal(err), "Expected error to be of type Internal")
	assert.Equal(suite.T(), models.Package{}, Package, "Expected Package to be empty")
}
