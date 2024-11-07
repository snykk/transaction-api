package v1

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	"github.com/snykk/transaction-api/internal/http/datatransfers/responses"
)

type ProductHandler struct {
	productUsecase V1Domains.ProductUsecase
	ristrettoCache caches.RistrettoCache
}

func NewProductHandler(productUsecase V1Domains.ProductUsecase, ristrettoCache caches.RistrettoCache) ProductHandler {
	return ProductHandler{
		productUsecase: productUsecase,
		ristrettoCache: ristrettoCache,
	}
}

func (c *ProductHandler) Store(ctx *gin.Context) {
	var productRequest requests.ProductRequest

	if err := ctx.ShouldBindJSON(&productRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	b, statusCode, err := c.productUsecase.Store(ctxx, productRequest.ToDomain())
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("products")

	NewSuccessResponse(ctx, statusCode, "product inserted successfully", map[string]interface{}{
		"product": responses.FromProductDomainV1(b),
	})
}

func (c *ProductHandler) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("products"); val != nil {
		NewSuccessResponse(ctx, http.StatusOK, "product data fetched successfully", map[string]interface{}{
			"products": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	listOfProducts, statusCode, err := c.productUsecase.GetAll(ctxx)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	productResponses := responses.ToProductResponseList(listOfProducts)

	if productResponses == nil {
		NewSuccessResponse(ctx, statusCode, "product data is empty", []int{})
		return
	}

	go c.ristrettoCache.Set("products", productResponses)

	NewSuccessResponse(ctx, statusCode, "product data fetched successfully", map[string]interface{}{
		"products": productResponses,
	})
}

func (c *ProductHandler) GetById(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))
	if val := c.ristrettoCache.Get(fmt.Sprintf("product/%d", id)); val != nil {
		NewSuccessResponse(ctx, http.StatusOK, fmt.Sprintf("product data with id %d fetched successfully", id), map[string]interface{}{
			"product": val,
		})
		return
	}

	ctxx := ctx.Request.Context()

	productDomain, statusCode, err := c.productUsecase.GetById(ctxx, id)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	productResponse := responses.FromProductDomainV1(productDomain)

	go c.ristrettoCache.Set(fmt.Sprintf("product/%d", id), productResponse)

	NewSuccessResponse(ctx, statusCode, fmt.Sprintf("product data with id %d fetched successfully", id), map[string]interface{}{
		"product": productResponse,
	})
}

func (c *ProductHandler) Update(ctx *gin.Context) {
	var productUpdateRequest requests.ProductRequest
	id, _ := strconv.Atoi(ctx.Param("id"))

	if err := ctx.ShouldBindJSON(&productUpdateRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctxx := ctx.Request.Context()
	productDomain := productUpdateRequest.ToDomain()
	newProduct, statusCode, err := c.productUsecase.Update(ctxx, productDomain, id)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("products", fmt.Sprintf("product/%d", id))

	NewSuccessResponse(ctx, statusCode, fmt.Sprintf("product data with id %d updated successfully", id), map[string]interface{}{
		"product": responses.FromProductDomainV1(newProduct),
	})
}

func (c *ProductHandler) Delete(ctx *gin.Context) {
	id, _ := strconv.Atoi(ctx.Param("id"))

	ctxx := ctx.Request.Context()
	statusCode, err := c.productUsecase.Delete(ctxx, id)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("products", fmt.Sprintf("product/%d", id))

	NewSuccessResponse(ctx, statusCode, fmt.Sprintf("product data with id %d deleted successfully", id), nil)
}
