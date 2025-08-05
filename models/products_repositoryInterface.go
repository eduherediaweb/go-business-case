package models

type ProductsRepositoryInterface interface {
	GetProductsByCriteria(filter Filter, pagination Pagination) ([]Product, int32, error)
	GetProductsByIdWithVariants(productID uint64) (Product, error)
}
