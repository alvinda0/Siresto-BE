package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// RequireInternalRole middleware untuk memastikan hanya internal users yang bisa akses
func RequireInternalRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		internalRole, exists := c.Get("internalRole")
		
		// Jika tidak ada internalRole atau kosong, berarti bukan internal user
		if !exists || internalRole == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. This endpoint is only accessible by internal users.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireExternalRole middleware untuk memastikan hanya external users yang bisa akses
func RequireExternalRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		externalRole, exists := c.Get("externalRole")
		
		// Jika tidak ada externalRole atau kosong, berarti bukan external user
		if !exists || externalRole == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "Access denied. This endpoint is only accessible by external users.",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireRoles middleware untuk memastikan user memiliki salah satu role yang diizinkan
func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		internalRole, _ := c.Get("internalRole")
		externalRole, _ := c.Get("externalRole")
		
		userRole := ""
		if internalRole != "" {
			userRole = internalRole.(string)
		} else if externalRole != "" {
			userRole = externalRole.(string)
		}

		// Cek apakah user role ada di allowed roles
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"error": "Access denied. You don't have permission to access this resource.",
		})
		c.Abort()
	}
}
