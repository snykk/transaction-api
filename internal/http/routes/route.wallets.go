package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1Usecase "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	V1PostgresRepository "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	V1Handler "github.com/snykk/transaction-api/internal/http/handlers/v1"
)

type walletRoutes struct {
	v1Handler       V1Handler.WalletHandler
	router          *gin.RouterGroup
	db              *sqlx.DB
	authMiddleware  gin.HandlerFunc
	adminMiddleware gin.HandlerFunc
}

func NewWalletRoute(router *gin.RouterGroup, db *sqlx.DB, ristrettoCache caches.RistrettoCache, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) *walletRoutes {
	V1WalletRepository := V1PostgresRepository.NewWalletRepository(db)
	V1WalletUsecase := V1Usecase.NewWalletUsecase(V1WalletRepository)
	V1WalletHandler := V1Handler.NewWalletHandler(V1WalletUsecase, ristrettoCache)

	return &walletRoutes{v1Handler: V1WalletHandler, router: router, db: db, authMiddleware: authMiddleware, adminMiddleware: adminMiddleware}
}

func (r *walletRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		walletRoute := V1Route.Group("/wallets")

		// authenticated user
		walletRoute.Use(r.authMiddleware)
		{
			walletRoute.POST("/init", r.v1Handler.Init)
			walletRoute.GET("/info", r.v1Handler.Info)
		}

		// admin only
		walletRoute.Use(r.adminMiddleware)
		{
			// admin only
			walletRoute.GET("", r.v1Handler.GetAll)
			// ...
		}
	}

}
