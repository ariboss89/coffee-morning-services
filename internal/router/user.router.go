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

func UserRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRouter := app.Group("/user")

	userRepository := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(userService)
	userRouter.Use(middleware.CORSMiddleware)
	userRouter.Use(middleware.VerifyJWT)
	userRouter.Use(middleware.IsBlackListed(rdb))
	userRouter.GET("/profile", userController.GetUserProfileById)
	userRouter.PATCH("/", userController.UpdateProfile)
}
