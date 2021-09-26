package transports

import (
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/go-resty/resty/v2"
)

type DownloaderClient interface {
	DownloadImage(URL string) ([]byte, error)
}

type downloaderClient struct {
	client *resty.Client
	logger *logrus.Logger
}

func NewDownloaderClient(logger *logrus.Logger) DownloaderClient {
	return &downloaderClient{
		client: resty.New(),
		logger: logger,
	}
}

func (c *downloaderClient) DownloadImage(URL string) ([]byte, error) {
	logger := c.logger.WithField("request", "download_image")

	res, err := c.client.R().Get(URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		err := fmt.Errorf("status code %d not expected", res.StatusCode())
		logger.WithError(err).Error(fmt.Sprintf("Failed downloading image: %s", URL))
		return nil, err
	}

	return res.Body(), nil
}
