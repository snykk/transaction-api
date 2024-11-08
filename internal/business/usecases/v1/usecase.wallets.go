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

func (uc *walletUsecase) GetAll(ctx context.Context) ([]V1Domains.WalletDomain, int, error) {
	wallets, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	return wallets, http.StatusOK, nil
}

func (walletUC *walletUsecase) Init(ctx context.Context, userId string) (domain V1Domains.WalletDomain, statusCode int, err error) {
	existingWallet, err := walletUC.repo.GetByUserId(ctx, userId)
	if err == nil {
		return existingWallet, http.StatusConflict, errors.New("wallet already exists for this user")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	walletWithoutRelationDom, err := walletUC.repo.CreateByUserId(ctx, userId)
	if err != nil {
		return walletWithoutRelationDom, http.StatusInternalServerError, err
	}

	walletWithRelationDom, err := walletUC.repo.GetByUserId(ctx, userId)
	if err != nil {
		return walletWithRelationDom, http.StatusInternalServerError, err
	}

	return walletWithRelationDom, http.StatusCreated, nil
}

func (walletUC *walletUsecase) GetByUserId(ctx context.Context, userId string) (V1Domains.WalletDomain, int, error) {
	walletDom, err := walletUC.repo.GetByUserId(ctx, userId)

	if err != nil {
		return V1Domains.WalletDomain{}, http.StatusNotFound, errors.New("wallet not found")
	}

	return walletDom, http.StatusOK, nil
}
