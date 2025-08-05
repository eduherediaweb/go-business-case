package variants

import (
	"github.com/mytheresa/go-hiring-challenge/app/api"
	"net/http"
	"strconv"
	"strings"

	"github.com/mytheresa/go-hiring-challenge/models"
)

type Response struct {
	Product Product `json:"product"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Product struct {
	ID       uint      `json:"id"`
	Code     string    `json:"code"`
	Price    float64   `json:"price"`
	Category Category  `json:"category,omitempty"`
	Variants []Variant `json:"variants,omitempty"`
}

type Variant struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

type VariantHandler struct {
	repo models.ProductsRepositoryInterface
}

func NewVariantsHandler(r models.ProductsRepositoryInterface) *VariantHandler {
	return &VariantHandler{
		repo: r,
	}
}

func (h *VariantHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	idStr := pathParts[2]
	productID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid product ID format")
		return
	}

	product, err := h.repo.GetProductsByIdWithVariants(productID)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Failed to fetch products")
		return
	}

	// Map response
	variants := make([]Variant, len(product.Variants))
	for i, variant := range product.Variants {
		price := variant.Price.InexactFloat64()
		if variant.Price.IsZero() {
			price = product.Price.InexactFloat64()
		}

		variants[i] = Variant{
			ID:    variant.ID,
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: price,
		}
	}

	response := Response{
		Product: Product{
			ID:    product.ID,
			Code:  product.Code,
			Price: product.Price.InexactFloat64(),
			Category: Category{
				Code: product.Category.Code,
				Name: product.Category.Name,
			},
			Variants: variants,
		},
	}

	// Return the products as a JSON response
	api.OKResponse(w, response)
}
