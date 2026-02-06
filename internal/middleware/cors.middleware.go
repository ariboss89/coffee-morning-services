package middleware

import (
	"log"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware(c *gin.Context) {
	defer c.Next()
	// Allowed Origin
	whiteListOrigin := []string{"http://127.0.0.1:8001", "http://localhost:8001", "http://127.0.0.1:6380", "http://localhost:5173", "http://localhost:8002", "http://127.0.0.1:8002", "http://localhost:8003", "http://127.0.0.1:8003"}
	origin := c.GetHeader("Origin")
	if slices.Contains(whiteListOrigin, origin) {
		c.Header("Access-Control-Allow-Origin", origin)
	} else {
		log.Printf("Origin is not in the Whitelist: %s", origin)
	}
	// Allowed Header
	allowedHeaders := []string{"Authorization", "Content-Type"}
	c.Header("Access-Control-Allow-Headers", strings.Join(allowedHeaders, ", "))
	// Allowed Methods
	allowedMethod := []string{http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete, http.MethodOptions, http.MethodHead}
	c.Header("Access-Control-Allow-Methods", strings.Join(allowedMethod, ", "))
	// Handling preflight
	if c.Request.Method == http.MethodOptions {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
}
