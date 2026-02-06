package main

import (
	"log"

	"github.com/ariboss89/coffee-morning-services/internal/config"
	"github.com/ariboss89/coffee-morning-services/internal/router"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title           Coffee Morning Services
// @version         1.0
// @description     Backend Services For Coffee Morning
// @host      			localhost:8002
// @BasePath  			/

// @securityDefinitions.apikey	BearerAuth
// @in													header
// @name 												Authorization
// @description 								Type "Bearer" followed by space and JWT Token
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to Load env")
		return
	}
	db, err := config.InitDb()
	rdb := config.InitRedis()
	defer rdb.Close()

	if err != nil {
		log.Println("Failed to Connect to Database")
		return
	}

	app := gin.Default()
	router.Init(app, db, rdb)
	app.Run(":8002")
}
