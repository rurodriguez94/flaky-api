package services

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"sync"

	"github.com/flaky-api/internal/models"
	"github.com/flaky-api/internal/transports"
)

var (
	ErrStatusNotExpected = errors.New("status not expected")
)

type HouseService interface {
	FetchHouses(pages int) ([]models.House, error)
	DownloadHouseImages(houses ...models.House) error
}

type houseService struct {
	homeVision transports.HomeVisionClient
	downloader transports.DownloaderClient
	filepath   string
}

func NewHouseService(homeVision transports.HomeVisionClient, downloader transports.DownloaderClient) HouseService {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	filepath := fmt.Sprintf("%s/images", dir)

	return &houseService{
		homeVision: homeVision,
		downloader: downloader,
		filepath:   filepath,
	}
}

func (s *houseService) FetchHouses(pages int) ([]models.House, error) {
	var houses []models.House
	ch := make(chan error, pages)
	wg := sync.WaitGroup{}

	for i := 1; i <= pages; i++ {
		wg.Add(1)

		page := i
		go func() {
			defer wg.Done()
			s.fetch(page, &houses, ch)
		}()
	}
	wg.Wait()
	close(ch)

	if err := <-ch; err != nil {
		return nil, err
	}

	// Sort slice by ID ascending
	sort.Slice(houses, func(i, j int) bool {
		return houses[i].ID < houses[j].ID
	})

	return houses, nil
}

func (s *houseService) DownloadHouseImages(houses ...models.House) error {
	ch := make(chan error, len(houses))
	wg := sync.WaitGroup{}

	for _, h := range houses {
		wg.Add(1)

		house := h

		go func() {
			defer wg.Done()
			img, err := s.downloader.DownloadImage(house.PhotoURL)
			if err != nil {
				ch <- err
				return
			}

			file, err := os.Create(fmt.Sprintf("%s/%s", s.filepath, house.Filename()))
			if err != nil {
				ch <- err
				return
			}
			defer file.Close()

			_, err = io.Copy(file, bytes.NewReader(img))
			if err != nil {
				ch <- err
				return
			}
		}()
	}
	wg.Wait()
	close(ch)

	if err := <-ch; err != nil {
		return err
	}

	return nil
}

func (s *houseService) fetch(page int, houses *[]models.House, ch chan<- error) {
	res, err := s.homeVision.FetchHouses(models.WithPage(page))
	if err != nil {
		ch <- err
		return
	}

	if !res.Status {
		ch <- fmt.Errorf("fetch house: %w", ErrStatusNotExpected)
		return
	}

	m := sync.Mutex{}

	m.Lock()
	*houses = append(*houses, res.Houses...)
	m.Unlock()
}
