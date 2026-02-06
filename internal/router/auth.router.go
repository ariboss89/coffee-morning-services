package router

import (
	"github.com/ariboss89/coffee-morning-services/internal/controller"
	"github.com/ariboss89/coffee-morning-services/internal/middleware"
	"github.com/ariboss89/coffee-morning-services/internal/repository"
	"github.com/ariboss89/coffee-morning-services/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func AuthRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	authRouter := app.Group("/auth")

	authRepository := repository.NewAuthRepository()
	authService := service.NewAuthService(authRepository, rdb, db)
	authController := controller.NewAuthController(authService)

	authRouter.POST("/register", authController.Register)
	authRouter.POST("/login", authController.Login)

	authRouter.Use(middleware.CORSMiddleware)
	authRouter.Use(middleware.IsBlackListed(rdb))
	// rbacMiddleware := middleware.AuthRole("user", "admin")

	// authRouter.PATCH("/password", rbacMiddleware, authController.UpdatePassword)
	authRouter.DELETE("/logout", authController.Logout)
}
