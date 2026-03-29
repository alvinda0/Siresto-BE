package entity

import (
	"time"

	"github.com/google/uuid"
)

type Promo struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	BranchID       *uuid.UUID `gorm:"type:uuid" json:"branch_id"` // nullable untuk company-level promo
	Name           string     `gorm:"type:varchar(100);not null" json:"name"`
	Code           string     `gorm:"type:varchar(50);not null" json:"code"`
	Type           string     `gorm:"type:varchar(20);not null" json:"type"` // percentage atau fixed
	Value          float64    `gorm:"type:decimal(15,2);not null" json:"value"`
	MaxDiscount    *float64   `gorm:"type:decimal(15,2)" json:"max_discount"`    // untuk percentage type
	MinTransaction *float64   `gorm:"type:decimal(15,2)" json:"min_transaction"` // minimum transaksi
	Quota          *int       `gorm:"type:int" json:"quota"`                     // jumlah maksimal penggunaan
	UsedCount      int        `gorm:"type:int;default:0" json:"used_count"`      // jumlah sudah digunakan
	StartDate      time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate        time.Time  `gorm:"type:date;not null" json:"end_date"`
	IsActive       bool       `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Branch  *Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (Promo) TableName() string {
	return "promos"
}
