package models

type HousesResponse struct {
	Houses []House `json:"houses"`
	Status bool    `json:"ok"`
}
