package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Product struct {
	ID          uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	BranchID    uuid.UUID      `gorm:"type:uuid;not null;index" json:"branch_id"`
	CategoryID  uuid.UUID      `gorm:"type:uuid;not null;index" json:"category_id"`
	Image       string         `gorm:"type:text" json:"image"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Stock       int            `gorm:"default:0" json:"stock"`
	Price       float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Position    string         `gorm:"type:varchar(100)" json:"position"`
	IsAvailable bool           `gorm:"default:true" json:"is_available"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Company  Company  `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"company,omitempty"`
	Branch   Branch   `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"branch,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:CASCADE" json:"category,omitempty"`
}
