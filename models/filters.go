package models

import (
	"net/url"
	"strconv"
)

type Filter struct {
	Category      string  `json:"category,omitempty"`
	PriceLessThan float64 `json:"price_less_than,omitempty"`
}

func CreateFilter(values url.Values) Filter {
	filter := Filter{}
	if category := values.Get("category"); category != "" {
		filter.Category = category
	}

	filter.PriceLessThan = 0
	if priceLessThanParam := values.Get("price_less_than"); priceLessThanParam != "" {
		if PriceLessThan, err := strconv.ParseFloat(priceLessThanParam, 64); err == nil && PriceLessThan >= 0 {
			filter.PriceLessThan = PriceLessThan
		}
	}

	return filter
}
