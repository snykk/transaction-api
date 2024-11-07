package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1Usecase "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	V1PostgresRepository "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	V1Handler "github.com/snykk/transaction-api/internal/http/handlers/v1"
)

type productRoutes struct {
	v1Handler       V1Handler.ProductHandler
	router          *gin.RouterGroup
	db              *sqlx.DB
	authMiddleware  gin.HandlerFunc
	adminMiddleware gin.HandlerFunc
}

func NewProductsRoute(router *gin.RouterGroup, db *sqlx.DB, ristrettoCache caches.RistrettoCache, authMiddleware gin.HandlerFunc, adminMiddleware gin.HandlerFunc) *productRoutes {
	V1ProductRepository := V1PostgresRepository.NewProductRepository(db)
	V1ProductUsecase := V1Usecase.NewProductUsecase(V1ProductRepository)
	V1ProductHandler := V1Handler.NewProductHandler(V1ProductUsecase, ristrettoCache)

	return &productRoutes{v1Handler: V1ProductHandler, router: router, db: db, authMiddleware: authMiddleware, adminMiddleware: adminMiddleware}
}

func (r *productRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		bookRoute := V1Route.Group("/products")

		// authenticated user
		bookRoute.Use(r.authMiddleware)
		{
			bookRoute.GET("", r.v1Handler.GetAll)
			bookRoute.GET("/:id", r.v1Handler.GetById)

			// admin only
			bookRoute.Use(r.adminMiddleware)
			{
				// admin only
				bookRoute.POST("", r.v1Handler.Store)
				bookRoute.PUT("/:id", r.v1Handler.Update)
				bookRoute.DELETE("/:id", r.v1Handler.Delete)
				// ...
			}

		}
	}

}
