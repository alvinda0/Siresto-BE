package entity

import (
	"time"

	"github.com/google/uuid"
)

type Category struct {
	ID             uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID      uuid.UUID   `json:"company_id" gorm:"type:uuid;not null;index"`
	Company        *Company    `json:"company,omitempty" gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE"`
	BranchID       *uuid.UUID  `json:"branch_id,omitempty" gorm:"type:uuid;index"`
	Branch         *Branch     `json:"branch,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE"`
	Name           string      `json:"name" gorm:"not null"`
	Description    string      `json:"description"`
	Position       int         `json:"position" gorm:"not null;default:1"`
	IsActive       bool        `json:"is_active" gorm:"default:true"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
