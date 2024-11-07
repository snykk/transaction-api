package v1

import (
	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type postgreWalletRepository struct {
	conn *sqlx.DB
}

func NewWalletRepository(conn *sqlx.DB) V1Domains.WalletRepository {
	return &postgreWalletRepository{
		conn: conn,
	}
}
