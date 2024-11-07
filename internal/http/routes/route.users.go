package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	V1Usecase "github.com/snykk/transaction-api/internal/business/usecases/v1"
	"github.com/snykk/transaction-api/internal/datasources/caches"
	V1PostgresRepository "github.com/snykk/transaction-api/internal/datasources/repositories/postgres/v1"
	V1Handler "github.com/snykk/transaction-api/internal/http/handlers/v1"
	"github.com/snykk/transaction-api/pkg/jwt"
	"github.com/snykk/transaction-api/pkg/mailer"
)

type usersRoutes struct {
	v1Handler      V1Handler.UserHandler
	router         *gin.RouterGroup
	db             *sqlx.DB
	authMiddleware gin.HandlerFunc
}

func NewUsersRoute(router *gin.RouterGroup, db *sqlx.DB, jwtService jwt.JWTService, redisCache caches.RedisCache, ristrettoCache caches.RistrettoCache, authMiddleware gin.HandlerFunc, mailer mailer.OTPMailer) *usersRoutes {
	V1UserRepository := V1PostgresRepository.NewUserRepository(db)
	V1UserUsecase := V1Usecase.NewUserUsecase(V1UserRepository, jwtService, mailer)
	V1UserHandler := V1Handler.NewUserHandler(V1UserUsecase, redisCache, ristrettoCache)

	return &usersRoutes{v1Handler: V1UserHandler, router: router, db: db, authMiddleware: authMiddleware}
}

func (r *usersRoutes) Routes() {
	// Routes V1
	V1Route := r.router.Group("/v1")
	{
		// auth
		V1AuhtRoute := V1Route.Group("/auth")
		V1AuhtRoute.POST("/regis", r.v1Handler.Regis)
		V1AuhtRoute.POST("/login", r.v1Handler.Login)
		V1AuhtRoute.POST("/send-otp", r.v1Handler.SendOTP)
		V1AuhtRoute.POST("/verif-otp", r.v1Handler.VerifOTP)

		// users
		userRoute := V1Route.Group("/users")
		userRoute.Use(r.authMiddleware)
		{
			userRoute.GET("/me", r.v1Handler.GetUserData)
			// ...
		}
	}

}
