package models

type ProductsRepositoryInterface interface {
	GetAllProducts() ([]Product, error)
	GetAllProductsPaginated(pagination Pagination) ([]Product, int32, error)
}
