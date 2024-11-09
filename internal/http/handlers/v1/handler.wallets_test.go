package v1_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
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
	"github.com/snykk/transaction-api/internal/http/datatransfers/responses"
	V1Handlers "github.com/snykk/transaction-api/internal/http/handlers/v1"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/snykk/transaction-api/pkg/helpers"
	"github.com/snykk/transaction-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtServiceWalletMock *mocks.JWTService
	walletRepoMock       *mocks.WalletRepository
	walletUsecase        V1Domains.WalletUsecase
	walletHandler        V1Handlers.WalletHandler
	ristrettoWalletMock  *mocks.RistrettoCache
	sWallet              *gin.Engine
	walletsDataFromDB    []V1Domains.WalletDomain
	walletDataFromDB     V1Domains.WalletDomain
	userWalletFromDB     V1Domains.UserDomain
	userWalletFromDB2    V1Domains.UserDomain
)

func setupWallet(t *testing.T) {
	// Initialize mock dependencies
	jwtServiceWalletMock = mocks.NewJWTService(t)
	ristrettoWalletMock = mocks.NewRistrettoCache(t)
	walletRepoMock = mocks.NewWalletRepository(t)
	walletUsecase = V1Usecases.NewWalletUsecase(walletRepoMock)
	walletHandler = V1Handlers.NewWalletHandler(walletUsecase, ristrettoWalletMock)

	// Mock users and wallets data
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

	// Setup Gin engine with middleware for authentication
	sWallet = gin.Default()
	// sWallet.Use(lazyAuthCommonWallet)
}

// Mock lazy authentication
func lazyAuthCommonWallet(ctx *gin.Context) {
	pass, _ := helpers.GenerateHash(walletDataFromDB.User.Password)
	jwtClaims := jwt.JwtCustomClaim{
		UserID:   walletDataFromDB.UserId,
		IsAdmin:  false,
		Email:    walletDataFromDB.User.Email,
		Password: pass,
		StandardClaims: dgriJWT.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    walletDataFromDB.User.Username,
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func lazyAuthAdminWallet(ctx *gin.Context) {
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

func TestInit(t *testing.T) {
	setupWallet(t)

	sWallet.Use(lazyAuthCommonWallet)

	// Define the route for testing
	sWallet.POST(constants.EndpointV1+"/wallets/init", walletHandler.Init)

	t.Run("Success - Wallet Creation", func(t *testing.T) {
		reqBody, _ := json.Marshal(new(int))

		// Set up mock expectations
		ristrettoWalletMock.On("Del", mock.AnythingOfType("string")).Once()
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, sql.ErrNoRows).Once()
		walletRepoMock.Mock.On("CreateWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(V1Domains.WalletDomain{}, nil).Once()
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/wallets/init", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sWallet.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "wallet created successfully")
	})

	t.Run("Failure - Wallet Already Exists", func(t *testing.T) {
		reqBody, _ := json.Marshal(new(int))

		// Mock the behavior for existing wallet
		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, "aaaa-bbbb-cccc").Return(walletDataFromDB, nil).Once()

		// Perform the HTTP request
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/wallets/init", bytes.NewReader(reqBody))
		r.Header.Set("Content-Type", "application/json")

		// Serve request
		sWallet.ServeHTTP(w, r)
		body := w.Body.String()

		// Assert the HTTP response
		assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "wallet already exists for this user")
	})
}

func TestGetAll(t *testing.T) {
	setupWallet(t)

	sWallet.Use(lazyAuthAdminWallet)

	// Define route
	sWallet.GET(constants.EndpointV1+"/wallets", walletHandler.GetAll)

	t.Run("When Success and Data Available in Cache", func(t *testing.T) {
		ristrettoWalletMock.On("Get", "wallets").Return(responses.ToWalletResponseList(walletsDataFromDB)).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets", nil)

		r.Header.Set("Content-Type", "application/json")

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "wallet data fetched successfully")
	})

	t.Run("When Success and Data Not Available in Cache", func(t *testing.T) {
		ristrettoWalletMock.On("Get", "wallets").Return(nil).Once()
		// walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return(walletsDataFromDB, http.StatusOK, nil).Once()
		walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return(walletsDataFromDB, nil).Once()
		ristrettoWalletMock.On("Set", "wallets", mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets", nil)

		r.Header.Set("Content-Type", "application/json")

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, body, "wallet data fetched successfully")
	})

	t.Run("When Wallet Data is Empty", func(t *testing.T) {
		ristrettoWalletMock.On("Get", "wallets").Return(nil).Once()
		walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return([]V1Domains.WalletDomain{}, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets", nil)

		r.Header.Set("Content-Type", "application/json")

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, body, "wallet data is empty")
	})

	t.Run("When Error Occurs", func(t *testing.T) {
		ristrettoWalletMock.On("Get", "wallets").Return(nil).Once()
		walletRepoMock.Mock.On("GetAllWallets", mock.Anything).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets", nil)

		r.Header.Set("Content-Type", "application/json")

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, body, "database error")
	})
}

func TestInfo(t *testing.T) {
	setupWallet(t)
	sWallet.Use(lazyAuthCommonWallet)
	// Define route
	sWallet.GET(constants.EndpointV1+"/wallets/info", walletHandler.Info)

	t.Run("When Success", func(t *testing.T) {
		ristrettoWalletMock.On("Get", mock.Anything).Return(nil).Once()

		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, mock.AnythingOfType("string")).Return(walletDataFromDB, nil).Once()

		ristrettoWalletMock.On("Set", mock.Anything, mock.Anything).Return(nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets/info", nil)

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "wallet fetched successfully")
	})

	t.Run("When Wallet Not Found", func(t *testing.T) {
		ristrettoWalletMock.On("Get", mock.Anything).Return(nil).Once()

		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, mock.AnythingOfType("string")).Return(V1Domains.WalletDomain{}, sql.ErrNoRows).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets/info", nil)

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		assert.Contains(t, body, "sql: no rows in result set")
	})

	t.Run("When Database Error Occurs", func(t *testing.T) {
		ristrettoWalletMock.On("Get", mock.Anything).Return(nil).Once()

		walletRepoMock.Mock.On("GetWalletByUserId", mock.Anything, mock.AnythingOfType("string")).Return(V1Domains.WalletDomain{}, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/wallets/info", nil)

		sWallet.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, body, "database error")
	})
}
