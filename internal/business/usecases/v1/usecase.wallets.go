package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

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

func (walletUC *walletUsecase) Init(ctx context.Context, userId string) (domain V1Domains.WalletDomain, statusCode int, err error) {
	existingWallet, err := walletUC.repo.GetWalletByUserId(ctx, userId)
	if err == nil {
		return existingWallet, http.StatusConflict, errors.New("wallet already exists for this user")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	walletWithoutRelation, err := walletUC.repo.CreateWalletByUserId(ctx, userId)
	if err != nil {
		return walletWithoutRelation, http.StatusInternalServerError, err
	}

	walletWithRelation, err := walletUC.repo.GetWalletByUserId(ctx, userId)
	if err != nil {
		return walletWithRelation, http.StatusInternalServerError, err
	}

	return walletWithRelation, http.StatusCreated, nil
}
