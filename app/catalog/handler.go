package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"net/http"
	"net/url"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/models"
)

const (
	LimitMin = 1
	LimitMax = 100
)

type Response struct {
	Products []Product `json:"products"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category,omitempty"`
}

type CatalogHandler struct {
	repo models.ProductsRepositoryInterface
}

func NewCatalogHandler(r *models.ProductsRepositoryInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: *r,
	}
}

type Pagination struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {

	paginationParams := paginationParams(r.URL.Query())

	res, err := h.repo.GetAllProducts()
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch products")
		return
	}

	// Map response
	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	response := Response{
		Products: products,
	}

	// Return the products as a JSON response
	api.OKResponse(w, response)
}

func paginationParams(values url.Values) Pagination {
	params := Pagination{
		Offset: 0,
		Limit:  10,
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
