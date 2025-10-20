package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func DeviceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Device ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			return
		}

		deviceID := strings.TrimPrefix(auth, "Device ")
		if deviceID == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "deviceID required"})
			return
		}

		c.Set("deviceID", deviceID)
		c.Next()
	}
}
