package entity

import (
	"time"

	"github.com/google/uuid"
)

type Branch struct {
	ID          uuid.UUID `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	CompanyID   uuid.UUID `json:"company_id" gorm:"type:uuid;not null"`
	Company     Company   `json:"company,omitempty" gorm:"foreignKey:CompanyID"`
	Name        string    `json:"name" gorm:"not null"`
	Address     string    `json:"address" gorm:"not null"`
	City        string    `json:"city"`
	Province    string    `json:"province"`
	PostalCode  string    `json:"postal_code"`
	Phone       string    `json:"phone"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
