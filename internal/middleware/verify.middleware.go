package middleware

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/ariboss89/coffee-morning-services/internal/response"
	pkg "github.com/ariboss89/coffee-morning-services/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func VerifyJWT(c *gin.Context) {
	bearerToken := c.GetHeader("Authorization")
	result := strings.Split(bearerToken, " ")
	if result[0] != "Bearer" {
		log.Println("token is not bearer token")
		response.Error(c, http.StatusUnauthorized, "Unauthorized Access")
		return
	}

	var jc pkg.JWTClaims

	_, err := jc.VerifyToken(result[1])
	if err != nil {
		log.Println(err.Error())
		if errors.Is(err, jwt.ErrTokenExpired) {
			response.Error(c, http.StatusUnauthorized, "Unauthorized Access")
			return
		}
		if errors.Is(err, jwt.ErrTokenInvalidIssuer) {
			response.Error(c, http.StatusUnauthorized, "Unauthorized Access")
			return
		}
		response.Error(c, http.StatusInternalServerError, "Internal Server Error")
		return
	}
	c.Set("token", jc)
	c.Set("tokenJWT", result[1])
	c.Next()
}
