package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	user := os.Getenv("RDS_USER")
	pwd := os.Getenv("RDS_PASS")
	host := os.Getenv("RDS_HOST")
	port := os.Getenv("RDS_PORT")
	db, _ := strconv.Atoi(os.Getenv("RDS_DB"))

	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Username: user,
		Password: pwd,
		DB:       db,
	})
}
