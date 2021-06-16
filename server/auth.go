package server

import (
	"net/http"
	"strings"

	"github.com/frixuu/bearpush"
	"github.com/gin-gonic/gin"
)

// ValidateToken serves as a Gin middleware for checking
// whether a user can trigger an artifact process.
func ValidateToken(appCtx *bearpush.Context) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header["Authorization"]
		if authHeader == nil || len(authHeader) != 1 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   1,
				"message": "Token not provided.",
			})
			return
		}

		token := strings.TrimPrefix(authHeader[0], "Bearer ")
		product := c.Param("product")
		p, ok := appCtx.Products[product]
		if !ok {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   2,
				"message": "One or more requested resources is not available.",
			})
			return
		}

		if !p.VerifyToken(token) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   3,
				"message": "You are not allowed to access one or more requested resources.",
			})
			return
		}
	}
}
