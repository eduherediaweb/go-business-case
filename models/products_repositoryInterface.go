package models

type ProductsRepositoryInterface interface {
	GetAllProducts() ([]Product, error)
}
