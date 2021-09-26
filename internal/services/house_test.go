package services

import (
	"errors"
	"testing"

	"github.com/flaky-api/internal/models"
	"github.com/stretchr/testify/assert"
)

type homeVisionClientMock struct {
	err error
	res *models.HousesResponse
}

func (c *homeVisionClientMock) FetchHouses(_ ...models.HouseParamsOption) (houseRes *models.HousesResponse, err error) {
	return c.res, c.err
}

func TestHouseService_FetchHouses_ReturnsUnexpectedError(t *testing.T) {
	// Arrange
	errUnexpected := errors.New("unexpected error")
	svc := houseService{homeVision: &homeVisionClientMock{err: errUnexpected}}

	// Act
	houses, err := svc.FetchHouses(1)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, errUnexpected, err)
	assert.Nil(t, houses)
}

func TestHouseService_FetchHouses_ReturnsStatusUnexpectedError(t *testing.T) {
	// Arrange
	res := &models.HousesResponse{Status: false}
	svc := houseService{homeVision: &homeVisionClientMock{res: res}}

	// Act
	houses, err := svc.FetchHouses(1)

	// Assert
	assert.Error(t, err)
	assert.True(t, errors.Is(err, ErrStatusNotExpected))
	assert.Nil(t, houses)
}

func TestHouseService_FetchHouses_ReturnsSuccess(t *testing.T) {
	// Arrange
	house := models.House{
		ID:       1,
		Address:  "Some Address",
		Owner:    "Owner1",
		Price:    1500,
		PhotoURL: "https://some-image.jpg",
	}
	res := &models.HousesResponse{
		Houses: []models.House{house},
		Status: true,
	}
	svc := houseService{homeVision: &homeVisionClientMock{res: res}}

	// Act
	houses, err := svc.FetchHouses(1)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, houses)
	assert.Equal(t, len(res.Houses), len(houses))
}
