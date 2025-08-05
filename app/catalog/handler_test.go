package catalog

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) GetProductsByCriteria(filters models.Filter, pagination models.Pagination) ([]models.Product, int32, error) {
	args := m.Called(filters, pagination)
	return args.Get(0).([]models.Product), args.Get(1).(int32), args.Error(2)
}

func (m *MockProductRepository) GetProductsByIdWithVariants(productID uint64) (models.Product, error) {
	args := m.Called(productID)
	return args.Get(0).(models.Product), args.Error(1)
}

func CreateTestProduct(id uint, code string, price float64, categoryCode, categoryName string) models.Product {
	return models.Product{
		ID:    id,
		Code:  code,
		Price: decimal.NewFromFloat(price),
		Category: models.Category{
			Code: categoryCode,
			Name: categoryName,
		},
	}
}

func CreateTestProducts() []models.Product {
	return []models.Product{
		CreateTestProduct(1, "PROD001", 99.99, "clothing", "Clothing"),
		CreateTestProduct(2, "PROD002", 149.99, "shoes", "Shoes"),
		CreateTestProduct(3, "PROD003", 29.99, "accessories", "Accessories"),
	}
}

func TestCatalogHandler_HandleGet_BasicFunctionality(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	testProducts := CreateTestProducts()
	mockRepo.On("GetProductsByCriteria", mock.Anything, mock.Anything).
		Return(testProducts, int32(3), nil)

	// Request
	req := httptest.NewRequest("GET", "/catalog", nil)
	recorder := httptest.NewRecorder()

	// Execute
	handler.HandleGet(recorder, req)

	// Assertions
	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))

	var response Response
	err := json.Unmarshal(recorder.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Len(t, response.Products, 3)
	assert.Equal(t, "PROD001", response.Products[0].Code)
	assert.Equal(t, testProducts[0].Price.InexactFloat64(), response.Products[0].Price)
	assert.Equal(t, "clothing", response.Products[0].Category.Code)
	assert.Equal(t, "Clothing", response.Products[0].Category.Name)

	assert.Equal(t, 0, response.Pagination.Offset)
	assert.Equal(t, 10, response.Pagination.Limit)
	assert.Equal(t, int32(3), response.Pagination.Total)

	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_WithPagination(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	testProducts := []models.Product{
		CreateTestProduct(3, "PROD003", 29.99, "accessories", "Accessories"),
	}

	mockRepo.On("GetProductsByCriteria", mock.Anything, mock.Anything).
		Return(testProducts, int32(10), nil)

	req := httptest.NewRequest("GET", "/catalog?offset=2&limit=1", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response Response
	json.Unmarshal(recorder.Body.Bytes(), &response)

	assert.Len(t, response.Products, 1)
	assert.Equal(t, 2, response.Pagination.Offset)
	assert.Equal(t, 1, response.Pagination.Limit)
	assert.Equal(t, int32(10), response.Pagination.Total)

	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_WithFilters(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	testProducts := []models.Product{
		CreateTestProduct(1, "PROD001", 45.99, "clothing", "Clothing"),
	}

	filter := models.Filter{Category: "clothing", PriceLessThan: 50}
	pagination := models.Pagination{Offset: 0, Limit: 10}
	mockRepo.On("GetProductsByCriteria", filter, pagination).
		Return(testProducts, int32(1), nil)

	req := httptest.NewRequest("GET", "/catalog?category=clothing&price_less_than=50", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response Response
	json.Unmarshal(recorder.Body.Bytes(), &response)

	assert.Len(t, response.Products, 1)
	assert.Equal(t, "clothing", response.Products[0].Category.Code)
	assert.Equal(t, int32(1), response.Pagination.Total)

	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_EmptyResults(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	filter := models.Filter{}
	pagination := models.Pagination{Offset: 0, Limit: 10}
	mockRepo.On("GetProductsByCriteria", filter, pagination).
		Return([]models.Product{}, int32(0), nil)

	req := httptest.NewRequest("GET", "/catalog", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response Response
	json.Unmarshal(recorder.Body.Bytes(), &response)

	assert.Len(t, response.Products, 0)
	assert.Equal(t, int32(0), response.Pagination.Total)

	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_RepositoryError(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	filter := models.Filter{}
	pagination := models.Pagination{Offset: 0, Limit: 10}
	mockRepo.On("GetProductsByCriteria", filter, pagination).
		Return([]models.Product{}, int32(0), errors.New("database error"))

	req := httptest.NewRequest("GET", "/catalog", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusInternalServerError, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "Failed to fetch products")
	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_InvalidPaginationParams(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	testProducts := CreateTestProducts()
	filter := models.Filter{}
	pagination := models.Pagination{Offset: 0, Limit: 100}
	mockRepo.On("GetProductsByCriteria", filter, pagination).
		Return(testProducts, int32(3), nil)

	req := httptest.NewRequest("GET", "/catalog?offset=-5&limit=200", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	var response Response
	json.Unmarshal(recorder.Body.Bytes(), &response)

	assert.Equal(t, 0, response.Pagination.Offset)
	assert.Equal(t, 100, response.Pagination.Limit)

	mockRepo.AssertExpectations(t)
}

func TestCatalogHandler_HandleGet_PaginationEdgeCases(t *testing.T) {
	testCases := []struct {
		name           string
		queryParams    string
		expectedOffset int
		expectedLimit  int
	}{
		{
			name:           "No params - use defaults",
			queryParams:    "",
			expectedOffset: 0,
			expectedLimit:  10,
		},
		{
			name:           "Limit too small",
			queryParams:    "limit=0",
			expectedOffset: 0,
			expectedLimit:  1, // minimum
		},
		{
			name:           "Limit too large",
			queryParams:    "limit=500",
			expectedOffset: 0,
			expectedLimit:  100, // maximum
		},
		{
			name:           "Invalid string params",
			queryParams:    "offset=abc&limit=xyz",
			expectedOffset: 0,  // default
			expectedLimit:  10, // default
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockRepo := new(MockProductRepository)
			handler := NewCatalogHandler(mockRepo)

			filter := models.Filter{}
			pagination := models.Pagination{
				Offset: tc.expectedOffset,
				Limit:  tc.expectedLimit,
			}

			mockRepo.On("GetProductsByCriteria", filter, pagination).
				Return([]models.Product{}, int32(0), nil)

			req := httptest.NewRequest("GET", "/catalog?"+tc.queryParams, nil)
			recorder := httptest.NewRecorder()

			handler.HandleGet(recorder, req)

			assert.Equal(t, http.StatusOK, recorder.Code)

			var response Response
			json.Unmarshal(recorder.Body.Bytes(), &response)

			assert.Equal(t, tc.expectedOffset, response.Pagination.Offset)
			assert.Equal(t, tc.expectedLimit, response.Pagination.Limit)

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestCatalogHandler_HandleGet_FilterEdgeCases(t *testing.T) {
	mockRepo := new(MockProductRepository)
	handler := NewCatalogHandler(mockRepo)

	filter := models.Filter{Category: "invalid", PriceLessThan: 0}
	pagination := models.Pagination{Offset: 0, Limit: 10}
	mockRepo.On("GetProductsByCriteria", filter, pagination).
		Return([]models.Product{}, int32(0), nil)

	req := httptest.NewRequest("GET", "/catalog?category=invalid&price_less_than=-10", nil)
	recorder := httptest.NewRecorder()

	handler.HandleGet(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)

	mockRepo.AssertExpectations(t)
}
