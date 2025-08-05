package variants

import (
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Mock repository
type MockVariantRepository struct {
	mock.Mock
}

func (m *MockVariantRepository) GetProductsByCriteria(filter models.Filter, pagination models.Pagination) ([]models.Product, int32, error) {
	args := m.Called(filter, pagination)
	return args.Get(0).([]models.Product), args.Get(1).(int32), args.Error(2)
}

func (m *MockVariantRepository) GetProductsByIdWithVariants(productID uint64) (models.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(models.Product), args.Error(1)
}

func TestVariantsHandler_HandleGet_Success(t *testing.T) {
	mockRepo := new(MockVariantRepository)
	handler := NewVariantsHandler(mockRepo)

	product := &models.Product{
		ID:    1,
		Code:  "PROD001",
		Price: decimal.NewFromFloat(99.99),
		Category: models.Category{
			Code: "clothing",
			Name: "Clothing",
		},
		Variants: []models.Variant{
			{
				ID:    1,
				Name:  "Size M",
				SKU:   "PROD001-M",
				Price: decimal.NewFromFloat(89.99),
			},
			{
				ID:    2,
				Name:  "Size L",
				SKU:   "PROD001-L",
				Price: decimal.Zero, // Should inherit product price
			},
		},
	}

	mockRepo.On("GetProductsByIdWithVariants", uint64(1)).Return(*product, nil)

	// Create request
	req := httptest.NewRequest("GET", "/catalog/1", nil)
	recorder := httptest.NewRecorder()

	// Execute
	handler.HandleGet(recorder, req)

	// Assertions
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "PROD001")
	assert.Contains(t, recorder.Body.String(), "89.99") // Variant specific price
	assert.Contains(t, recorder.Body.String(), "99.99") // Inherited price

	mockRepo.AssertExpectations(t)
}

func TestVariantsHandler_HandleGet_PriceInheritance(t *testing.T) {
	mockRepo := new(MockVariantRepository)
	handler := NewVariantsHandler(mockRepo)

	product := &models.Product{
		ID:    1,
		Code:  "PROD001",
		Price: decimal.NewFromFloat(50.00),
		Variants: []models.Variant{
			{
				ID:    1,
				Name:  "Variant without price",
				SKU:   "SKU001",
				Price: decimal.Zero, // Should inherit 50.00
			},
		},
	}

	mockRepo.On("GetProductsByIdWithVariants", uint64(1)).Return(*product, nil)

	req := httptest.NewRequest("GET", "/catalog/1", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "50")
}

func TestVariantsHandler_HandleGet_InvalidID(t *testing.T) {
	mockRepo := new(MockVariantRepository)
	handler := NewVariantsHandler(mockRepo)

	req := httptest.NewRequest("GET", "/catalog/invalid", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Invalid product ID format")
}

func TestVariantsHandler_HandleGet_ProductNotFound(t *testing.T) {
	mockRepo := new(MockVariantRepository)
	handler := NewVariantsHandler(mockRepo)

	mockRepo.On("GetProductsByIdWithVariants", uint64(999)).Return(models.Product{}, gorm.ErrRecordNotFound)

	req := httptest.NewRequest("GET", "/catalog/999", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusNotFound, recorder.Code)

	mockRepo.AssertExpectations(t)
}
