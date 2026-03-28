package middleware

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// WebSocketAuthMiddleware authenticates WebSocket connections using token from query parameter
func WebSocketAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get token from query parameter
		tokenString := c.Query("token")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token required in query parameter"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("role_type", claims.RoleType)
			c.Set("internal_role", claims.InternalRole)
			c.Set("external_role", claims.ExternalRole)
			
			// Set company_id dan branch_id untuk WebSocket room
			if claims.CompanyID != nil {
				c.Set("company_id", *claims.CompanyID)
			}
			if claims.BranchID != nil {
				c.Set("branch_id", *claims.BranchID)
			}
		}

		c.Next()
	}
}
