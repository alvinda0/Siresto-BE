package entity

import "github.com/google/uuid"

type CreateOrderRequest struct {
	CustomerName  string           `json:"customer_name"`
	CustomerPhone string           `json:"customer_phone"`
	TableNumber   string           `json:"table_number" binding:"required"`
	Notes         string           `json:"notes"`
	ReferralCode  string           `json:"referral_code"`
	OrderMethod   OrderMethod      `json:"order_method" binding:"required"`
	PromoCode     string           `json:"promo_code"`
	OrderItems    []OrderItemInput `json:"order_items" binding:"required,min=1"`
}

type QuickOrderRequest struct {
	TableNumber string           `json:"table_number" binding:"required"`
	OrderMethod OrderMethod      `json:"order_method" binding:"required"`
	OrderItems  []OrderItemInput `json:"order_items" binding:"required,min=1"`
}

type AddOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Note      string    `json:"note"`
}

type CreatePublicOrderRequest struct {
	CompanyID     uuid.UUID        `json:"company_id" binding:"required"`
	BranchID      uuid.UUID        `json:"branch_id" binding:"required"`
	CustomerName  string           `json:"customer_name"`
	CustomerPhone string           `json:"customer_phone"`
	TableNumber   string           `json:"table_number" binding:"required"`
	Notes         string           `json:"notes"`
	ReferralCode  string           `json:"referral_code"`
	OrderMethod   OrderMethod      `json:"order_method" binding:"required"`
	PromoCode     string           `json:"promo_code"`
	OrderItems    []OrderItemInput `json:"order_items" binding:"required,min=1"`
}

type OrderItemInput struct {
	ProductID uuid.UUID `json:"product_id" binding:"required"`
	Quantity  int       `json:"quantity" binding:"required,min=1"`
	Note      string    `json:"note"`
}

type UpdateOrderRequest struct {
	CustomerName  string           `json:"customer_name"`
	CustomerPhone string           `json:"customer_phone"`
	TableNumber   string           `json:"table_number"`
	Notes         string           `json:"notes"`
	OrderMethod   OrderMethod      `json:"order_method"`
	Status        OrderStatus      `json:"status"`
	OrderItems    []OrderItemInput `json:"order_items"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" binding:"required"`
}

type OrderResponse struct {
	ID             uuid.UUID        `json:"id"`
	CompanyID      uuid.UUID        `json:"company_id"`
	BranchID       uuid.UUID        `json:"branch_id"`
	CustomerName   string           `json:"customer_name"`
	CustomerPhone  string           `json:"customer_phone"`
	TableNumber    string           `json:"table_number"`
	Notes          string           `json:"notes"`
	ReferralCode   string           `json:"referral_code"`
	OrderMethod    OrderMethod      `json:"order_method"`
	PromoCode      string           `json:"promo_code"`       // Comma-separated promo codes
	PromoID        *uuid.UUID       `json:"promo_id"`
	DiscountAmount float64          `json:"discount_amount"`  // Total discount dari semua promo
	PromoDetails   []PromoDetailDTO `json:"promo_details,omitempty"` // Detail semua promo yang digunakan
	Status         OrderStatus      `json:"status"`
	SubtotalAmount float64          `json:"subtotal_amount"` // Total item sebelum diskon & pajak
	TaxAmount      float64          `json:"tax_amount"`      // Total pajak
	TotalAmount    float64          `json:"total_amount"`    // (Subtotal - Diskon) + Tax
	PaymentMethod  PaymentMethod    `json:"payment_method"`
	PaymentStatus  PaymentStatus    `json:"payment_status"`
	PaidAmount     float64          `json:"paid_amount"`
	ChangeAmount   float64          `json:"change_amount"`
	PaymentNote    string           `json:"payment_note"`
	PaidAt         string           `json:"paid_at"`
	TaxDetails     []TaxDetailDTO   `json:"tax_details"`     // Detail perhitungan pajak
	OrderItems     []OrderItemDTO   `json:"order_items"`
	CreatedAt      string           `json:"created_at"`
	UpdatedAt      string           `json:"updated_at"`
}

type TaxDetailDTO struct {
	TaxID      uuid.UUID `json:"tax_id"`
	TaxName    string    `json:"tax_name"`
	Percentage float64   `json:"percentage"`
	Priority   int       `json:"priority"`
	BaseAmount float64   `json:"base_amount"` // Jumlah yang dikenakan pajak
	TaxAmount  float64   `json:"tax_amount"`  // Hasil perhitungan pajak
}

type OrderItemDTO struct {
	ID          uuid.UUID    `json:"id"`
	ProductID   uuid.UUID    `json:"product_id"`
	ProductName string       `json:"product_name"`
	Quantity    int          `json:"quantity"`
	Price       float64      `json:"price"`
	Subtotal    float64      `json:"subtotal"`
	Note        string       `json:"note"`
}

type PromoDetailDTO struct {
	PromoID        uuid.UUID `json:"promo_id"`
	PromoName      string    `json:"promo_name"`
	PromoCode      string    `json:"promo_code"`
	PromoType      string    `json:"promo_type"`       // percentage atau fixed
	PromoValue     float64   `json:"promo_value"`      // nilai promo (% atau nominal)
	DiscountAmount float64   `json:"discount_amount"`  // jumlah diskon yang didapat
	MaxDiscount    *float64  `json:"max_discount"`     // maksimal diskon (untuk percentage)
	MinTransaction *float64  `json:"min_transaction"`  // minimum transaksi
}

type ProcessPaymentRequest struct {
	PaymentMethod PaymentMethod `json:"payment_method" binding:"required"`
	PaidAmount    float64       `json:"paid_amount" binding:"required,gt=0"`
	PromoCode     string        `json:"promo_code"` // Opsional: apply promo saat payment
	PaymentNote   string        `json:"payment_note"`
}

type PaymentResponse struct {
	OrderID        uuid.UUID        `json:"order_id"`
	PaymentMethod  PaymentMethod    `json:"payment_method"`
	PaymentStatus  PaymentStatus    `json:"payment_status"`
	SubtotalAmount float64          `json:"subtotal_amount"`
	DiscountAmount float64          `json:"discount_amount"`
	TaxAmount      float64          `json:"tax_amount"`
	TotalAmount    float64          `json:"total_amount"`
	PaidAmount     float64          `json:"paid_amount"`
	ChangeAmount   float64          `json:"change_amount"`
	PaymentNote    string           `json:"payment_note"`
	PaidAt         string           `json:"paid_at"`
	PromoDetails   []PromoDetailDTO `json:"promo_details,omitempty"` // Detail semua promo
	TaxDetails     []TaxDetailDTO   `json:"tax_details"`
}
