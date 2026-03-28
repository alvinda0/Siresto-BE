package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type APILog struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Method        string         `gorm:"type:varchar(10);not null" json:"method"`
	Path          string         `gorm:"type:varchar(255);not null" json:"path"`
	StatusCode    int            `gorm:"not null" json:"status_code"`
	ResponseTime  int64          `gorm:"not null" json:"response_time"` // in milliseconds
	IPAddress     string         `gorm:"type:varchar(45)" json:"ip_address"`
	UserAgent     string         `gorm:"type:text" json:"user_agent"`
	AccessFrom    string         `gorm:"type:varchar(50)" json:"access_from"` // website, mobile, postman, etc
	UserID        *uuid.UUID     `gorm:"type:uuid;index" json:"user_id,omitempty"`
	CompanyID     *uuid.UUID     `gorm:"type:uuid;index" json:"company_id,omitempty"`
	BranchID      *uuid.UUID     `gorm:"type:uuid;index" json:"branch_id,omitempty"`
	RequestBody   string         `gorm:"type:text" json:"request_body,omitempty"`
	ResponseBody  string         `gorm:"type:text" json:"response_body,omitempty"`
	ErrorMessage  string         `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
