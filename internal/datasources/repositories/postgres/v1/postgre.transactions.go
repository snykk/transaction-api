package v1

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/constants"
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

func (r *postgreTransactionRepository) GetAll(ctx context.Context) ([]V1Domains.TransactionDomain, error) {
	query := `SELECT transaction_id, wallet_id, amount, transaction_type, stock, created_at FROM transactionsss`
	var transactionRecord []records.Transaction
	err := r.conn.SelectContext(ctx, &transactionRecord, query)
	if err != nil {
		return nil, err
	}

	return records.ToArrayOfTransactionV1Domain(&transactionRecord), nil
}

func (r *postgreTransactionRepository) Deposit(ctx context.Context, transactionDom V1Domains.TransactionDomain) (V1Domains.TransactionDomain, error) {
	// Mulai transaksi database
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Tingkat isolasi tertinggi
		ReadOnly:  false,                 // Transaksi diperbolehkan melakukan perubahan data
	}

	tx, err := r.conn.BeginTxx(ctx, txOptions)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
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
	var wallet records.Wallet
	err = tx.GetContext(ctx, &wallet, queryGetWallet, transactionDom.Wallet.UserId)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Hitung saldo baru
	newBalance := wallet.Balance + transactionDom.Amount

	// Update saldo wallet
	queryUpdateBalance := `
		UPDATE wallets SET balance = $1, updated_at = $2
		WHERE wallet_id = $3
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, newBalance, time.Now(), wallet.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Buat transaksi baru dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction records.Transaction
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at
	`
	err = tx.GetContext(ctx, &newTransaction, queryCreateTransaction, wallet.Id, transactionDom.Amount, constants.TransactionTypeDeposit, time.Now())

	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	return newTransaction.ToV1Domain(), nil
}

func (r *postgreTransactionRepository) Withdraw(ctx context.Context, transactionDom V1Domains.TransactionDomain) (V1Domains.TransactionDomain, error) {
	// Mulai transaksi database
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Tingkat isolasi tertinggi
		ReadOnly:  false,                 // Transaksi diperbolehkan melakukan perubahan data
	}

	tx, err := r.conn.BeginTxx(ctx, txOptions)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
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
	var wallet records.Wallet
	err = tx.GetContext(ctx, &wallet, queryGetWallet, transactionDom.Wallet.UserId)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Cek apakah saldo mencukupi untuk withdraw dan tidak akan menyebabkan balance menjadi negatif
	if wallet.Balance < transactionDom.Amount {
		return V1Domains.TransactionDomain{}, errors.New("insufficient balance for withdraw")
	}

	// Hitung saldo baru setelah withdraw
	newBalance := wallet.Balance - transactionDom.Amount

	// Update saldo wallet
	queryUpdateBalance := `
		UPDATE wallets SET balance = $1, updated_at = $2
		WHERE wallet_id = $3
	`
	_, err = tx.ExecContext(ctx, queryUpdateBalance, newBalance, time.Now(), wallet.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Buat transaksi baru dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction records.Transaction
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at
	`
	err = tx.GetContext(ctx, &newTransaction, queryCreateTransaction, wallet.Id, transactionDom.Amount, constants.TransactionTypeWithdraw, time.Now())
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	return newTransaction.ToV1Domain(), nil
}

func (r *postgreTransactionRepository) Purchase(ctx context.Context, trasanctionDom V1Domains.TransactionDomain) (V1Domains.TransactionDomain, error) {
	// Mulai transaksi database
	txOptions := &sql.TxOptions{
		Isolation: sql.LevelSerializable, // Tingkat isolasi tertinggi
		ReadOnly:  false,                 // Transaksi diperbolehkan melakukan perubahan data
	}

	tx, err := r.conn.BeginTxx(ctx, txOptions)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
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
	var wallet records.Wallet
	err = tx.GetContext(ctx, &wallet, queryGetWallet, trasanctionDom.Wallet.UserId)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Ambil data produk berdasarkan productId untuk mendapatkan harga
	queryGetProduct := `
		SELECT product_id, price, stock
		FROM products
		WHERE product_id = $1
	`
	var product records.Product
	err = tx.GetContext(ctx, &product, queryGetProduct, trasanctionDom.ProductId)
	if err != nil {
		return V1Domains.TransactionDomain{}, errors.New("product not found")
	}

	// Validasi apakah stock cukup
	if product.Stock < *trasanctionDom.Quantity {
		return V1Domains.TransactionDomain{}, errors.New("insufficient product stock")
	}

	// Hitung total price berdasarkan quantity
	totalPrice := product.Price * float64(*trasanctionDom.Quantity)

	// Validasi apakah saldo cukup untuk pembelian
	if wallet.Balance < totalPrice {
		return V1Domains.TransactionDomain{}, errors.New("insufficient balance")
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
		return V1Domains.TransactionDomain{}, err
	}

	// Kurangi stock produk setelah pembelian
	newStock := product.Stock - *trasanctionDom.Quantity
	queryUpdateProductStock := `
		UPDATE products SET stock = $1 WHERE product_id = $2
	`
	_, err = tx.ExecContext(ctx, queryUpdateProductStock, newStock, product.Id)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	// Buat transaksi baru untuk pembelian dan dapatkan semua data transaksi yang dihasilkan oleh database
	var newTransaction records.Transaction
	queryCreateTransaction := `
		INSERT INTO transactions (transaction_id, wallet_id, amount, transaction_type, created_at, product_id, quantity)
		VALUES (uuid_generate_v4(), $1, $2, $3, $4, $5, $6)
		RETURNING transaction_id, wallet_id, amount, transaction_type, created_at, product_id, quantity
	`
	err = tx.GetContext(ctx, &newTransaction, queryCreateTransaction, wallet.Id, totalPrice, "purchase", time.Now(), trasanctionDom.ProductId, trasanctionDom.Quantity)
	if err != nil {
		return V1Domains.TransactionDomain{}, err
	}

	return newTransaction.ToV1Domain(), nil
}
