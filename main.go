package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/flaky-api/internal/config"
	"github.com/flaky-api/internal/services"
	"github.com/flaky-api/internal/transports"
)

func main() {
	log.Println("Starting the program")

	// Flags used to set configs not mandatory
	pagesFlag := flag.Int("pages", config.DefaultPages, "to define the number of pages")
	retryFlag := flag.Int("retries", config.DefaultMaxRetries, "to define the number of http retries")

	flag.Parse()

	log.Println("Configuration",
		fmt.Sprintf("Pages: %d", *pagesFlag),
		fmt.Sprintf("MaxRetries: %d", *retryFlag),
	)

	houseService := makeHouseService(*retryFlag)

	log.Println("Fetching houses starting")
	houses, err := houseService.FetchHouses(*pagesFlag)
	if err != nil {
		log.Fatal("Failed fetching houses: ", err)
	}

	log.Println("Downloading houses starting")
	err = houseService.DownloadHouseImages(houses...)
	if err != nil {
		log.Fatal("Failed downloading house images: ", err)
	}
}

func makeHouseService(maxRetries int) services.HouseService {
	homeVisionClient := transports.NewHomeVisionClient(maxRetries, config.DefaultRetryWaitTime)
	downloaderClient := transports.NewDownloaderClient()

	return services.NewHouseService(homeVisionClient, downloaderClient)
}
