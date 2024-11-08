package v1

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
)

type TransactionDomain struct {
	Id              string
	WalletId        string
	Wallet          WalletDomain
	ProductId       *int // Nullable, karena transaksi deposit tidak melibatkan produk
	Product         ProductDomain
	Amount          float64
	Quantity        int
	TransactionType string
	// Status          string
	// Description     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionUsecase interface {
	Deposit(ctx context.Context, transactionDom *TransactionDomain) (domain TransactionDomain, statusCode int, err error)
	Withdraw(ctx context.Context, transactionDom *TransactionDomain) (domain TransactionDomain, statusCode int, err error)
}

type TransactionRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	Store(ctx context.Context, transactionDom *TransactionDomain) (TransactionDomain, error)
}
