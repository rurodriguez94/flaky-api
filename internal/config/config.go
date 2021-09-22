package config

import "time"

const (
	DefaultMaxRetries    = 10
	DefaultRetryWaitTime = 200 * time.Millisecond
	DefaultPages         = 10
)

const (
	HomeVisionHost = "http://app-homevision-staging.herokuapp.com"
	HousePath      = "/api_project/houses"
)
