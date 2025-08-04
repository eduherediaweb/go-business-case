package catalog

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"net/http"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Products   []Product         `json:"products"`
	Pagination models.Pagination `json:"pagination" validate:"required"`
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

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {

	paginationParams := models.PaginationParams(r.URL.Query())

	res, total, err := h.repo.GetAllProductsPaginated(paginationParams)
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
		Pagination: models.Pagination{
			Offset: paginationParams.Offset,
			Limit:  paginationParams.Limit,
			Total:  total,
		},
	}

	// Return the products as a JSON response
	api.OKResponse(w, response)
}
