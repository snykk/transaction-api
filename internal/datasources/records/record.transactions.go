package records

import (
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type Transaction struct {
	Id              string    `db:"transaction_id"`
	WalletId        string    `db:"wallet_id"`
	Wallet          Wallet    `db:"wallet"`
	ProductId       *int      `db:"product_id"` // Nullable, karena transaksi deposit tidak melibatkan produk
	Product         Product   `db:"product"`
	Amount          float64   `db:"amount"`
	Quantity        *int      `db:"quantity"` // Nullable, karena transaksi deposit tidak melibatkan quantity
	TransactionType string    `db:"transaction_type"`
	CreatedAt       time.Time `db:"created_at"`
	UpdatedAt       time.Time `db:"updated_at"`
}

// Mapper
func (p *Transaction) ToV1Domain() V1Domains.TransactionDomain {
	return V1Domains.TransactionDomain{
		Id:              p.Id,
		WalletId:        p.WalletId,
		Wallet:          p.Wallet.ToV1Domain(),
		ProductId:       p.ProductId,
		Product:         p.Product.ToV1Domain(),
		Amount:          p.Amount,
		Quantity:        p.Quantity,
		TransactionType: p.TransactionType,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func FromTransactionV1Domain(p *V1Domains.TransactionDomain) Transaction {
	return Transaction{
		Id:              p.Id,
		WalletId:        p.WalletId,
		Wallet:          FromWalletV1Domain(&p.Wallet),
		ProductId:       p.ProductId,
		Product:         FromProductsV1Domain(&p.Product),
		Amount:          p.Amount,
		Quantity:        p.Quantity,
		TransactionType: p.TransactionType,
		CreatedAt:       p.CreatedAt,
		UpdatedAt:       p.UpdatedAt,
	}
}

func ToArrayOfTransactionV1Domain(u *[]Transaction) []V1Domains.TransactionDomain {
	var result []V1Domains.TransactionDomain

	for _, val := range *u {
		result = append(result, val.ToV1Domain())
	}

	return result
}
