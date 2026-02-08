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

func InteractionRouter(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	interactionRouter := app.Group("/interaction")

	interactionRepository := repository.NewInteractionRepository()
	interactionService := service.NewInteractionService(interactionRepository, db, rdb)
	interactionController := controller.NewInteractionController(interactionService)
	interactionRouter.Use(middleware.CORSMiddleware)
	interactionRouter.Use(middleware.VerifyJWT)
	interactionRouter.Use(middleware.IsBlackListed(rdb))
	interactionRouter.POST("/content", interactionController.PostContent)
	interactionRouter.POST("/following", interactionController.FollowingUser)
	interactionRouter.POST("/like", interactionController.LikePosts)
}
