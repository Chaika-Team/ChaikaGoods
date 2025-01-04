package unit_tests

import (
	"context"
	"errors"

	"github.com/Chaika-Team/ChaikaGoods/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func (suite *ServiceTestSuite) TestGetPacketByID_Success() {
	packetID := int64(1)
	expectedPacket := createTestPackage(packetID, "Test Package")

	suite.mockRepo.On("GetPackageByID", mock.Anything, mock.MatchedBy(func(packet *models.Package) bool {
		return packet.ID == packetID
	})).
		Run(func(args mock.Arguments) {
			p := args.Get(1).(*models.Package)
			*p = expectedPacket
		}).
		Return(nil).
		Once()

	packet, err := suite.svc.GetPacketByID(context.Background(), packetID)

	assert.NoError(suite.T(), err, "Expected no error when getting packet by ID")
	assert.Equal(suite.T(), expectedPacket, packet, "Expected packet to match the mocked packet")
}

func (suite *ServiceTestSuite) TestGetPacketByID_NotFound() {
	packetID := int64(2)
	expectedError := errors.New("packet not found")

	suite.mockRepo.On("GetPackageByID", mock.Anything, mock.MatchedBy(func(packet *models.Package) bool {
		return packet.ID == packetID
	})).
		Return(expectedError).
		Once()

	packet, err := suite.svc.GetPacketByID(context.Background(), packetID)

	assert.Error(suite.T(), err, "Expected error when packet is not found")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), models.Package{}, packet, "Expected packet to be empty")
}

func (suite *ServiceTestSuite) TestGetPacketByID_RepositoryError() {
	packetID := int64(3)
	expectedError := errors.New("database error")

	suite.mockRepo.On("GetPackageByID", mock.Anything, mock.MatchedBy(func(packet *models.Package) bool {
		return packet.ID == packetID
	})).
		Return(expectedError).
		Once()

	packet, err := suite.svc.GetPacketByID(context.Background(), packetID)

	assert.Error(suite.T(), err, "Expected error when repository returns an error during GetPacketByID")
	assert.Equal(suite.T(), expectedError, err, "Expected error to match the mocked error")
	assert.Equal(suite.T(), models.Package{}, packet, "Expected packet to be empty")
}
