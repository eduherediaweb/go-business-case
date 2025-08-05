package models

import (
	"net/url"
	"testing"
)

func TestCreateFilter_CategoryOnly(t *testing.T) {
	values := url.Values{}
	values.Set("category", "shoes")
	filter := CreateFilter(values)
	if filter.Category != "shoes" {
		t.Errorf("expected category 'shoes', got '%s'", filter.Category)
	}
	if filter.PriceLessThan != 0 {
		t.Errorf("expected price_less_than 0, got %f", filter.PriceLessThan)
	}
}

func TestCreateFilter_PriceLessThanOnly(t *testing.T) {
	values := url.Values{}
	values.Set("price_less_than", "99.99")
	filter := CreateFilter(values)
	if filter.Category != "" {
		t.Errorf("expected empty category, got '%s'", filter.Category)
	}
	if filter.PriceLessThan != 99.99 {
		t.Errorf("expected price_less_than 99.99, got %f", filter.PriceLessThan)
	}
}

func TestCreateFilter_CategoryAndPriceLessThan(t *testing.T) {
	values := url.Values{}
	values.Set("category", "bags")
	values.Set("price_less_than", "150")
	filter := CreateFilter(values)
	if filter.Category != "bags" {
		t.Errorf("expected category 'bags', got '%s'", filter.Category)
	}
	if filter.PriceLessThan != 150 {
		t.Errorf("expected price_less_than 150, got %f", filter.PriceLessThan)
	}
}

func TestCreateFilter_InvalidPriceLessThan(t *testing.T) {
	values := url.Values{}
	values.Set("price_less_than", "notanumber")
	filter := CreateFilter(values)
	if filter.PriceLessThan != 0 {
		t.Errorf("expected price_less_than 0 for invalid input, got %f", filter.PriceLessThan)
	}
}

func TestCreateFilter_NegativePriceLessThan(t *testing.T) {
	values := url.Values{}
	values.Set("price_less_than", "-10")
	filter := CreateFilter(values)
	if filter.PriceLessThan != 0 {
		t.Errorf("expected price_less_than 0 for negative input, got %f", filter.PriceLessThan)
	}
}

func TestCreateFilter_EmptyValues(t *testing.T) {
	values := url.Values{}
	filter := CreateFilter(values)
	if filter.Category != "" {
		t.Errorf("expected empty category, got '%s'", filter.Category)
	}
	if filter.PriceLessThan != 0 {
		t.Errorf("expected price_less_than 0, got %f", filter.PriceLessThan)
	}
}
