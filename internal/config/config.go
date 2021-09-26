package config

import (
	"flag"
	"time"
)

const (
	DefaultMaxRetries    = 10
	DefaultRetryWaitTime = 200 * time.Millisecond
	DefaultPages         = 10
	DefaultStopOnFailure = false
)

const (
	HomeVisionHost = "http://app-homevision-staging.herokuapp.com"
	HousePath      = "/api_project/houses"

	ImagesFolder = "images"
)

type Config struct {
	MaxRetries    int
	RetryWaitTime time.Duration
	Pages         int
	StopOnFailure bool
}

func NewConfig() *Config {
	// Flags used to set configs not mandatory
	pagesFlag := flag.Int("pages", DefaultPages, "to define the number of pages")
	retryFlag := flag.Int("retries", DefaultMaxRetries, "to define the number of http retries")
	stopOnFailFlag := flag.Bool("stopOnFail", DefaultStopOnFailure, "stop the execution when the max retries of some request fail")

	flag.Parse()

	return &Config{
		MaxRetries:    *retryFlag,
		RetryWaitTime: DefaultRetryWaitTime,
		Pages:         *pagesFlag,
		StopOnFailure: *stopOnFailFlag,
	}
}
