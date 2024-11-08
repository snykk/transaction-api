package v1

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/constants"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	"github.com/snykk/transaction-api/internal/http/datatransfers/requests"
	"github.com/snykk/transaction-api/internal/http/datatransfers/responses"
	"github.com/snykk/transaction-api/pkg/jwt"
)

type TransactionHandler struct {
	transactionUsecase V1Domains.TransactionUsecase
	ristrettoCache     caches.RistrettoCache
}

func NewTransactionHandler(transactionUsecase V1Domains.TransactionUsecase, ristrettoCache caches.RistrettoCache) TransactionHandler {
	return TransactionHandler{
		transactionUsecase: transactionUsecase,
		ristrettoCache:     ristrettoCache,
	}
}

func (c *TransactionHandler) Deposit(ctx *gin.Context) {
	var walletDepositRequest requests.TransactionDepositOrWithdrawRequest

	// get authenticated user from context
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)

	if err := ctx.ShouldBindJSON(&walletDepositRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	trDom := walletDepositRequest.ToDomain()
	trDom.Wallet.UserId = userClaims.UserID

	ctxx := ctx.Request.Context()
	transactionDom, statusCode, err := c.transactionUsecase.Deposit(ctxx, trDom)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("transactions")

	NewSuccessResponse(ctx, statusCode, "product inserted successfully", map[string]interface{}{
		"product": responses.FromTransactionDomainV1(transactionDom),
	})
}

func (c *TransactionHandler) Withdraw(ctx *gin.Context) {
	var walletWithdrawRequest requests.TransactionDepositOrWithdrawRequest

	// Ambil data pengguna yang sudah diotentikasi dari konteks
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)

	// Validasi request JSON
	if err := ctx.ShouldBindJSON(&walletWithdrawRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// Konversi request menjadi domain model
	trDom := walletWithdrawRequest.ToDomain()
	trDom.Wallet.UserId = userClaims.UserID

	// Panggil usecase untuk melakukan withdraw
	ctxx := ctx.Request.Context()
	transactionDom, statusCode, err := c.transactionUsecase.Withdraw(ctxx, trDom)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	// Hapus cache yang berkaitan dengan transaksi
	go c.ristrettoCache.Del("transactions")

	// Kirim respons sukses
	NewSuccessResponse(ctx, statusCode, "withdraw completed successfully", map[string]interface{}{
		"transaction": responses.FromTransactionDomainV1(transactionDom),
	})
}

func (c *TransactionHandler) Purchase(ctx *gin.Context) {
	// 1. Bind request body ke struct TransactionPurchaseRequest
	var purchaseRequest requests.TransactionPurchaseRequest

	// 2. Get authenticated user from context
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)

	// 3. Validasi input request
	if err := ctx.ShouldBindJSON(&purchaseRequest); err != nil {
		NewErrorResponse(ctx, http.StatusBadRequest, err.Error())
		return
	}

	// 4. Mapping request ke domain dan menambahkan UserId dari authenticated user
	trDom := purchaseRequest.ToDomain()
	trDom.Wallet.UserId = userClaims.UserID

	// 5. Memanggil usecase untuk melakukan purchase
	ctxx := ctx.Request.Context()
	transactionDom, statusCode, err := c.transactionUsecase.Purchase(ctxx, trDom)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	// 6. Menghapus cache transaksi jika diperlukan (misal menggunakan ristretto)
	go c.ristrettoCache.Del("transactions")
	go c.ristrettoCache.Del("products", fmt.Sprintf("product/%d", *transactionDom.ProductId))

	// 7. Mengembalikan response sukses
	NewSuccessResponse(ctx, statusCode, "purchase successful", map[string]interface{}{
		"transaction": responses.FromTransactionDomainV1(transactionDom),
	})
}
