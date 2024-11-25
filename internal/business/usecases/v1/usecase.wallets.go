package v1

import (
	"context"
	"database/sql"
	"errors"
	"net/http"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/utils"
)

type walletUsecase struct {
	repo V1Domains.WalletRepository
}

func NewWalletUsecase(repo V1Domains.WalletRepository) V1Domains.WalletUsecase {
	return &walletUsecase{
		repo,
	}
}

func (uc *walletUsecase) GetAllWallets(ctx context.Context) ([]V1Domains.WalletDomain, int, error) {
	wallets, err := uc.repo.GetAllWallets(ctx)

	if err != nil {
		return []V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	return wallets, http.StatusOK, nil
}

func (walletUC *walletUsecase) Init(ctx context.Context, userId string) (domain V1Domains.WalletDomain, statusCode int, err error) {
	existingWallet, err := walletUC.repo.GetWalletByUserId(ctx, userId)
	if err == nil {
		return existingWallet, http.StatusConflict, errors.New("wallet already exists for this user")
	}

	if !errors.Is(err, sql.ErrNoRows) {
		return V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	_, err = walletUC.repo.CreateWalletByUserId(ctx, userId)
	if err != nil {
		return V1Domains.WalletDomain{}, http.StatusInternalServerError, err
	}

	walletWithRelationDom, err := walletUC.repo.GetWalletByUserId(ctx, userId)
	if err != nil {
		return walletWithRelationDom, http.StatusInternalServerError, err
	}

	return walletWithRelationDom, http.StatusCreated, nil
}

func (walletUC *walletUsecase) GetWalletByUserId(ctx context.Context, userId string) (V1Domains.WalletDomain, int, error) {
	walletDom, err := walletUC.repo.GetWalletByUserId(ctx, userId)

	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.WalletDomain{}, statusCode, err
	}

	return walletDom, http.StatusOK, nil
}
