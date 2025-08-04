package repository

import (
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

func (r *ProductsRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	if err := r.db.Preload("Category").Preload("Variants").Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductsRepository) GetAllProductsPaginated(pagination models.Pagination) ([]models.Product, int32, error) {
	var products []models.Product
	var total int64

	if err := r.db.Model(&models.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := r.db.Preload("Category").Preload("Variants").
		Offset(pagination.Offset).Limit(pagination.Limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, int32(total), nil
}
