package requests

import (
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type ProductRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float32 `json:"price" binding:"required,gt=0"`  // price lebih besar dari 0
	Stock       int     `json:"stock" binding:"required,gte=0"` // stock tidak negatif

}

func (productRequest *ProductRequest) ToDomain() *V1Domains.ProductDomain {
	return &V1Domains.ProductDomain{
		Name:        productRequest.Name,
		Description: productRequest.Description,
		Price:       productRequest.Price,
		Stock:       productRequest.Stock,
	}
}

type ProductUpdateRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description" binding:"required"`
	Price       float32 `json:"price" binding:"required,gt=0"`  // price lebih besar dari 0
	Stock       int     `json:"stock" binding:"required,gte=0"` // stock tidak negatif
}

func (p *ProductUpdateRequest) ToDomain() *V1Domains.ProductDomain {
	return &V1Domains.ProductDomain{
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
	}
}
