package handler

import "github.com/gin-gonic/gin"

func Health(c *gin.Context) {
	OK(c, gin.H{"status": "ok"})
}
