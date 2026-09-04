package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AccessToken(token string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if token == "" {
			c.Next()
			return
		}
		if c.Request.URL.Path == "/api/health" {
			c.Next()
			return
		}
		got := c.GetHeader("X-Access-Token")
		if got == "" {
			auth := c.GetHeader("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				got = strings.TrimPrefix(auth, "Bearer ")
			}
		}
		if got != token {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未授权",
				"data":    nil,
			})
			return
		}
		c.Next()
	}
}
