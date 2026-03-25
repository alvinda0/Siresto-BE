package entity

import (
	"time"

	"github.com/google/uuid"
)

type RoleType string

const (
	RoleTypeInternal RoleType = "INTERNAL"
	RoleTypeExternal RoleType = "EXTERNAL"
)

type Role struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"unique;not null"` // SUPER_ADMIN, OWNER, CASHIER, etc
	DisplayName string    `json:"display_name" gorm:"not null"`
	Type        RoleType  `json:"type" gorm:"not null"` // INTERNAL or EXTERNAL
	Description string    `json:"description"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
