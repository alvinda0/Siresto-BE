package entity

import (
	"time"

	"github.com/google/uuid"
)

// TransactionReportDTO untuk response report transaksi
type TransactionReportDTO struct {
	ID              uuid.UUID     `json:"id"`
	OrderNumber     string        `json:"order_number"`
	CustomerName    string        `json:"customer_name"`
	CustomerPhone   string        `json:"customer_phone"`
	TableNumber     string        `json:"table_number"`
	OrderMethod     OrderMethod   `json:"order_method"`
	Status          OrderStatus   `json:"status"`
	PaymentStatus   PaymentStatus `json:"payment_status"`
	PaymentMethod   string        `json:"payment_method"`
	SubtotalAmount  float64       `json:"subtotal_amount"`
	TaxAmount       float64       `json:"tax_amount"`
	DiscountAmount  float64       `json:"discount_amount"`
	TotalAmount     float64       `json:"total_amount"`
	PaidAmount      float64       `json:"paid_amount"`
	ChangeAmount    float64       `json:"change_amount"`
	PromoCode       string        `json:"promo_code,omitempty"`
	CompanyName     string        `json:"company_name"`
	BranchName      string        `json:"branch_name"`
	CreatedAt       time.Time     `json:"created_at"`
	PaidAt          *time.Time    `json:"paid_at,omitempty"`
}

// TransactionReportFilter untuk filter report
type TransactionReportFilter struct {
	StartDate     string `form:"start_date"`     // Format: 2006-01-02
	EndDate       string `form:"end_date"`       // Format: 2006-01-02
	StartTime     string `form:"start_time"`     // Format: 15:04
	EndTime       string `form:"end_time"`       // Format: 15:04
	Search        string `form:"search"`         // Search by customer name, phone, order number
	Status        string `form:"status"`         // Filter by order status
	PaymentStatus string `form:"payment_status"` // Filter by payment status
	PaymentMethod string `form:"payment_method"` // Filter by payment method
	OrderMethod   string `form:"order_method"`   // Filter by order method
	Page          int    `form:"page"`
	Limit         int    `form:"limit"`
}

// ToReportDTO converts Order to TransactionReportDTO
func (o *Order) ToReportDTO() TransactionReportDTO {
	dto := TransactionReportDTO{
		ID:             o.ID,
		OrderNumber:    o.ID.String()[:8], // First 8 chars of UUID as order number
		CustomerName:   o.CustomerName,
		CustomerPhone:  o.CustomerPhone,
		TableNumber:    o.TableNumber,
		OrderMethod:    o.OrderMethod,
		Status:         o.Status,
		PaymentStatus:  o.PaymentStatus,
		PaymentMethod:  string(o.PaymentMethod), // Convert PaymentMethod type to string
		SubtotalAmount: o.SubtotalAmount,
		TaxAmount:      o.TaxAmount,
		DiscountAmount: o.DiscountAmount,
		TotalAmount:    o.TotalAmount,
		PaidAmount:     o.PaidAmount,
		ChangeAmount:   o.ChangeAmount,
		PromoCode:      o.PromoCode,
		CreatedAt:      o.CreatedAt,
		PaidAt:         o.PaidAt,
	}

	// Add company and branch names if preloaded
	if o.Company.ID != uuid.Nil {
		dto.CompanyName = o.Company.Name
	}
	if o.Branch.ID != uuid.Nil {
		dto.BranchName = o.Branch.Name
	}

	return dto
}
