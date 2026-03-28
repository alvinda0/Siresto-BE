package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// APILogListDTO untuk response GET all logs (tanpa response_body)
type APILogListDTO struct {
	ID            uuid.UUID  `json:"id"`
	Method        string     `json:"method"`
	Path          string     `json:"path"`
	StatusCode    int        `json:"status_code"`
	ResponseTime  int64      `json:"response_time"` // in milliseconds
	IPAddress     string     `json:"ip_address"`
	UserAgent     string     `json:"user_agent"`
	AccessFrom    string     `json:"access_from"`
	UserID        *uuid.UUID `json:"user_id,omitempty"`
	CompanyID     *uuid.UUID `json:"company_id,omitempty"`
	BranchID      *uuid.UUID `json:"branch_id,omitempty"`
	RequestBody   string     `json:"request_body,omitempty"`
	ErrorMessage  string     `json:"error_message,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
}

// APILogDetailDTO untuk response GET log by ID (dengan response_body)
type APILogDetailDTO struct {
	ID            uuid.UUID      `json:"id"`
	Method        string         `json:"method"`
	Path          string         `json:"path"`
	StatusCode    int            `json:"status_code"`
	ResponseTime  int64          `json:"response_time"` // in milliseconds
	IPAddress     string         `json:"ip_address"`
	UserAgent     string         `json:"user_agent"`
	AccessFrom    string         `json:"access_from"`
	UserID        *uuid.UUID     `json:"user_id,omitempty"`
	CompanyID     *uuid.UUID     `json:"company_id,omitempty"`
	BranchID      *uuid.UUID     `json:"branch_id,omitempty"`
	RequestBody   string         `json:"request_body,omitempty"`
	ResponseBody  string         `json:"response_body,omitempty"` // Hanya ada di detail
	ErrorMessage  string         `json:"error_message,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at,omitempty"`
}

// ToListDTO converts APILog to APILogListDTO
func (log *APILog) ToListDTO() APILogListDTO {
	return APILogListDTO{
		ID:           log.ID,
		Method:       log.Method,
		Path:         log.Path,
		StatusCode:   log.StatusCode,
		ResponseTime: log.ResponseTime,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		AccessFrom:   log.AccessFrom,
		UserID:       log.UserID,
		CompanyID:    log.CompanyID,
		BranchID:     log.BranchID,
		RequestBody:  log.RequestBody,
		ErrorMessage: log.ErrorMessage,
		CreatedAt:    log.CreatedAt,
	}
}

// ToDetailDTO converts APILog to APILogDetailDTO
func (log *APILog) ToDetailDTO() APILogDetailDTO {
	return APILogDetailDTO{
		ID:           log.ID,
		Method:       log.Method,
		Path:         log.Path,
		StatusCode:   log.StatusCode,
		ResponseTime: log.ResponseTime,
		IPAddress:    log.IPAddress,
		UserAgent:    log.UserAgent,
		AccessFrom:   log.AccessFrom,
		UserID:       log.UserID,
		CompanyID:    log.CompanyID,
		BranchID:     log.BranchID,
		RequestBody:  log.RequestBody,
		ResponseBody: log.ResponseBody,
		ErrorMessage: log.ErrorMessage,
		CreatedAt:    log.CreatedAt,
		UpdatedAt:    log.UpdatedAt,
		DeletedAt:    log.DeletedAt,
	}
}
