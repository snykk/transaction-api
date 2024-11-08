package records

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type Wallet struct {
	Id        string     `db:"wallet_id"`
	UserId    string     `db:"user_id"`
	Balance   float64    `db:"balance"`
	User      Users      `db:"user"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt *time.Time `db:"updated_at"`
}

// Mapper
func (p *Wallet) ToV1Domain() V1Domains.WalletDomain {
	return V1Domains.WalletDomain{
		Id:        p.Id,
		UserId:    p.UserId,
		Balance:   p.Balance,
		User:      p.User.ToV1Domain(),
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func FromWalletV1Domain(p *V1Domains.WalletDomain) Wallet {
	return Wallet{
		Id:        p.Id,
		UserId:    p.UserId,
		Balance:   p.Balance,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func ToArrayOfWalletV1Domain(u *[]Wallet) []V1Domains.WalletDomain {
	var result []V1Domains.WalletDomain

	for _, val := range *u {
		result = append(result, val.ToV1Domain())
	}

	return result
}
