package v1

import (
	"context"
	"time"
)

type ProductDomain struct {
	Id          int
	Name        string
	Description string
	Price       float64
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type ProductUsecase interface {
	GetAllProducts(ctx context.Context) (domains []ProductDomain, statusCode int, err error)
	StoreProduct(ctx context.Context, product *ProductDomain) (domain ProductDomain, statusCode int, err error)
	GetProductById(ctx context.Context, id int) (domain ProductDomain, statusCode int, err error)
	UpdateProduct(ctx context.Context, product *ProductDomain, id int) (domain ProductDomain, statusCode int, err error)
	DeleteProduct(ctx context.Context, id int) (statusCode int, err error)
}

type ProductRepository interface {
	GetAllProducts(ctx context.Context) ([]ProductDomain, error)
	StoreProduct(ctx context.Context, product *ProductDomain) (ProductDomain, error)
	GetProductById(ctx context.Context, id int) (ProductDomain, error)
	UpdateProduct(ctx context.Context, product *ProductDomain) (err error)
	DeleteProduct(ctx context.Context, id int) error
}
