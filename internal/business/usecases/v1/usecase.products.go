package v1

import (
	"context"
	"errors"
	"net/http"

	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
)

type productUsecase struct {
	repo V1Domains.ProductRepository
}

func NewProductUsecase(repo V1Domains.ProductRepository) V1Domains.ProductUsecase {
	return &productUsecase{
		repo,
	}
}

func (uc *productUsecase) GetAll(ctx context.Context) ([]V1Domains.ProductDomain, int, error) {
	products, err := uc.repo.GetAll(ctx)

	if err != nil {
		return []V1Domains.ProductDomain{}, http.StatusInternalServerError, err
	}

	return products, http.StatusOK, nil
}

func (uc *productUsecase) Store(ctx context.Context, product *V1Domains.ProductDomain) (V1Domains.ProductDomain, int, error) {
	result, err := uc.repo.Store(ctx, product)
	if err != nil {
		return result, http.StatusInternalServerError, err
	}
	return result, http.StatusCreated, nil
}

func (uc *productUsecase) GetById(ctx context.Context, id int) (V1Domains.ProductDomain, int, error) {
	result, err := uc.repo.GetById(ctx, id)

	if err != nil {
		return V1Domains.ProductDomain{}, http.StatusNotFound, errors.New("product not found")
	}

	return result, http.StatusOK, nil
}

func (uc *productUsecase) Update(ctx context.Context, product *V1Domains.ProductDomain, id int) (V1Domains.ProductDomain, int, error) {
	product.Id = id
	if err := uc.repo.Update(ctx, product); err != nil {
		return V1Domains.ProductDomain{}, http.StatusInternalServerError, err
	}

	newProduct, err := uc.repo.GetById(ctx, id)
	if err != nil {
		return V1Domains.ProductDomain{}, http.StatusNotFound, err
	}

	return newProduct, http.StatusOK, err
}

func (uc *productUsecase) Delete(ctx context.Context, id int) (int, error) {
	_, err := uc.repo.GetById(ctx, id)
	if err != nil { // check wheter data is exists or not
		return http.StatusNotFound, errors.New("product not found")
	}
	err = uc.repo.Delete(ctx, id)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
