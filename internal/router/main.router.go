package router

import (
	"net/http"

	"github.com/ariboss89/coffee-morning-services/internal/middleware"
	"github.com/ariboss89/coffee-morning-services/internal/response"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	_ "github.com/ariboss89/coffee-morning-services/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Init(app *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	app.Use(middleware.CORSMiddleware)
	AuthRouter(app, db, rdb)
	UserRouter(app, db, rdb)
	app.NoRoute(func(ctx *gin.Context) {
		response.Error(ctx, http.StatusNotFound, "No route here, go to http://localhost:8002/swagger/index.html")
	})
	app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	app.Static("/static/img", "public")
}
