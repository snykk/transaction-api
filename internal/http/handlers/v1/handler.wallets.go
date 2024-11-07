package v1

import (
	V1Domains "github.com/snykk/transaction-api/internal/business/domains/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
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
