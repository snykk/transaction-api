package v1

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/datasources/records"
)

type postgreWalletRepository struct {
	conn *sqlx.DB
}

func NewWalletRepository(conn *sqlx.DB) V1Domains.WalletRepository {
	return &postgreWalletRepository{
		conn: conn,
	}
}

func (r *postgreWalletRepository) CreateByUserId(ctx context.Context, userId string) (V1Domains.WalletDomain, error) {
	query := `
        INSERT INTO wallets (wallet_id, user_id, balance, created_at)
        VALUES (uuid_generate_v4(), $1, 0, $2)
        RETURNING wallet_id, user_id, balance, created_at, updated_at
    `
	var result records.Wallet
	now := time.Now()
	err := r.conn.GetContext(ctx, &result, query, userId, now)
	if err != nil {
		return V1Domains.WalletDomain{}, err
	}

	return result.ToV1Domain(), nil
}

func (r *postgreWalletRepository) GetByUserId(ctx context.Context, userId string) (V1Domains.WalletDomain, error) {
	query := `
        SELECT 
            w.wallet_id, w.user_id, w.balance, w.created_at, w.updated_at,
				u.user_id AS "user.user_id", u.username AS "user.username", u.email AS "user.email", 
				u.password AS "user.password", u.active AS "user.active", u.role_id AS "user.role_id", 
				u.created_at AS "user.created_at", u.updated_at AS "user.updated_at"
         FROM wallets w
        INNER JOIN users u ON w.user_id = u.user_id
        WHERE w.user_id = $1
    `

	var result records.Wallet
	err := r.conn.GetContext(ctx, &result, query, userId)
	if err != nil {
		return V1Domains.WalletDomain{}, err
	}

	return result.ToV1Domain(), nil
}
