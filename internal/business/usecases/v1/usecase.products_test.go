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
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	"github.com/snykk/transaction-api/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	productRepoMock    *mocks.ProductRepository
	productUsecase     V1Domains.ProductUsecase
	productsDataFromDB []V1Domains.ProductDomain
	productDataFromDB  V1Domains.ProductDomain
)

func setupProduct(t *testing.T) {
	productRepoMock = mocks.NewProductRepository(t)
	productUsecase = V1Usecases.NewProductUsecase(productRepoMock)

	currentTime := time.Now()
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
}

func TestStoreProduct(t *testing.T) {
	setupProduct(t)

	req := requests.ProductRequest{
		Name:        "keyboard",
		Price:       20.0,
		Stock:       234,
		Description: "lorem ipsum dolor sit amet",
	}

	t.Run("When Success Store Product Data", func(t *testing.T) {
		// Mock repository untuk mengembalikan data produk yang berhasil disimpan
		productRepoMock.Mock.On("StoreProduct", mock.Anything, mock.AnythingOfType("*v1.ProductDomain")).Return(productDataFromDB, nil).Once()

		// Memanggil fungsi StoreProduct
		result, statusCode, err := productUsecase.StoreProduct(context.Background(), req.ToDomain())

		// Assertions
		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusCreated, statusCode, "Status code should be Created (201)")

		// Pastikan hasilnya sesuai dengan data yang dimock
		assert.Equal(t, productDataFromDB.Id, result.Id, "Product ID should match")
		assert.Equal(t, productDataFromDB.Name, result.Name, "Product name should match")
		assert.Equal(t, productDataFromDB.Description, result.Description, "Product description should match")
		assert.Equal(t, productDataFromDB.Price, result.Price, "Product price should match")
		assert.Equal(t, productDataFromDB.Stock, result.Stock, "Product stock should match")

		// Pastikan atribut waktu tidak bernilai nil
		assert.NotNil(t, result.CreatedAt, "CreatedAt should not be nil")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Failure to Store Product Data", func(t *testing.T) {
		// Mock repository untuk mengembalikan error saat menyimpan produk
		productRepoMock.Mock.On("StoreProduct", mock.Anything, mock.AnythingOfType("*v1.ProductDomain")).Return(V1Domains.ProductDomain{}, errors.New("create product failed")).Once()

		// Memanggil fungsi StoreProduct
		result, statusCode, err := productUsecase.StoreProduct(context.Background(), req.ToDomain())

		// Assertions untuk memastikan adanya error
		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Equal(t, 0, result.Id, "Product ID should be zero on failure")
		assert.Equal(t, "create product failed", err.Error(), "Error message should match")
	})
}

func TestGetAllProduct(t *testing.T) {
	setupProduct(t)

	t.Run("When Success Get Products Data", func(t *testing.T) {
		productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return(productsDataFromDB, nil).Once()

		result, statusCode, err := productUsecase.GetAllProducts(context.Background())

		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")
		assert.Len(t, result, len(productsDataFromDB), "The number of products in result should match mock data")

		// Validate each product
		for i, product := range result {
			t.Run(fmt.Sprintf("Check Product %d", i+1), func(t *testing.T) {
				assert.Equal(t, productsDataFromDB[i].Id, product.Id, "Product ID should match")
				assert.Equal(t, productsDataFromDB[i].Name, product.Name, "Product name should match")
				assert.Equal(t, productsDataFromDB[i].Description, product.Description, "Product description should match")
				assert.Equal(t, productsDataFromDB[i].Price, product.Price, "Product price should match")
				assert.Equal(t, productsDataFromDB[i].Stock, product.Stock, "Product stock should match")
				assert.NotNil(t, product.CreatedAt, "CreatedAt should not be nil")
				assert.NotNil(t, product.UpdatedAt, "UpdatedAt should not be nil")
			})
		}
	})

	t.Run("When Failure Get Products Data", func(t *testing.T) {
		productRepoMock.Mock.On("GetAllProducts", mock.Anything).Return([]V1Domains.ProductDomain{}, errors.New("get all products failed")).Once()

		result, statusCode, err := productUsecase.GetAllProducts(context.Background())

		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Empty(t, result, "Result should be empty on failure")
		assert.Equal(t, "get all products failed", err.Error(), "Error message should match")
	})
}

func TestGetProductById(t *testing.T) {
	setupProduct(t)

	t.Run("When Success Get Product Data", func(t *testing.T) {
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(productDataFromDB, nil).Once()

		result, statusCode, err := productUsecase.GetProductById(context.Background(), productDataFromDB.Id)

		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")
		assert.Equal(t, productDataFromDB, result, "Product data should match the mock data")
	})

	t.Run("When Failure Product doesn't exist", func(t *testing.T) {
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(V1Domains.ProductDomain{}, sql.ErrNoRows).Once()

		result, statusCode, err := productUsecase.GetProductById(context.Background(), productDataFromDB.Id)

		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusNotFound, statusCode, "Status code should be Not Found (404)")
		assert.Equal(t, V1Domains.ProductDomain{}, result, "Result should be empty on product not found")
		assert.ErrorIs(t, err, sql.ErrNoRows, "Error should be sql.ErrNoRows")
	})
}

func TestDeleteProduct(t *testing.T) {
	setupProduct(t)

	t.Run("When Success Delete Product Data", func(t *testing.T) {
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(productDataFromDB, nil).Once()
		productRepoMock.Mock.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int")).Return(nil).Once()

		statusCode, err := productUsecase.DeleteProduct(context.Background(), productDataFromDB.Id)

		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")
	})

	t.Run("When Failure Delete Product Data", func(t *testing.T) {
		t.Run("Product doesn't exist", func(t *testing.T) {
			productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(V1Domains.ProductDomain{}, errors.New("product not found")).Once()

			statusCode, err := productUsecase.DeleteProduct(context.Background(), 1)

			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusNotFound, statusCode, "Status code should be Not Found (404)")
			assert.EqualError(t, err, "product not found", "Error message should match")
		})

		t.Run("Failed Delete Product", func(t *testing.T) {
			productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(productDataFromDB, nil).Once()
			productRepoMock.Mock.On("DeleteProduct", mock.Anything, mock.AnythingOfType("int")).Return(errors.New("failed")).Once()

			statusCode, err := productUsecase.DeleteProduct(context.Background(), 1)

			assert.NotNil(t, err, "Error should not be nil")
			assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
			assert.EqualError(t, err, "failed", "Error message should match")
		})
	})
}

func TestUpdateProduct(t *testing.T) {
	setupProduct(t)

	t.Run("When Success Update Product", func(t *testing.T) {
		currentTime := time.Now()
		updatedProductFromDB := productDataFromDB
		updatedProductFromDB.UpdatedAt = &currentTime
		productRepoMock.Mock.On("UpdateProduct", mock.Anything, mock.AnythingOfType("*v1.ProductDomain")).Return(nil).Once()
		productRepoMock.Mock.On("GetProductById", mock.Anything, mock.AnythingOfType("int")).Return(updatedProductFromDB, nil).Once()

		result, statusCode, err := productUsecase.UpdateProduct(context.Background(), &productDataFromDB, productDataFromDB.Id)

		assert.Nil(t, err, "Error should be nil")
		assert.Equal(t, http.StatusOK, statusCode, "Status code should be OK (200)")
		assert.Equal(t, updatedProductFromDB, result, "Updated product should match the mock data")
		assert.NotNil(t, result.UpdatedAt, "UpdatedAt should not be nil")
	})

	t.Run("When Failure Update Product", func(t *testing.T) {
		productRepoMock.Mock.On("UpdateProduct", mock.Anything, mock.AnythingOfType("*v1.ProductDomain")).Return(errors.New("update failed")).Once()

		result, statusCode, err := productUsecase.UpdateProduct(context.Background(), &productDataFromDB, productDataFromDB.Id)

		assert.NotNil(t, err, "Error should not be nil")
		assert.Equal(t, http.StatusInternalServerError, statusCode, "Status code should be Internal Server Error (500)")
		assert.Equal(t, V1Domains.ProductDomain{}, result, "Result should be empty on update failure")
		assert.EqualError(t, err, "update failed", "Error message should match")
	})
}
