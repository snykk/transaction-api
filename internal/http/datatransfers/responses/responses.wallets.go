package responses

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type WalletResponse struct {
	Id        int        `json:"wallet_id"`
	UserId    int        `json:"user_id"`
	Balance   float64    `json:"balance"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func FromWalletDomainV1(b V1Domains.WalletDomain) WalletResponse {
	return WalletResponse{
		Id: b.Id,

		CreatedAt: b.CreatedAt,
		UpdatedAt: b.UpdatedAt,
	}
}

func ToWalletResponseList(domains []V1Domains.WalletDomain) []WalletResponse {
	var result []WalletResponse

	for _, val := range domains {
		result = append(result, FromWalletDomainV1(val))
	}

	return result
}
