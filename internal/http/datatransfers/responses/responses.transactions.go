package responses

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type TransactionResponse struct {
	Id              string                   `json:"transaction_id"`
	WalletId        string                   `json:"wallet_id"`
	Wallet          *V1Domains.WalletDomain  `json:"wallet,omitempty"`
	ProductId       *int                     `json:"product_id,omitempty"`
	Product         *V1Domains.ProductDomain `json:"product,omitempty"`
	Amount          float64                  `json:"amount"`
	Quantity        *int                     `json:"quantity,omitempty"`
	TransactionType string                   `json:"transaction_type"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       *time.Time               `json:"updated_at,omitempty"`
}

func FromTransactionDomainV1(b V1Domains.TransactionDomain) TransactionResponse {
	return TransactionResponse{
		Id:              b.Id,
		WalletId:        b.WalletId,
		ProductId:       b.ProductId,
		Amount:          b.Amount,
		Quantity:        &b.Quantity,
		TransactionType: b.TransactionType,
		CreatedAt:       b.CreatedAt,
		UpdatedAt:       &b.UpdatedAt,
	}
}

func ToTransactionResponseList(domains []V1Domains.TransactionDomain) []TransactionResponse {
	var result []TransactionResponse

	for _, val := range domains {
		result = append(result, FromTransactionDomainV1(val))
	}

	return result
}
