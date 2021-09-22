package transports

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/flaky-api/internal/config"

	"github.com/flaky-api/internal/models"
	"github.com/go-resty/resty/v2"
)

var (
	ErrMaxRetriesReached = errors.New("max retries reached")
)

type HomeVisionClient interface {
	FetchHouses(opts ...models.HouseParamsOption) (*models.HousesResponse, error)
}

type homeVisionClient struct {
	client *resty.Client
}

func NewHomeVisionClient(retries int, retryWaitTime time.Duration) HomeVisionClient {
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
	}
}

func (c *homeVisionClient) FetchHouses(opts ...models.HouseParamsOption) (houseRes *models.HousesResponse, err error) {
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
		return nil, err
	}

	if res.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status code %d not expected: %w", res.StatusCode(), ErrMaxRetriesReached)
	}

	if houseRes == nil {
		return nil, fmt.Errorf("house response is empty")
	}

	return
}
