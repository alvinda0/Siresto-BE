package entity

import "github.com/google/uuid"

type CreatePromoRequest struct {
	BranchID       *uuid.UUID `json:"branch_id"` // optional, null = company level
	Name           string     `json:"name" binding:"required"`
	Code           string     `json:"code" binding:"required"`
	Type           string     `json:"type" binding:"required,oneof=percentage fixed"`
	Value          float64    `json:"value" binding:"required,min=0"`
	MaxDiscount    *float64   `json:"max_discount" binding:"omitempty,min=0"`
	MinTransaction *float64   `json:"min_transaction" binding:"omitempty,min=0"`
	Quota          *int       `json:"quota" binding:"omitempty,min=1"`
	StartDate      string     `json:"start_date" binding:"required"` // format: YYYY-MM-DD
	EndDate        string     `json:"end_date" binding:"required"`   // format: YYYY-MM-DD
	IsActive       *bool      `json:"is_active"`
}

type UpdatePromoRequest struct {
	Name           string   `json:"name"`
	Code           string   `json:"code"`
	Type           string   `json:"type" binding:"omitempty,oneof=percentage fixed"`
	Value          float64  `json:"value" binding:"omitempty,min=0"`
	MaxDiscount    *float64 `json:"max_discount" binding:"omitempty,min=0"`
	MinTransaction *float64 `json:"min_transaction" binding:"omitempty,min=0"`
	Quota          *int     `json:"quota" binding:"omitempty,min=1"`
	StartDate      string   `json:"start_date"` // format: YYYY-MM-DD
	EndDate        string   `json:"end_date"`   // format: YYYY-MM-DD
	IsActive       *bool    `json:"is_active"`
}

type PromoResponse struct {
	ID             uuid.UUID  `json:"id"`
	CompanyID      uuid.UUID  `json:"company_id"`
	CompanyName    string     `json:"company_name"`
	BranchID       *uuid.UUID `json:"branch_id"`
	BranchName     *string    `json:"branch_name"`
	Name           string     `json:"name"`
	Code           string     `json:"code"`
	Type           string     `json:"type"`
	Value          float64    `json:"value"`
	MaxDiscount    *float64   `json:"max_discount"`
	MinTransaction *float64   `json:"min_transaction"`
	Quota          *int       `json:"quota"`
	UsedCount      int        `json:"used_count"`
	RemainingQuota *int       `json:"remaining_quota"` // calculated field
	StartDate      string     `json:"start_date"`
	EndDate        string     `json:"end_date"`
	IsActive       bool       `json:"is_active"`
	IsExpired      bool       `json:"is_expired"`   // calculated field
	IsAvailable    bool       `json:"is_available"` // calculated field
	CreatedAt      string     `json:"created_at"`
	UpdatedAt      string     `json:"updated_at"`
}
