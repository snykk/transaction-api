package v1

import (
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type walletUsecase struct {
	repo V1Domains.WalletRepository
}

func NewWalletUsecase(repo V1Domains.WalletRepository) V1Domains.WalletUsecase {
	return &walletUsecase{
		repo,
	}
}
