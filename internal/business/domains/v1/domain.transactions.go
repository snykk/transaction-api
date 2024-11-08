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
	Quantity        *int
	TransactionType string
	// Status          string
	// Description     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionUsecase interface {
	GetAll(ctx context.Context) (domains []TransactionDomain, statusCode int, err error)
	Deposit(ctx context.Context, transactionDom *TransactionDomain) (domain TransactionDomain, statusCode int, err error)
	Withdraw(ctx context.Context, transactionDom *TransactionDomain) (domain TransactionDomain, statusCode int, err error)
	Purchase(ctx context.Context, transactionData *TransactionDomain) (domain TransactionDomain, statusCode int, err error)
	History(ctx context.Context, userId string) (domains []TransactionDomain, statusCode int, err error)
}

type TransactionRepository interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, error)
	GetAll(ctx context.Context) ([]TransactionDomain, error)
	GetByUserId(ctx context.Context, userId string) ([]TransactionDomain, error)
}
