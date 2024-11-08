package v1

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type postgreTransactionRepository struct {
	conn *sqlx.DB
}

func NewTransactionRepository(conn *sqlx.DB) V1Domains.TransactionRepository {
	return &postgreTransactionRepository{
		conn: conn,
	}
}

func (r *postgreTransactionRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Tingkat isolasi tertinggi
		ReadOnly:  false,                 // Transaksi diperbolehkan melakukan perubahan data
	}

	return r.conn.BeginTxx(ctx, txOptions)
}

func (r *postgreTransactionRepository) Store(ctx context.Context, transactionDom *V1Domains.TransactionDomain) (V1Domains.TransactionDomain, error) {
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, status, description, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6)
		RETURNING transaction_id, wallet_id, amount, transaction_type, status, description, created_at
	`

	row := r.conn.QueryRowContext(ctx, queryCreateTransaction,
		transactionDom.WalletId,
		transactionDom.Amount,
		transactionDom.TransactionType,
		transactionDom.CreatedAt,
	)

	var newTransaction V1Domains.TransactionDomain

	err := row.Scan(
		&newTransaction.Id,
		&newTransaction.WalletId,
		&newTransaction.Amount,
		&newTransaction.TransactionType,
		&newTransaction.CreatedAt,
	)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	return newTransaction, nil
}
