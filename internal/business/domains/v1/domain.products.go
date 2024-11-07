package v1

import (
	"context"
	"time"
)

type ProductDomain struct {
	Id          int
	Name        string
	Description string
	Price       float32
	Stock       int
	CreatedAt   time.Time
	UpdatedAt   *time.Time
}

type ProductUsecase interface {
	GetAll(ctx context.Context) (domains []ProductDomain, statusCode int, err error)
	Store(ctx context.Context, product *ProductDomain) (domain ProductDomain, statusCode int, err error)
	GetById(ctx context.Context, id int) (domain ProductDomain, statusCode int, err error)
	Update(ctx context.Context, product *ProductDomain, id int) (domain ProductDomain, statusCode int, err error)
	Delete(ctx context.Context, id int) (statusCode int, err error)
}

type ProductRepository interface {
	GetAll(ctx context.Context) ([]ProductDomain, error)
	Store(ctx context.Context, product *ProductDomain) (ProductDomain, error)
	GetById(ctx context.Context, id int) (ProductDomain, error)
	Update(ctx context.Context, product *ProductDomain) (err error)
	Delete(ctx context.Context, id int) error
}
