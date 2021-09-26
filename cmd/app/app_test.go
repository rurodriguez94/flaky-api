package app

import (
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/flaky-api/internal/config"
	"github.com/flaky-api/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type houseServiceMock struct {
	err error
	res []models.House
}

func (srv houseServiceMock) FetchHouses(_ int) ([]models.House, error) {
	return srv.res, srv.err
}

func (srv houseServiceMock) DownloadHouseImages(_ ...models.House) error {
	return srv.err
}

func TestApp_MakeHouseService_CreatesImageFolder(t *testing.T) {
	// Arrange
	cfg := config.Config{MaxRetries: config.DefaultMaxRetries, StopOnFailure: config.DefaultStopOnFailure}
	path, err := os.Getwd()
	imagesFolder := fmt.Sprintf("%s/%s", path, config.ImagesFolder)
	require.NoError(t, err)

	// Act
	houseService, err := makeHouseService(nil, cfg.MaxRetries)
	require.NoError(t, err)

	fileInfo, err := os.Stat(imagesFolder)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, houseService)
	assert.NotNil(t, fileInfo)
	assert.True(t, fileInfo.IsDir())
	assert.NoError(t, os.RemoveAll(imagesFolder))
}

func TestApp_Run_ReturnsSuccess(t *testing.T) {
	// Arrange
	app := &application{
		HouseService: houseServiceMock{},
		Config: &config.Config{
			MaxRetries:    config.DefaultMaxRetries,
			StopOnFailure: config.DefaultStopOnFailure,
		},
		Logger: logrus.New(),
	}

	// Act
	err := app.Run()

	// Assert
	assert.NoError(t, err)
}
