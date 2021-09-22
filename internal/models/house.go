package models

import (
	"fmt"
	"strconv"
)

type House struct {
	ID       uint64  `json:"id"`
	Address  string  `json:"address"`
	Owner    string  `json:"owner"`
	Price    float64 `json:"price"`
	PhotoURL string  `json:"photoURL"`
}

func (h *House) Filename() string {
	return fmt.Sprintf("id-%d-%s.jpg", h.ID, h.Address)
}

type HouseParams struct {
	Page string
}

type HouseParamsOption func(*HouseParams)

func WithPage(page int) HouseParamsOption {
	return func(p *HouseParams) {
		p.Page = strconv.Itoa(page)
	}
}
