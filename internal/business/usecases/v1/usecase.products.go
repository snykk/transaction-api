package v1

import (
	"context"
	"net/http"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/utils"
)

type productUsecase struct {
	repo V1Domains.ProductRepository
}

func NewProductUsecase(repo V1Domains.ProductRepository) V1Domains.ProductUsecase {
	return &productUsecase{
		repo,
	}
}

func (uc *productUsecase) GetAllProducts(ctx context.Context) ([]V1Domains.ProductDomain, int, error) {
	products, err := uc.repo.GetAllProducts(ctx)

	if err != nil {
		return []V1Domains.ProductDomain{}, http.StatusInternalServerError, err
	}

	return products, http.StatusOK, nil
}

func (uc *productUsecase) StoreProduct(ctx context.Context, product *V1Domains.ProductDomain) (V1Domains.ProductDomain, int, error) {
	result, err := uc.repo.StoreProduct(ctx, product)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}
	return result, http.StatusCreated, nil
}

func (uc *productUsecase) GetProductById(ctx context.Context, id int) (V1Domains.ProductDomain, int, error) {
	result, err := uc.repo.GetProductById(ctx, id)

	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.ProductDomain{}, statusCode, err
	}

	return result, http.StatusOK, nil
}

func (uc *productUsecase) UpdateProduct(ctx context.Context, product *V1Domains.ProductDomain, id int) (V1Domains.ProductDomain, int, error) {
	product.Id = id
	if err := uc.repo.UpdateProduct(ctx, product); err != nil {
		return V1Domains.ProductDomain{}, http.StatusInternalServerError, err
	}

	newProduct, err := uc.repo.GetProductById(ctx, id)
	if err != nil {
		statusCode, _ := utils.MapDBError(err)
		return V1Domains.ProductDomain{}, statusCode, err
	}

	return newProduct, http.StatusOK, err
}

func (uc *productUsecase) DeleteProduct(ctx context.Context, id int) (int, error) {
	_, err := uc.repo.GetProductById(ctx, id)
	if err != nil { // check wheter data is exists or not
		statusCode, _ := utils.MapDBError(err)
		return statusCode, err
	}
	err = uc.repo.DeleteProduct(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusNoContent, nil
}
