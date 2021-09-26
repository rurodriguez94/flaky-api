package app

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/flaky-api/internal/config"
	"github.com/flaky-api/internal/services"
	"github.com/flaky-api/internal/transports"
)

type application struct {
	HouseService services.HouseService
	Config       *config.Config
	Logger       *logrus.Logger
}

func NewApplication(cfg *config.Config, logger *logrus.Logger) *application {
	houseService, err := makeHouseService(logger, cfg.MaxRetries)
	if err != nil {
		logger.WithError(err).Fatal("Failed making house service")
	}

	return &application{
		HouseService: houseService,
		Config:       cfg,
		Logger:       logger,
	}
}

func (app *application) Run() error {
	logger := app.Logger
	cfg := app.Config

	logger.WithFields(logrus.Fields{
		"pages":           app.Config.Pages,
		"max_retries":     app.Config.MaxRetries,
		"stop_on_failure": app.Config.StopOnFailure,
	}).Info("Configuration")

	logger.Info("Fetching houses starting")
	houses, err := app.HouseService.FetchHouses(app.Config.Pages)
	if err != nil {
		logger.WithError(err).Error("Failed fetching houses")
		if cfg.StopOnFailure {
			return err
		}
	}

	logger.WithField("houses", len(houses)).Info("Fetching houses finished")

	logger.Info("Downloading houses starting")
	err = app.HouseService.DownloadHouseImages(houses...)
	if err != nil {
		logger.WithError(err).Error("Failed downloading house images")
		if cfg.StopOnFailure {
			return err
		}
	}

	return nil
}

func makeHouseService(logger *logrus.Logger, maxRetries int) (services.HouseService, error) {
	homeVisionClient := transports.NewHomeVisionClient(logger, maxRetries, config.DefaultRetryWaitTime)
	downloaderClient := transports.NewDownloaderClient(logger)

	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// Create images folder if it does not exist
	imagesPath := filepath.Join(dir, config.ImagesFolder)
	err = os.MkdirAll(imagesPath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return services.NewHouseService(homeVisionClient, downloaderClient, imagesPath), nil
}
