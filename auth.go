package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func ValidateToken() gin.HandlerFunc {
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

		valid := false

		// Do work with token and product
		log.Printf("Request for product '%s', token '%s'", product, token)
		valid = true

		if !valid {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error":   2,
				"message": "Invalid token.",
			})
			return
		}
	}
}
