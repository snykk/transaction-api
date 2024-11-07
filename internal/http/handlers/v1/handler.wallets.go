package v1

import (
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
