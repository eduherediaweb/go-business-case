package models

import (
	"net/url"
	"testing"
)

func TestPaginationParams_Defaults(t *testing.T) {
	values := url.Values{}
	params := PaginationParams(values)
	if params.Offset != OffsetDefault {
		t.Errorf("expected OffsetDefault, got %d", params.Offset)
	}
	if params.Limit != LimitDefault {
		t.Errorf("expected LimitDefault, got %d", params.Limit)
	}
}

func TestPaginationParams_ValidOffsetAndLimit(t *testing.T) {
	values := url.Values{}
	values.Set("offset", "10")
	values.Set("limit", "20")
	params := PaginationParams(values)
	if params.Offset != 10 {
		t.Errorf("expected offset 10, got %d", params.Offset)
	}
	if params.Limit != 20 {
		t.Errorf("expected limit 20, got %d", params.Limit)
	}
}

func TestPaginationParams_NegativeOffset(t *testing.T) {
	values := url.Values{}
	values.Set("offset", "-5")
	params := PaginationParams(values)
	if params.Offset != OffsetDefault {
		t.Errorf("expected OffsetDefault for negative offset, got %d", params.Offset)
	}
}

func TestPaginationParams_LimitBelowMin(t *testing.T) {
	values := url.Values{}
	values.Set("limit", "0")
	params := PaginationParams(values)
	if params.Limit != LimitMin {
		t.Errorf("expected LimitMin for limit below min, got %d", params.Limit)
	}
}

func TestPaginationParams_LimitAboveMax(t *testing.T) {
	values := url.Values{}
	values.Set("limit", "1000")
	params := PaginationParams(values)
	if params.Limit != LimitMax {
		t.Errorf("expected LimitMax for limit above max, got %d", params.Limit)
	}
}

func TestPaginationParams_InvalidOffsetAndLimit(t *testing.T) {
	values := url.Values{}
	values.Set("offset", "abc")
	values.Set("limit", "xyz")
	params := PaginationParams(values)
	if params.Offset != OffsetDefault {
		t.Errorf("expected OffsetDefault for invalid offset, got %d", params.Offset)
	}
	if params.Limit != LimitDefault {
		t.Errorf("expected LimitDefault for invalid limit, got %d", params.Limit)
	}
}
