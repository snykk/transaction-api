package v1_test

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	V1Usecases "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	walletRepoMock    *mocks.WalletRepository
	walletUsecase     V1Domains.WalletUsecase
	walletsDataFromDB []V1Domains.WalletDomain
	walletDataFromDB  V1Domains.WalletDomain
	userWalletFromDB  V1Domains.UserDomain
	userWalletFromDB2 V1Domains.UserDomain
)

func setupWallet(t *testing.T) {
	walletRepoMock = mocks.NewWalletRepository(t)
	walletUsecase = V1Usecases.NewWalletUsecase(walletRepoMock)

	currentTime := time.Now()

	userWalletFromDB = V1Domains.UserDomain{
		ID:       "aaaa-bbbb-cccc",
		Username: "patrick",
		Email:    "patrick@gmail.com",
		Password: "sjdflakjsdfldks",
		Active:   true,
		RoleID:   2,
	}
	userWalletFromDB2 = V1Domains.UserDomain{
		ID:       "asdf-asdf-asdf",
		Username: "asdf",
		Email:    "asdf@gmail.com",
		Password: "asdfsdafsasad",
		Active:   true,
		RoleID:   2,
	}

	walletsDataFromDB = []V1Domains.WalletDomain{
		{
			Id:        "xxxx-yyyy-zzzz",
			UserId:    "aaaa-bbbb-cccc",
			Balance:   0,
			User:      userWalletFromDB,
			CreatedAt: currentTime,
			UpdatedAt: &currentTime,
		},
		{
			Id:        "yyyy-zzzz-aaaa",
			UserId:    "asdf-asdf-asdf",
			Balance:   0,
			User:      userWalletFromDB2,
			CreatedAt: currentTime,
			UpdatedAt: &currentTime,
		},
	}

	walletDataFromDB = walletsDataFromDB[0]
}

func TestGetAllWallets(t *testing.T) {
	setupWallet(t)

	t.Run("When Success Get Wallets Data", func(t *testing.T) {
		// Mock the repository to return a list of wallets from the database
		walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return(walletsDataFromDB, nil).Once()

		// Call the method
		result, statusCode, err := walletUsecase.GetAllWallets(context.Background())

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")

		// Assert that the result contains the same data as the mock response
		assert.Len(t, result, len(walletsDataFromDB), "The number of wallets in the result should match the mock data")

		// Check individual wallets
		for i, wallet := range result {
			t.Run(fmt.Sprintf("Check Wallet %d", i+1), func(t *testing.T) {
				assert.Equal(t, walletsDataFromDB[i].Id, wallet.Id, "Wallet ID should match")
				assert.Equal(t, walletsDataFromDB[i].UserId, wallet.UserId, "User ID should match")
				assert.Equal(t, walletsDataFromDB[i].Balance, wallet.Balance, "Balance should match")
				assert.Equal(t, walletsDataFromDB[i].User.Email, wallet.User.Email, "User email should match")
				assert.NotNil(t, wallet.CreatedAt, "CreatedAt should not be nil")
				assert.NotNil(t, wallet.UpdatedAt, "UpdatedAt should not be nil")
			})
		}
	})

	t.Run("When Failure Get Wallets Data", func(t *testing.T) {
		// Mock the repository to return an error when fetching all wallets
		walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return([]V1Domains.WalletDomain{}, errors.New("get all wallets failed")).Once()

		// Call the method
		result, statusCode, err := walletUsecase.GetAllWallets(context.Background())

		// Assertions
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Equal(t, []V1Domains.WalletDomain{}, result, "Result should be empty when there's an error")
		assert.Equal(t, "get all wallets failed", err.Error(), "Error message should match")
	})
}

func TestInit(t *testing.T) {
	setupWallet(t)

	t.Run("When Wallet Doesn't Exist and Creation Succeeds", func(t *testing.T) {
		// Mock the behavior when no wallet exists for the user and creation is successful
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, sql.ErrNoRows).Once()
		walletRepoMock.Mock.On("CreateWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()

		// Call the method
		result, statusCode, err := walletUsecase.Init(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusCreated, statusCode, "Status code should be Created (201)")
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, walletDataFromDB.Id, result.Id, "Wallet ID should match")
		assert.Equal(t, walletDataFromDB.UserId, result.UserId, "User ID should match")
		assert.Equal(t, walletDataFromDB.Balance, result.Balance, "Wallet balance should match")
		assert.Equal(t, walletDataFromDB.User.Email, result.User.Email, "User email should match")
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Wallet Already Exists", func(t *testing.T) {
		// Mock the behavior when wallet already exists for the user
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()

		// Call the method
		result, statusCode, err := walletUsecase.Init(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusConflict, statusCode, "Status code should be Conflict (409)")
		assert.Equal(t, "wallet already exists for this user", err.Error(), "Error message should match")
		assert.Equal(t, walletDataFromDB, result, "Result should be same")
	})

	t.Run("When Error Occurs While Creating Wallet", func(t *testing.T) {
		// Mock the behavior when an error occurs while creating the wallet
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, sql.ErrNoRows).Once()
		walletRepoMock.Mock.On("CreateWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, errors.New("create wallet failed")).Once()

		// Call the method
		result, statusCode, err := walletUsecase.Init(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Equal(t, "create wallet failed", err.Error(), "Error message should match")
		assert.Equal(t, V1Domains.WalletDomain{}, result, "Result should be empty")
	})
}

func TestGetWalletByUserId(t *testing.T) {
	setupWallet(t)

	t.Run("When Success Get Wallet By User ID", func(t *testing.T) {
		// Mock the behavior when getting the wallet by user ID is successful
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()

		// Call the method
		result, statusCode, err := walletUsecase.GetWalletByUserId(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")
		assert.NotNil(t, result, "Result should not be nil")
		assert.Equal(t, walletDataFromDB.Id, result.Id, "Wallet ID should match")
		assert.Equal(t, walletDataFromDB.UserId, result.UserId, "User ID should match")
		assert.Equal(t, walletDataFromDB.Balance, result.Balance, "Wallet balance should match")
		assert.Equal(t, walletDataFromDB.User.Email, result.User.Email, "User email should match")
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Wallet Not Found", func(t *testing.T) {
		// Mock the behavior when the wallet does not exist for the given user ID
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, sql.ErrNoRows).Once()

		// Call the method
		result, statusCode, err := walletUsecase.GetWalletByUserId(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusNotFound, statusCode, "Status code should be Not Found (404)")
		assert.Equal(t, "sql: no rows in result set", err.Error(), "Error message should match")
		assert.Equal(t, V1Domains.WalletDomain{}, result, "Result should be empty")
	})

	t.Run("When Error Occurs While Fetching Wallet", func(t *testing.T) {
		// Mock the behavior when an error occurs while fetching the wallet
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, errors.New("failed to fetch wallet")).Once()

		// Call the method
		result, statusCode, err := walletUsecase.GetWalletByUserId(context.Background(), "aaaa-bbbb-cccc")

		// Assertions
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Equal(t, "failed to fetch wallet", err.Error(), "Error message should match")
		assert.Equal(t, V1Domains.WalletDomain{}, result, "Result should be empty")
	})
}
