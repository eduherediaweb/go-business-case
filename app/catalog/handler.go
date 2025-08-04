package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products []Product `json:"products"`
}

type Product struct {
	Code  string  `json:"code"`
	Price float64 `json:"price"`
}

type CatalogHandler struct {
	repo models.ProductsRepositoryInterface
}

func NewCatalogHandler(r *models.ProductsRepositoryInterface) *CatalogHandler {
	return &CatalogHandler{
		repo: *r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
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
		}
	}

	response := Response{
		Products: products,
	}

	// Return the products as a JSON response
	api.OKResponse(w, response)
}
