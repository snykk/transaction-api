package responses

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type ProductResponse struct {
	Id          int        `json:"product_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Stock       int        `json:"stock"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func FromProductDomainV1(b V1Domains.ProductDomain) ProductResponse {
	return ProductResponse{
		Id:          b.Id,
		Name:        b.Name,
		Description: b.Description,
		Price:       b.Price,
		Stock:       b.Stock,
		CreatedAt:   b.CreatedAt,
		UpdatedAt:   b.UpdatedAt,
	}
}

func ToProductResponseList(domains []V1Domains.ProductDomain) []ProductResponse {
	var result []ProductResponse

	for _, val := range domains {
		result = append(result, FromProductDomainV1(val))
	}

	return result
}
