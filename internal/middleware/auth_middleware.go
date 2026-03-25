package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID       uuid.UUID  `json:"user_id"`
	Email        string     `json:"email"`
	InternalRole string     `json:"internal_role,omitempty"`
	ExternalRole string     `json:"external_role,omitempty"`
	CompanyID    *uuid.UUID `json:"company_id,omitempty"`
	BranchID     *uuid.UUID `json:"branch_id,omitempty"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok {
			c.Set("userID", claims.UserID)
			c.Set("email", claims.Email)
			c.Set("internalRole", claims.InternalRole)
			c.Set("externalRole", claims.ExternalRole)
			c.Set("companyID", claims.CompanyID)
			c.Set("branchID", claims.BranchID)
		}

		c.Next()
	}
}
