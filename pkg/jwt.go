package pkg

import (
	"os"
	"time"

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

func GenerateJWT(userID uuid.UUID, email string, internalRole, externalRole string, companyID, branchID *uuid.UUID) (string, error) {
	claims := JWTClaims{
		UserID:       userID,
		Email:        email,
		InternalRole: internalRole,
		ExternalRole: externalRole,
		CompanyID:    companyID,
		BranchID:     branchID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
