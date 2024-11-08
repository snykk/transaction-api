package v1

import (
	"context"
	"errors"
	"net/http"
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type transactionUsecase struct {
	repo       V1Domains.TransactionRepository
	walletRepo V1Domains.WalletRepository
}

func NewTransactionUsecase(repo V1Domains.TransactionRepository, walletRepo V1Domains.WalletRepository) V1Domains.TransactionUsecase {
	return &transactionUsecase{
		repo,
		walletRepo,
	}
}

func (txUC *transactionUsecase) Deposit(ctx context.Context, transactionData *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi jumlah deposit harus lebih dari 0
	if transactionData.Amount <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("amount must be greater than zero")
	}

	// Mulai transaksi database
	tx, err := txUC.repo.BeginTx(ctx)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Pastikan transaksi di-rollback jika terjadi error atau panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Panic diteruskan setelah rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Ambil data wallet berdasarkan userID
	queryGetWallet := `
		SELECT wallet_id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var wallet V1Domains.WalletDomain
	err = tx.QueryRowContext(ctx, queryGetWallet, transactionData.Wallet.UserId).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusNotFound, errors.New("wallet not found")
	}

	// Hitung saldo baru
	newBalance := wallet.Balance + transactionData.Amount

	// Update saldo wallet
	queryUpdateBalance := `
		UPDATE wallets SET balance = $1, updated_at = $2
		WHERE wallet_id = $3
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, newBalance, time.Now(), wallet.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Buat transaksi baru dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction V1Domains.TransactionDomain
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at
	`
	err = tx.QueryRowContext(ctx, queryCreateTransaction, wallet.Id, transactionData.Amount, "deposit", time.Now()).
		Scan(&newTransaction.Id, &newTransaction.WalletId, &newTransaction.Amount, &newTransaction.TransactionType, &newTransaction.CreatedAt)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransaction, http.StatusOK, nil
}

func (txUC *transactionUsecase) Withdraw(ctx context.Context, transactionData *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi jumlah withdraw harus lebih dari 0
	if transactionData.Amount <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("amount must be greater than zero")
	}

	// Mulai transaksi database
	tx, err := txUC.repo.BeginTx(ctx)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Pastikan transaksi di-rollback jika terjadi error atau panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Panic diteruskan setelah rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Ambil data wallet berdasarkan userID
	queryGetWallet := `
		SELECT wallet_id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var wallet V1Domains.WalletDomain
	err = tx.QueryRowContext(ctx, queryGetWallet, transactionData.Wallet.UserId).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusNotFound, errors.New("wallet not found")
	}

	// Cek apakah saldo mencukupi untuk withdraw dan tidak akan menyebabkan balance menjadi negatif
	if wallet.Balance < transactionData.Amount {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("insufficient balance for withdraw")
	}

	// Hitung saldo baru setelah withdraw
	newBalance := wallet.Balance - transactionData.Amount

	// Update saldo wallet
	queryUpdateBalance := `
		UPDATE wallets SET balance = $1, updated_at = $2
		WHERE wallet_id = $3
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, newBalance, time.Now(), wallet.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Buat transaksi baru dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction V1Domains.TransactionDomain
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at
	`
	err = tx.QueryRowContext(ctx, queryCreateTransaction, wallet.Id, transactionData.Amount, "withdraw", time.Now()).
		Scan(&newTransaction.Id, &newTransaction.WalletId, &newTransaction.Amount, &newTransaction.TransactionType, &newTransaction.CreatedAt)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransaction, http.StatusOK, nil
}

func (txUC *transactionUsecase) Purchase(ctx context.Context, transactionData *V1Domains.TransactionDomain) (domain V1Domains.TransactionDomain, statusCode int, err error) {
	// Validasi quantity harus lebih dari 0
	if *transactionData.Quantity <= 0 {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("quantity must be greater than zero")
	}

	// Mulai transaksi database
	tx, err := txUC.repo.BeginTx(ctx)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Pastikan transaksi di-rollback jika terjadi error atau panic
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p) // Panic diteruskan setelah rollback
		} else if err != nil {
			tx.Rollback()
		} else {
			err = tx.Commit()
		}
	}()

	// Ambil data wallet berdasarkan userID
	queryGetWallet := `
		SELECT wallet_id, user_id, balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
	`
	var wallet V1Domains.WalletDomain
	err = tx.QueryRowContext(ctx, queryGetWallet, transactionData.Wallet.UserId).Scan(&wallet.Id, &wallet.UserId, &wallet.Balance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusNotFound, errors.New("wallet not found")
	}

	// Ambil data produk berdasarkan productId untuk mendapatkan harga
	queryGetProduct := `
		SELECT product_id, price, stock
		FROM products
		WHERE product_id = $1
	`
	var product V1Domains.ProductDomain
	err = tx.QueryRowContext(ctx, queryGetProduct, transactionData.ProductId).Scan(&product.Id, &product.Price, &product.Stock)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusNotFound, errors.New("product not found")
	}

	// Validasi apakah stock cukup
	if product.Stock < *transactionData.Quantity {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("insufficient product stock")
	}

	// Hitung total price berdasarkan quantity
	totalPrice := product.Price * float64(*transactionData.Quantity)

	// Validasi apakah saldo cukup untuk pembelian
	if wallet.Balance < totalPrice {
		return V1Domains.TransactionDomain{}, http.StatusBadRequest, errors.New("insufficient balance")
	}

	// Hitung saldo baru setelah pembelian
	newBalance := wallet.Balance - totalPrice

	// Update saldo wallet
	queryUpdateBalance := `
		UPDATE wallets SET balance = $1, updated_at = $2
		WHERE wallet_id = $3
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, newBalance, time.Now(), wallet.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Kurangi stock produk setelah pembelian
	newStock := product.Stock - *transactionData.Quantity
	queryUpdateProductStock := `
		UPDATE products SET stock = $1 WHERE product_id = $2
	`
	_, err = tx.ExecContext(ctx, queryUpdateProductStock, newStock, product.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Buat transaksi baru untuk pembelian dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction V1Domains.TransactionDomain
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at, product_id, quantity)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at, product_id, quantity
	`
	err = tx.QueryRowContext(ctx, queryCreateTransaction, wallet.Id, totalPrice, "purchase", time.Now(), transactionData.ProductId, transactionData.Quantity).
		Scan(&newTransaction.Id, &newTransaction.WalletId, &newTransaction.Amount, &newTransaction.TransactionType, &newTransaction.CreatedAt, &newTransaction.ProductId, &newTransaction.Quantity)
	if err != nil {
		return V1Domains.TransactionDomain{}, http.StatusInternalServerError, err
	}

	// Mengembalikan transaksi yang baru disimpan
	return newTransaction, http.StatusOK, nil
}
