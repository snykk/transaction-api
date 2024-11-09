package v1

import (
	"context"
	"net/http"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/utils"
)

type transactionUsecase struct {
	repo V1Domains.TransactionRepository
}

func NewTransactionUsecase(repo V1Domains.TransactionRepository) V1Domains.TransactionUsecase {
	return &transactionUsecase{
		repo,
	}
}

func (txUC *transactionUsecase) Deposit(ctx context.Context, transactionDom *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi jumlah deposit harus lebih dari 0
	if transactionDom.Amount <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, ErrAmountMustGreateThanZero
	}

	newTransactionDom, err := txUC.repo.Deposit(ctx, *transactionDom)
	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.TransactionDomain{}, statusCode, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransactionDom, http.StatusCreated, nil
}

func (txUC *transactionUsecase) Withdraw(ctx context.Context, transactionData *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi jumlah withdraw harus lebih dari 0
	if transactionData.Amount <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, ErrAmountMustGreateThanZero
	}

	newTransactionDom, err := txUC.repo.Withdraw(ctx, *transactionData)
	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.TransactionDomain{}, statusCode, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransactionDom, http.StatusCreated, nil
}

func (txUC *transactionUsecase) Purchase(ctx context.Context, transactionData *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi quantity harus lebih dari 0
	if *transactionData.Quantity <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, ErrQuantityMustGreaterThanZero
	}

	newTransactionDom, err := txUC.repo.Purchase(ctx, *transactionData)
	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.TransactionDomain{}, statusCode, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransactionDom, http.StatusCreated, nil
}

func (uc *transactionUsecase) History(ctx context.Context, userId string) ([]V1Domains.TransactionDomain, int, error) {
	userTransactionHistoryDom, err := uc.repo.GetByUserId(ctx, userId)

	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return []V1Domains.TransactionDomain{}, statusCode, err
	}

	return userTransactionHistoryDom, http.StatusOK, nil
}

func (uc *transactionUsecase) GetAll(ctx context.Context) ([]V1Domains.TransactionDomain, int, error) {
	transactionDom, err := uc.repo.GetAll(ctx)

	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return []V1Domains.TransactionDomain{}, statusCode, err
	}

	return transactionDom, http.StatusOK, nil
}
