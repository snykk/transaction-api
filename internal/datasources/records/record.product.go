package records

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type Product struct {
	Id          int        `db:"product_id"`
	Name        string     `db:"name"`
	Description string     `db:"description"`
	Price       float64    `db:"price"`
	Stock       int        `db:"stock"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   *time.Time `db:"updated_at"`
}

// Mapper
func (p *Product) ToV1Domain() V1Domains.ProductDomain {
	return V1Domains.ProductDomain{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func FromProductsV1Domain(p *V1Domains.ProductDomain) Product {
	return Product{
		Id:          p.Id,
		Name:        p.Name,
		Description: p.Description,
		Price:       p.Price,
		Stock:       p.Stock,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func ToArrayOfProductsV1Domain(u *[]Product) []V1Domains.ProductDomain {
	var result []V1Domains.ProductDomain

	for _, val := range *u {
		result = append(result, val.ToV1Domain())
	}

	return result
}
