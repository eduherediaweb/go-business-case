package models

import (
	"net/url"
	"strconv"
)

const (
	LimitMin      = 1
	LimitMax      = 100
	OffsetDefault = 0
	LimitDefault  = 10
)

type Pagination struct {
	Offset int   `json:"offset"`
	Limit  int   `json:"limit"`
	Total  int32 `json:"total"`
}

func PaginationParams(values url.Values) Pagination {
	params := Pagination{
		Offset: OffsetDefault,
		Limit:  LimitDefault,
	}

	if offsetStr := values.Get("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}

	if limitStr := values.Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			if limit < LimitMin {
				params.Limit = LimitMin
			} else if limit > LimitMax {
				params.Limit = LimitMax
			} else {
				params.Limit = limit
			}
		}
	}

	return params
}
