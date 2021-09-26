package main

import (
	"github.com/sirupsen/logrus"

	"github.com/flaky-api/cmd/app"
	"github.com/flaky-api/internal/config"
)

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	logger.Info("Starting the program")

	cfg := config.NewConfig()
	application := app.NewApplication(cfg, logger)

	if err := application.Run(); err != nil {
		logger.WithError(err).Fatal("Execution stopped")
	}

	logger.Info("Successful execution")
}
