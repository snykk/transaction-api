package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/constants"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	"github.com/snykk/transaction-api/internal/http/datatransfers/responses"
	"github.com/snykk/transaction-api/pkg/jwt"
)

type WalletHandler struct {
	walletUsecase  V1Domains.WalletUsecase
	ristrettoCache caches.RistrettoCache
}

func NewWalletHandler(walletUsecase V1Domains.WalletUsecase, ristrettoCache caches.RistrettoCache) WalletHandler {
	return WalletHandler{
		walletUsecase:  walletUsecase,
		ristrettoCache: ristrettoCache,
	}
}

func (c *WalletHandler) GetAll(ctx *gin.Context) {
	if val := c.ristrettoCache.Get("wallets"); val != nil {
		NewSuccessResponse(ctx, http.StatusOK, "wallet data fetched successfully", map[string]interface{}{
			"wallets": val,
		})
		return
	}

	ctxx := ctx.Request.Context()
	listOfWalletDom, statusCode, err := c.walletUsecase.GetAll(ctxx)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	walletResponseList := responses.ToWalletResponseList(listOfWalletDom)

	if walletResponseList == nil {
		NewSuccessResponse(ctx, statusCode, "wallet data is empty", []int{})
		return
	}

	go c.ristrettoCache.Set("wallets", walletResponseList)

	NewSuccessResponse(ctx, statusCode, "wallet data fetched successfully", map[string]interface{}{
		"wallets": walletResponseList,
	})
}

func (c *WalletHandler) Init(ctx *gin.Context) {
	// get authenticated user from context
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)

	ctxx := ctx.Request.Context()
	walletDom, statusCode, err := c.walletUsecase.Init(ctxx, userClaims.UserID)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	go c.ristrettoCache.Del("wallets")

	NewSuccessResponse(ctx, statusCode, "wallet created successfully", map[string]interface{}{
		"wallet": responses.FromWalletDomainV1(walletDom),
	})
}

func (c *WalletHandler) Info(ctx *gin.Context) {
	// get authenticated user from context
	userClaims := ctx.MustGet(constants.CtxAuthenticatedUserKey).(jwt.JwtCustomClaim)

	ctxx := ctx.Request.Context()
	walletDom, statusCode, err := c.walletUsecase.GetByUserId(ctxx, userClaims.UserID)
	if err != nil {
		NewErrorResponse(ctx, statusCode, err.Error())
		return
	}

	// remove user relation
	walletDom.User = V1Domains.UserDomain{}

	NewSuccessResponse(ctx, statusCode, "wallet fetched successfully", map[string]interface{}{
		"wallet": responses.FromWalletDomainV1(walletDom),
	})
}
