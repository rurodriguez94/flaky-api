package transports

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/sirupsen/logrus"

	"github.com/flaky-api/internal/config"
	"github.com/flaky-api/internal/models"
)

var (
	ErrMaxRetriesReached = errors.New("max retries reached")
)

type HomeVisionClient interface {
	FetchHouses(opts ...models.HouseParamsOption) (*models.HousesResponse, error)
}

type homeVisionClient struct {
	client *resty.Client
	logger *logrus.Logger
}

func NewHomeVisionClient(logger *logrus.Logger, retries int, retryWaitTime time.Duration) HomeVisionClient {
	client := resty.New().
		SetHostURL(config.HomeVisionHost).
		SetRetryCount(retries).
		SetRetryWaitTime(retryWaitTime)

	retryCondition := func(res *resty.Response, err error) bool {
		return res.StatusCode() != http.StatusOK
	}

	client.AddRetryCondition(retryCondition)

	return &homeVisionClient{
		client: client,
		logger: logger,
	}
}

func (c *homeVisionClient) FetchHouses(opts ...models.HouseParamsOption) (houseRes *models.HousesResponse, err error) {
	logger := c.logger.WithField("request", "fetch_houses")
	req := c.client.R()

	var params models.HouseParams
	for _, opt := range opts {
		opt(&params)
	}

	if params.Page != "" {
		req.SetQueryParam("page", params.Page)
	}

	res, err := req.SetResult(&houseRes).Get(config.HousePath)
	if err != nil {
		logger.WithError(err).Error("Failed executing request")
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		err := fmt.Errorf("status code %d not expected: %w", res.StatusCode(), ErrMaxRetriesReached)
		logger.WithError(err).Error(fmt.Sprintf("Failed fetching page: %s", params.Page))
		return nil, err
	}

	if houseRes == nil {
		err := fmt.Errorf("house response is empty")
		logger.WithError(err).Error()
		return nil, err
	}

	logger.Info(fmt.Sprintf("Fetching page %s successful", params.Page))

	return
}
