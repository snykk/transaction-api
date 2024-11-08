package v1

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/datasources/records"
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

func (r *postgreTransactionRepository) GetByUserId(ctx context.Context, userId string) ([]V1Domains.TransactionDomain, error) {
	query := `
		SELECT 
			t.transaction_id,
			t.wallet_id,
			t.amount,
			t.transaction_type,
			t.created_at
		FROM 
			transactions t
		JOIN 
			wallets w ON t.wallet_id = w.wallet_id
		WHERE 
			w.user_id = $1
		ORDER BY 
			t.created_at DESC
	`

	fmt.Println(userId)

	var transactions []records.Transaction
	err := r.conn.SelectContext(ctx, &transactions, query, userId)
	if err != nil {
		return nil, err
	}

	return records.ToArrayOfTransactionV1Domain(&transactions), nil
}
