package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1Usecase "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	V1PostgresRepository "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	V1Handler "github.com/snykk/transaction-api/internal/http/handlers/v1"
)

type transactionRoutes struct {
	v1Handler       V1Handler.TransactionHandler
	router          *gin.RouterGroup
	db              *sqlx.DB
	authMiddleware  gin.HandlerFunc
	adminMiddleware gin.HandlerFunc
}

func NewTransactionRoute(router *gin.RouterGroup, db *sqlx.DB, ristrettoCache caches.RistrettoCache, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) *transactionRoutes {
	V1TransactionRepository := V1PostgresRepository.NewTransactionRepository(db)

	V1TransactionUsecase := V1Usecase.NewTransactionUsecase(V1TransactionRepository)
	V1TransactionHandler := V1Handler.NewTransactionHandler(V1TransactionUsecase, ristrettoCache)

	return &transactionRoutes{v1Handler: V1TransactionHandler, router: router, db: db, authMiddleware: authMiddleware, adminMiddleware: adminMiddleware}
}

func (r *transactionRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		transactionRoute := V1Route.Group("/transaction")

		// authenticated user
		transactionRoute.Use(r.authMiddleware)
		{
			transactionRoute.GET("/history", r.v1Handler.History)

			transactionRoute.POST("/deposit", r.v1Handler.Deposit)
			transactionRoute.POST("/withdraw", r.v1Handler.Withdraw)
			transactionRoute.POST("/purchase", r.v1Handler.Purchase)
		}

		// admin only
		transactionRoute.Use(r.adminMiddleware)
		{
			// admin only
			transactionRoute.GET("", r.v1Handler.GetAll)
			// ...
		}
	}

}
