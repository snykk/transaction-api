package v1_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	V1Usecases "github.com/snykk/transaction-api/internal/business/usecases/v1"
	PostgresRepo "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	transactionRepoMock    *mocks.TransactionRepository
	transactionUsecase     V1Domains.TransactionUsecase
	transactionsDataFromDB []V1Domains.TransactionDomain
	transactionDataFromDB  V1Domains.TransactionDomain
)

func setupTransaction(t *testing.T) {
	transactionRepoMock = mocks.NewTransactionRepository(t)
	transactionUsecase = V1Usecases.NewTransactionUsecase(transactionRepoMock)

	productId1 := 1
	quantity1 := 200
	productId2 := 2
	quantity2 := 400

	current_time := time.Now()

	transactionsDataFromDB = []V1Domains.TransactionDomain{
		{
			Id:       "asdfdasf",
			WalletId: "adfasf",
			Wallet: V1Domains.WalletDomain{
				Id:        "sdfas",
				UserId:    "asdfsa",
				Balance:   200,
				CreatedAt: current_time,
				User: V1Domains.UserDomain{
					ID:        "asdfads",
					Username:  "asdfsdafjl",
					Email:     "asdfaf",
					Password:  "asdfasflj",
					Active:    true,
					RoleID:    2,
					CreatedAt: current_time,
				},
			},
			ProductId:       &productId1, // Nullable, karena transaksi deposit tidak melibatkan transaksi
			Product:         V1Domains.ProductDomain{},
			Amount:          200,
			Quantity:        &quantity1,
			TransactionType: "asdf",
			CreatedAt:       current_time,
		},
		{
			Id:       "asdfdasf",
			WalletId: "adfasf",
			Wallet: V1Domains.WalletDomain{
				Id:        "sdfas",
				UserId:    "asdfsa",
				Balance:   200,
				CreatedAt: current_time,
				User: V1Domains.UserDomain{
					ID:        "asdfads",
					Username:  "asdfsdafjl",
					Email:     "asdfaf",
					Password:  "asdfasflj",
					Active:    true,
					RoleID:    2,
					CreatedAt: current_time,
				},
			},
			ProductId:       &productId2, // Nullable, karena transaksi deposit tidak melibatkan transaksi
			Product:         V1Domains.ProductDomain{},
			Amount:          200,
			Quantity:        &quantity2,
			TransactionType: "asdf",
			CreatedAt:       current_time,
		},
	}

	transactionDataFromDB = transactionsDataFromDB[0]
}

func TestDepositTransaction(t *testing.T) {
	setupTransaction(t)

	t.Run("When Success Transaction Depo", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: transactionDataFromDB.Amount,
		}

		// Mock repository untuk mengembalikan data transaksi yang berhasil disimpan
		transactionRepoMock.Mock.On("Deposit", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		// Memanggil method Deposit
		result, statusCode, err := transactionUsecase.Deposit(context.Background(), req.ToDomain())

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusCreated, statusCode, "Status code should be Created (201)")

		// Pastikan hasilnya sesuai dengan data yang dimock
		assert.Equal(t, transactionDataFromDB.Id, result.Id, "Transaction ID should match")
		assert.Equal(t, transactionDataFromDB.WalletId, result.WalletId, "Transaction wallet_id should match")
		assert.Equal(t, transactionDataFromDB.Wallet, result.Wallet, "Transaction wallet should match")
		assert.Equal(t, transactionDataFromDB.Amount, result.Amount, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.Quantity, result.Quantity, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.TransactionType, result.TransactionType, "Transaction amount should match")

		// Pastikan atribut waktu tidak bernilai nil
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Failure | Invalid Amount", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: -200,
		}

		// Memanggil method Deposit
		result, statusCode, err := transactionUsecase.Deposit(context.Background(), req.ToDomain())

		// Assertions
		// Assertions untuk memastikan adanya error
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusBadRequest, statusCode, "Status code should be Bad Request (400)")
		assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
		assert.Equal(t, err, V1Usecases.ErrAmountMustGreateThanZero, "Error message should match")

	})

}

func TestWithdrawTransaction(t *testing.T) {
	setupTransaction(t)

	req := requests.TransactionDepositOrWithdrawRequest{
		Amount: transactionDataFromDB.Amount,
	}

	t.Run("When Success Transaction Withdraw", func(t *testing.T) {
		// Mock repository untuk mengembalikan data transaksi yang berhasil disimpan
		transactionRepoMock.Mock.On("Withdraw", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		// Memanggil method Withdraw
		result, statusCode, err := transactionUsecase.Withdraw(context.Background(), req.ToDomain())

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusCreated, statusCode, "Status code should be Created (201)")

		// Pastikan hasilnya sesuai dengan data yang dimock
		assert.Equal(t, transactionDataFromDB.Id, result.Id, "Transaction ID should match")
		assert.Equal(t, transactionDataFromDB.WalletId, result.WalletId, "Transaction wallet_id should match")
		assert.Equal(t, transactionDataFromDB.Wallet, result.Wallet, "Transaction wallet should match")
		assert.Equal(t, transactionDataFromDB.Amount, result.Amount, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.Quantity, result.Quantity, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.TransactionType, result.TransactionType, "Transaction amount should match")

		// Pastikan atribut waktu tidak bernilai nil
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Failure | Invalid Amount", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: -200,
		}

		// Memanggil method Deposit
		result, statusCode, err := transactionUsecase.Withdraw(context.Background(), req.ToDomain())

		// Assertions
		// Assertions untuk memastikan adanya error
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusBadRequest, statusCode, "Status code should be Bad Request (400)")
		assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
		assert.Equal(t, err, V1Usecases.ErrAmountMustGreateThanZero, "Error message should match")

	})

}

func TestPurchaseTransaction(t *testing.T) {
	setupTransaction(t)

	t.Run("When Success Transaction Purchase", func(t *testing.T) {
		req := requests.TransactionPurchaseRequest{
			ProductId: 1,
			Quantity:  200,
		}

		// Mock repository untuk mengembalikan data transaksi yang berhasil disimpan
		transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		// Memanggil method Purchase
		result, statusCode, err := transactionUsecase.Purchase(context.Background(), req.ToDomain())

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusCreated, statusCode, "Status code should be Created (201)")

		// Pastikan hasilnya sesuai dengan data yang dimock
		assert.Equal(t, transactionDataFromDB.Id, result.Id, "Transaction ID should match")
		assert.Equal(t, transactionDataFromDB.WalletId, result.WalletId, "Transaction wallet_id should match")
		assert.Equal(t, transactionDataFromDB.Wallet, result.Wallet, "Transaction wallet should match")
		assert.Equal(t, transactionDataFromDB.ProductId, result.ProductId, "Transaction product_id should match")
		assert.Equal(t, transactionDataFromDB.Product, result.Product, "Transaction product should match")
		assert.Equal(t, transactionDataFromDB.Amount, result.Amount, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.Quantity, result.Quantity, "Transaction amount should match")
		assert.Equal(t, transactionDataFromDB.TransactionType, result.TransactionType, "Transaction amount should match")

		// Pastikan atribut waktu tidak bernilai nil
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Failure", func(t *testing.T) {
		t.Run("Invalid Quantity", func(t *testing.T) {
			req := requests.TransactionPurchaseRequest{
				ProductId: 1,
				Quantity:  -2,
			}

			// Memanggil method Deposit
			result, statusCode, err := transactionUsecase.Purchase(context.Background(), req.ToDomain())

			// Assertions
			// Assertions untuk memastikan adanya error
			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusBadRequest, statusCode, "Status code should be Bad Request (400)")
			assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
			assert.Equal(t, err, V1Usecases.ErrQuantityMustGreaterThanZero, "Error message should match")
		})

		t.Run("Product Not Found", func(t *testing.T) {
			req := requests.TransactionPurchaseRequest{
				ProductId: 999,
				Quantity:  200,
			}

			transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(V1Domains.TransactionDomain{}, PostgresRepo.ErrProductNotFound).Once()

			// Memanggil method Deposit
			result, statusCode, err := transactionUsecase.Purchase(context.Background(), req.ToDomain())

			// Assertions
			// Assertions untuk memastikan adanya error
			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusNotFound, statusCode, "Status code should be Bad Request (404)")
			assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
			assert.Equal(t, err, PostgresRepo.ErrProductNotFound, "Error message should match")
		})

		t.Run("Insufficient Product Stock", func(t *testing.T) {
			req := requests.TransactionPurchaseRequest{
				ProductId: 1,
				Quantity:  999,
			}

			transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(V1Domains.TransactionDomain{}, PostgresRepo.ErrInsufficientProductStock).Once()

			// Memanggil method Deposit
			result, statusCode, err := transactionUsecase.Purchase(context.Background(), req.ToDomain())

			// Assertions
			// Assertions untuk memastikan adanya error
			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusUnprocessableEntity, statusCode, "Status code should be Bad Request (422)")
			assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
			assert.Equal(t, err, PostgresRepo.ErrInsufficientProductStock, "Error message should match")
		})

		t.Run("Insufficient Ballance", func(t *testing.T) {
			req := requests.TransactionPurchaseRequest{
				ProductId: 1,
				Quantity:  200,
			}

			transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(V1Domains.TransactionDomain{}, PostgresRepo.ErrInsufficientBalance).Once()

			// Memanggil method Deposit
			result, statusCode, err := transactionUsecase.Purchase(context.Background(), req.ToDomain())

			// Assertions
			// Assertions untuk memastikan adanya error
			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusUnprocessableEntity, statusCode, "Status code should be Bad Request (422)")
			assert.Equal(t, "", result.Id, "Transaction ID should be blank string on failure")
			assert.Equal(t, err, PostgresRepo.ErrInsufficientBalance, "Error message should match")
		})

	})

}
