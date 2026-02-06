package middleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/ariboss89/coffee-morning-services/internal/dto"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func IsBlackListed(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		tokenJWT := strings.Split(auth, "Bearer ")

		token := tokenJWT[1]
		//token, _ := c.Get("tokenJWT")
		rkey := "ari:tickitz:logout" + fmt.Sprint(token)
		rsc := rdb.Get(c, rkey)

		if rsc.Err() == nil {
			tokenStore := rsc.Val()
			if token == tokenStore {
				log.Println("token is blacklisted")
				c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ResponseError{
					Message: "Unauthorized Token",
					Status:  "False",
					Error:   "Unauthorized Token",
				})
				return
			}
		}
		c.Next()
	}
}
