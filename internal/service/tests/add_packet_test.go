package tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestAddPacket_Success() {
	newPacket := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedID := int64(1)

	suite.mockRepo.On("CreatePackage", mock.Anything, newPacket).
		Run(func(args mock.Arguments) {
			packet := args.Get(1).(*models.Package)
			packet.ID = expectedID
			// Дополнительные проверки
			assert.NotEmpty(suite.T(), packet.PackageName, "Package name should not be empty")
			assert.NotEmpty(suite.T(), packet.Description, "Package description should not be empty")
			assert.Len(suite.T(), packet.Content, 2, "Package should have two contents")
		}).
		Return(nil).
		Once()

	createdID, err := suite.svc.AddPacket(context.Background(), newPacket)

	assert.NoError(suite.T(), err, "Expected no error when adding packet")
	assert.Equal(suite.T(), expectedID, createdID, "Expected created ID to match the mocked ID")
	assert.Equal(suite.T(), expectedID, newPacket.ID, "Expected packet ID to be set to the mocked ID")
}

func (suite *ServiceTestSuite) TestAddPacket_RepositoryError() {
	newPacket := &models.Package{
		PackageName: "New Package",
		Description: "New Package Description",
		Content: []models.PackageContent{
			{ProductID: 1, Quantity: 2},
			{ProductID: 2, Quantity: 3},
		},
	}
	expectedError := errors.New("database error")

	suite.mockRepo.On("CreatePackage", mock.Anything, newPacket).
		Return(expectedError).
		Once()

	createdID, err := suite.svc.AddPacket(context.Background(), newPacket)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during AddPacket")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), int64(0), createdID, "Expected created ID to be 0 on error")
}
