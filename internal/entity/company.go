package entity

import (
	"time"

	"github.com/google/uuid"
)

type CompanyType string

const (
	CompanyTypePT         CompanyType = "PT"
	CompanyTypeCV         CompanyType = "CV"
	CompanyTypePerorangan CompanyType = "PERORANGAN"
)

type Company struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name        string      `json:"name" gorm:"not null"`
	Type        CompanyType `json:"type" gorm:"not null"`
	OwnerID     uuid.UUID   `json:"owner_id" gorm:"type:uuid;not null"`
	Owner       User        `json:"owner" gorm:"foreignKey:OwnerID;constraint:OnDelete:CASCADE"`
	Branches    []Branch    `json:"branches,omitempty" gorm:"foreignKey:CompanyID"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}
