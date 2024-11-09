package v1_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dgriJWT "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	V1Usecases "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/config"
	"github.com/snykk/transaction-api/internal/constants"
	PostgresRepo "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	V1Handlers "github.com/snykk/transaction-api/internal/http/handlers/v1"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/snykk/transaction-api/pkg/helpers"
	"github.com/snykk/transaction-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtServiceTransactionMock *mocks.JWTService
	transactionRepoMock       *mocks.TransactionRepository
	transactionUsecase        V1Domains.TransactionUsecase
	transactionHandler        V1Handlers.TransactionHandler
	ristrettoTransactiontMock *mocks.RistrettoCache
	sTransaction              *gin.Engine
	transactionsDataFromDB    []V1Domains.TransactionDomain
	transactionDataFromDB     V1Domains.TransactionDomain
)

func setupTransaction(t *testing.T) {
	// Initialize mock dependencies
	jwtServiceTransactionMock = mocks.NewJWTService(t)
	ristrettoTransactiontMock = mocks.NewRistrettoCache(t)
	transactionRepoMock = mocks.NewTransactionRepository(t)
	transactionUsecase = V1Usecases.NewTransactionUsecase(transactionRepoMock)
	transactionHandler = V1Handlers.NewTransactionHandler(transactionUsecase, ristrettoTransactiontMock)

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

	// Setup Gin engine with middleware for authentication
	sTransaction = gin.Default()
	// sTransaction.Use(lazyAuthCommonTransaction)
}

// Mock lazy authentication
func lazyAuthCommonTransaction(ctx *gin.Context) {
	pass, _ := helpers.GenerateHash(transactionDataFromDB.Wallet.User.Password)
	jwtClaims := jwt.JwtCustomClaim{
		UserID:   transactionDataFromDB.Wallet.User.ID,
		IsAdmin:  false,
		Email:    transactionDataFromDB.Wallet.User.Email,
		Password: pass,
		StandardClaims: dgriJWT.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    transactionDataFromDB.Wallet.User.Username,
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func lazyAuthAdminTransaction(ctx *gin.Context) {
	pass, _ := helpers.GenerateHash("asdfsasaf")
	jwtClaims := jwt.JwtCustomClaim{
		UserID:   "asdfsda",
		IsAdmin:  true,
		Email:    "asdf@gmail.com",
		Password: pass,
		StandardClaims: dgriJWT.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    "asdfafasa",
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func TestDeposit(t *testing.T) {
	setupTransaction(t)

	sTransaction.Use(lazyAuthCommonTransaction)

	// Define the route for testing
	sTransaction.POST(constants.EndpointV1+"/transactions/deposit", transactionHandler.Deposit)

	t.Run("Success - Deposit Transaction", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: 200,
		}

		reqBody, _ := json.Marshal(req)

		// Set up mock expectations
		transactionRepoMock.Mock.On("Deposit", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/deposit", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "deposit completed successfully")
	})

	t.Run("Failure - Invalid Amount", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: -200,
		}
		reqBody, _ := json.Marshal(req)

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/deposit", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "Field validation for 'Amount' failed on the 'gt'")
	})
}

func TestWithdraw(t *testing.T) {
	setupTransaction(t)

	sTransaction.Use(lazyAuthCommonTransaction)

	// Define the route for testing
	sTransaction.POST(constants.EndpointV1+"/transactions/withdraw", transactionHandler.Withdraw)

	t.Run("Success - Withdraw Transaction", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: 200,
		}

		reqBody, _ := json.Marshal(req)

		// Set up mock expectations
		transactionRepoMock.Mock.On("Withdraw", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/withdraw", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "withdraw completed successfully")
	})

	t.Run("Failure - Invalid Amount", func(t *testing.T) {
		req := requests.TransactionDepositOrWithdrawRequest{
			Amount: -200,
		}
		reqBody, _ := json.Marshal(req)

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/withdraw", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "Field validation for 'Amount' failed on the 'gt'")
	})
}

func TestWithPurchase(t *testing.T) {
	setupTransaction(t)

	sTransaction.Use(lazyAuthCommonTransaction)

	// Define the route for testing
	sTransaction.POST(constants.EndpointV1+"/transactions/purchase", transactionHandler.Purchase)

	t.Run("Success - Purchase Transaction", func(t *testing.T) {
		req := requests.TransactionPurchaseRequest{
			ProductId: 1,
			Quantity:  200,
		}

		reqBody, _ := json.Marshal(req)

		// Set up mock expectations
		transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(transactionDataFromDB, nil).Once()

		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string"), mock.AnythingOfType("string"), mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Once()
		ristrettoTransactiontMock.On("Del", mock.AnythingOfType("string")).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/purchase", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "purchase successful")
	})

	t.Run("Failure - Invalid Quantity", func(t *testing.T) {
		req := requests.TransactionPurchaseRequest{
			ProductId: 1,
			Quantity:  -1,
		}
		reqBody, _ := json.Marshal(req)

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/purchase", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "Error:Field validation for 'Quantity' failed on the 'gt'")
	})

	t.Run("Failure - Product Not Found", func(t *testing.T) {
		req := requests.TransactionPurchaseRequest{
			ProductId: 999,
			Quantity:  200,
		}

		reqBody, _ := json.Marshal(req)

		transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(V1Domains.TransactionDomain{}, PostgresRepo.ErrProductNotFound).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/purchase", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "product not found")
	})

	t.Run("Failure - Insufficient Product Stock", func(t *testing.T) {
		req := requests.TransactionPurchaseRequest{
			ProductId: 1,
			Quantity:  999,
		}

		reqBody, _ := json.Marshal(req)

		transactionRepoMock.Mock.On("Purchase", mock.Anything, mock.AnythingOfType("v1.TransactionDomain")).Return(V1Domains.TransactionDomain{}, PostgresRepo.ErrInsufficientProductStock).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/transactions/purchase", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sTransaction.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusUnprocessableEntity, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "insufficient product stock")
	})
}
