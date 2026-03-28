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

type OrderResponse struct {
	ID            uuid.UUID       `json:"id"`
	CompanyID     uuid.UUID       `json:"company_id"`
	BranchID      uuid.UUID       `json:"branch_id"`
	CustomerName  string          `json:"customer_name"`
	CustomerPhone string          `json:"customer_phone"`
	TableNumber   string          `json:"table_number"`
	Notes         string          `json:"notes"`
	ReferralCode  string          `json:"referral_code"`
	OrderMethod   OrderMethod     `json:"order_method"`
	PromoCode     string          `json:"promo_code"`
	Status        OrderStatus     `json:"status"`
	TotalAmount   float64         `json:"total_amount"`
	OrderItems    []OrderItemDTO  `json:"order_items"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
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
