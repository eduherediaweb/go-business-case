package repository

import (
	"errors"
	"fmt"
	"github.com/mytheresa/go-hiring-challenge/models"
	"gorm.io/gorm"
)

type ProductsRepository struct {
	db *gorm.DB
}

func NewProductsDatabaseRepository(db *gorm.DB) models.ProductsRepositoryInterface {
	return &ProductsRepository{
		db: db,
	}
}

func (r *ProductsRepository) GetProductsByCriteria(filter models.Filter, pagination models.Pagination) ([]models.Product, int32, error) {
	var products []models.Product
	var total int64

	query := r.db.Model(&models.Product{}).Preload("Category").Preload("Variants")

	if filter.Category != "" {
		query = query.Joins("JOIN categories ON products.category_id = categories.id").
			Where("categories.code = ?", filter.Category)
	}

	if filter.PriceLessThan > 0 {
		query = query.Where("price < ?", filter.PriceLessThan)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Offset(pagination.Offset).Limit(pagination.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int32(total), nil
}

func (r *ProductsRepository) GetProductsByIdWithVariants(productID uint64) (models.Product, error) {
	var product models.Product

	if err := r.db.Preload("Category").Preload("Variants").
		First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return product, fmt.Errorf("product with ID %d not found", productID)
		}
		return product, err
	}

	return product, nil
}
