package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/Chaika-Team/ChaikaGoods/internal/myerr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestAddPackage_Success() {
	newPackage := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedID := int64(1)

	suite.mockRepo.On("CreatePackage", mock.Anything, mock.MatchedBy(func(p *models.Package) bool {
		return p.PackageName == newPackage.PackageName &&
			p.Description == newPackage.Description &&
			len(p.Content) == len(newPackage.Content)
	})).
		Run(func(args mock.Arguments) {
			Package := args.Get(1).(*models.Package)
			Package.ID = expectedID
			// Дополнительные проверки
			assert.NotEmpty(suite.T(), Package.PackageName, "Package name should not be empty")
			assert.NotEmpty(suite.T(), Package.Description, "Package description should not be empty")
			assert.Len(suite.T(), Package.Content, 2, "Package should have two contents")
		}).
		Return(nil).
		Once()

	createdID, err := suite.svc.AddPackage(context.Background(), newPackage)

	assert.NoError(suite.T(), err, "Expected no error when adding Package")
	assert.Equal(suite.T(), expectedID, createdID, "Expected created ID to match the mocked ID")
	assert.Equal(suite.T(), expectedID, newPackage.ID, "Expected Package ID to be set to the mocked ID")
}

func (suite *ServiceTestSuite) TestAddPackage_AlreadyExists() {
	newPackage := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.Conflict("Package with name New Package already exists", nil)

	suite.mockRepo.On("CreatePackage", mock.Anything, mock.MatchedBy(func(p *models.Package) bool {
		return p.PackageName == newPackage.PackageName &&
			p.Description == newPackage.Description &&
			len(p.Content) == len(newPackage.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddPackage(context.Background(), newPackage)

	assert.Error(suite.T(), err, "Expected error when Package already exists")
	assert.True(suite.T(), myerr.IsConflict(err), "Expected error to be of type Conflict")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}

func (suite *ServiceTestSuite) TestAddPackage_ProductNotFound() {
	newPackage := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.NotFound("Product with ID 1 not found", nil)

	suite.mockRepo.On("CreatePackage", mock.Anything, mock.MatchedBy(func(p *models.Package) bool {
		return p.PackageName == newPackage.PackageName &&
			p.Description == newPackage.Description &&
			len(p.Content) == len(newPackage.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddPackage(context.Background(), newPackage)

	assert.Error(suite.T(), err, "Expected error when Product is not found")
	assert.True(suite.T(), myerr.IsNotFound(err), "Expected error to be of type NotFound")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}

func (suite *ServiceTestSuite) TestAddPackage_RepositoryError() {
	newPackage := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := myerr.Internal("Database error", errors.New("connection failed"))

	suite.mockRepo.On("CreatePackage", mock.Anything, mock.MatchedBy(func(p *models.Package) bool {
		return p.PackageName == newPackage.PackageName &&
			p.Description == newPackage.Description &&
			len(p.Content) == len(newPackage.Content)
	})).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddPackage(context.Background(), newPackage)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during AddPackage")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}
