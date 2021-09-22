package transports

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type DownloaderClient interface {
	DownloadImage(URL string) ([]byte, error)
}

type downloaderClient struct {
	client *resty.Client
}

func NewDownloaderClient() DownloaderClient {
	return &downloaderClient{client: resty.New()}
}

func (c *downloaderClient) DownloadImage(URL string) ([]byte, error) {
	res, err := c.client.R().Get(URL)
	if err != nil {
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code %d not expected", res.StatusCode())
	}

	return res.Body(), nil
}
