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
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	"github.com/snykk/transaction-api/internal/http/datatransfers/responses"
	V1Handlers "github.com/snykk/transaction-api/internal/http/handlers/v1"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/snykk/transaction-api/pkg/helpers"
	"github.com/snykk/transaction-api/pkg/jwt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jwtServiceProductMock *mocks.JWTService
	productRepoMock       *mocks.ProductRepository
	ProductUsecase        V1Domains.ProductUsecase
	ProductHandler        V1Handlers.ProductHandler
	ristrettoProductMock  *mocks.RistrettoCache
	sProduct              *gin.Engine
	productsDataFromDB    []V1Domains.ProductDomain
	productDataFromDB     V1Domains.ProductDomain
)

func setupProduct(t *testing.T) {
	// Initialize mock dependencies
	jwtServiceProductMock = mocks.NewJWTService(t)
	ristrettoProductMock = mocks.NewRistrettoCache(t)
	productRepoMock = mocks.NewProductRepository(t)
	ProductUsecase = V1Usecases.NewProductUsecase(productRepoMock)
	ProductHandler = V1Handlers.NewProductHandler(ProductUsecase, ristrettoProductMock)

	// Mock users and products data
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

	productsDataFromDB = []V1Domains.ProductDomain{
		{
			Id:          1,
			Name:        "keyboard",
			Description: "lorem ipsum dolor sit amet",
			Price:       20.0,
			Stock:       234,
			CreatedAt:   currentTime,
			UpdatedAt:   &currentTime,
		},
		{
			Id:          2,
			Name:        "mouse",
			Description: "lorem ipsum dolor sit amet",
			Price:       29.0,
			Stock:       10,
			CreatedAt:   currentTime,
			UpdatedAt:   &currentTime,
		},
	}

	productDataFromDB = productsDataFromDB[0]

	// Setup Gin engine with middleware for authentication
	sProduct = gin.Default()
}

// Mock lazy authentication
func lazyAuthCommonProduct(ctx *gin.Context) {
	pass, _ := helpers.GenerateHash("sadfaas")
	jwtClaims := jwt.JwtCustomClaim{
		UserID:   "adsfdas",
		IsAdmin:  false,
		Email:    "asdfasf",
		Password: pass,
		StandardClaims: dgriJWT.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTExpired)).Unix(),
			Issuer:    "asdfas",
			IssuedAt:  time.Now().Unix(),
		},
	}
	ctx.Set(constants.CtxAuthenticatedUserKey, jwtClaims)
}

func lazyAuthAdminProduct(ctx *gin.Context) {
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

func TestGetAllProduct(t *testing.T) {
	setupProduct(t)

	sProduct.Use(lazyAuthCommonProduct)

	// Define route
	sProduct.GET(constants.EndpointV1+"/products", ProductHandler.GetAll)

	t.Run("When Success and Data Available in Cache", func(t *testing.T) {
		ristrettoProductMock.On("Get", "products").Return(responses.ToProductResponseList(productsDataFromDB)).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products", nil)

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, w.Result().Header.Get("Content-Type"), "application/json")
		assert.Contains(t, body, "product data fetched successfully")
	})

	t.Run("When Success and Data Not Available in Cache", func(t *testing.T) {
		ristrettoProductMock.On("Get", "products").Return(nil).Once()
		// productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return(productsDataFromDB, http.StatusOK, nil).Once()
		productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return(productsDataFromDB, nil).Once()
		ristrettoProductMock.On("Set", "products", mock.Anything).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products", nil)

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, body, "product data fetched successfully")
	})

	t.Run("When Product Data is Empty", func(t *testing.T) {
		ristrettoProductMock.On("Get", "products").Return(nil).Once()
		productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return([]V1Domains.ProductDomain{}, nil).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products", nil)

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Contains(t, body, "product data is empty")
	})

	t.Run("When Error Occurs", func(t *testing.T) {
		ristrettoProductMock.On("Get", "products").Return(nil).Once()
		productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return(nil, errors.New("database error")).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products", nil)

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		body := w.Body.String()

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		assert.Contains(t, body, "database error")
	})
}

func TestStoreProduct(t *testing.T) {
	setupProduct(t)

	sProduct.Use(lazyAuthAdminProduct)
	sProduct.POST(constants.EndpointV1+"/products", ProductHandler.Store)

	t.Run("When Success", func(t *testing.T) {
		req := requests.ProductRequest{
			Name:        "keyboard",
			Description: "lorem ipsum dolor sit amet",
			Price:       20.0,
			Stock:       234,
		}

		reqBody, _ := json.Marshal(req)

		productRepoMock.Mock.On("StoreProduct", mock.Anything, mock.Anything).Return(productDataFromDB, nil).Once()
		ristrettoProductMock.On("Del", "products").Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/products", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusCreated, w.Result().StatusCode)
	})

	t.Run("When Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, constants.EndpointV1+"/products", nil)

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

func TestGetProductById(t *testing.T) {
	setupProduct(t)

	sProduct.Use(lazyAuthCommonProduct)
	sProduct.GET(constants.EndpointV1+"/products/:id", ProductHandler.GetById)

	t.Run("When Success and Data Available in Cache", func(t *testing.T) {
		ristrettoProductMock.On("Get", mock.Anything).Return(responses.FromProductDomainV1(productDataFromDB)).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products/1", nil)

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("When Product Not Found", func(t *testing.T) {
		ristrettoProductMock.On("Get", mock.Anything).Return(nil).Once()
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.Anything).Return(V1Domains.ProductDomain{}, sql.ErrNoRows).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, constants.EndpointV1+"/products/1", nil)

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestUpdateProduct(t *testing.T) {
	setupProduct(t)

	sProduct.Use(lazyAuthAdminProduct)

	sProduct.PUT(constants.EndpointV1+"/products/:id", ProductHandler.Update)

	t.Run("When Success", func(t *testing.T) {
		req := requests.ProductRequest{
			Name:        "keyboard ye",
			Description: "lorem ipsum dolor sit amet",
			Price:       20.0,
			Stock:       234,
		}

		reqBody, _ := json.Marshal(req)

		productRepoMock.Mock.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.Anything).Return(productDataFromDB, nil).Once()

		ristrettoProductMock.On("Del", "products", "product/product_id:1").Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, constants.EndpointV1+"/products/1", bytes.NewReader(reqBody))

		r.Header.Set("Content-Type", "application/json")

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})

	t.Run("When Invalid JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, constants.EndpointV1+"/products/1", nil)

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("When Product Not Found", func(t *testing.T) {
		req := requests.ProductRequest{
			Name:        "keyboard ye",
			Description: "lorem ipsum dolor sit amet",
			Price:       20.0,
			Stock:       234,
		}

		reqBody, _ := json.Marshal(req)
		productRepoMock.Mock.On("UpdateProduct", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.Anything).Return(V1Domains.ProductDomain{}, sql.ErrNoRows).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPut, constants.EndpointV1+"/products/1", bytes.NewReader(reqBody))

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestDeleteProduct(t *testing.T) {
	setupProduct(t)

	sProduct.Use(lazyAuthAdminProduct)
	sProduct.DELETE(constants.EndpointV1+"/products/:id", ProductHandler.Delete)

	t.Run("When Success", func(t *testing.T) {
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.Anything).Return(productDataFromDB, nil).Once()
		productRepoMock.Mock.On("DeleteProduct", mock.Anything, mock.Anything).Return(nil).Once()
		ristrettoProductMock.On("Del", "products", "product/product_id:1").Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, constants.EndpointV1+"/products/1", nil)

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNoContent, w.Result().StatusCode)
	})

	t.Run("When Product Not Found", func(t *testing.T) {
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.Anything).Return(V1Domains.ProductDomain{}, sql.ErrNoRows).Once()

		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodDelete, constants.EndpointV1+"/products/1", nil)

		sProduct.ServeHTTP(w, r)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}
