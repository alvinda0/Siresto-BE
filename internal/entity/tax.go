package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tax struct {
	ID          uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID   uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	BranchID    *uuid.UUID `gorm:"type:uuid" json:"branch_id"` // nullable untuk company-level tax
	NamaPajak   string     `gorm:"type:varchar(100);not null" json:"nama_pajak"`
	TipePajak   string     `gorm:"type:varchar(10);not null" json:"tipe_pajak"` // sc atau pb1
	Presentase  float64    `gorm:"type:decimal(5,2);not null" json:"presentase"`
	Deskripsi   string     `gorm:"type:text" json:"deskripsi"`
	Status      string     `gorm:"type:varchar(20);default:'active'" json:"status"` // active/inactive
	Prioritas   int        `gorm:"type:int;default:0" json:"prioritas"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	
	// Relations
	Company Company `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Branch  *Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}

func (Tax) TableName() string {
	return "taxes"
}
