package entity

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string     `json:"name" gorm:"not null"`
	Email     string     `json:"-" gorm:"unique;not null"`
	Password  string     `json:"-" gorm:"not null"`
	RoleID    uuid.UUID  `json:"role_id" gorm:"type:uuid;not null"`
	Role      Role       `json:"role" gorm:"foreignKey:RoleID;constraint:OnDelete:RESTRICT"`
	CompanyID *uuid.UUID `json:"company_id,omitempty" gorm:"type:uuid;index"`
	Company   *Company   `json:"company,omitempty" gorm:"foreignKey:CompanyID;constraint:OnDelete:SET NULL"`
	BranchID  *uuid.UUID `json:"branch_id,omitempty" gorm:"type:uuid;index"`
	Branch    *Branch    `json:"branch,omitempty" gorm:"foreignKey:BranchID;constraint:OnDelete:SET NULL"`
	IsActive  bool       `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// IsInternalUser cek apakah user adalah internal platform
func (u *User) IsInternalUser() bool {
	return u.Role.Type == RoleTypeInternal
}

// IsExternalUser cek apakah user adalah client restoran
func (u *User) IsExternalUser() bool {
	return u.Role.Type == RoleTypeExternal
}
