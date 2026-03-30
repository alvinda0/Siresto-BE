package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderMethod string

const (
	OrderMethodDineIn   OrderMethod = "DINE_IN"
	OrderMethodTakeAway OrderMethod = "TAKE_AWAY"
	OrderMethodDelivery OrderMethod = "DELIVERY"
)

type OrderStatus string

const (
	OrderStatusPending    OrderStatus = "PENDING"
	OrderStatusConfirmed  OrderStatus = "CONFIRMED"
	OrderStatusProcessing OrderStatus = "PROCESSING"
	OrderStatusReady      OrderStatus = "READY"
	OrderStatusCompleted  OrderStatus = "COMPLETED"
	OrderStatusCancelled  OrderStatus = "CANCELLED"
)

type PaymentMethod string

const (
	PaymentMethodQRIS          PaymentMethod = "QRIS"
	PaymentMethodCash          PaymentMethod = "TUNAI"
	PaymentMethodDebit         PaymentMethod = "DEBIT"
	PaymentMethodCredit        PaymentMethod = "CREDIT"
	PaymentMethodGoPay         PaymentMethod = "GOPAY"
	PaymentMethodOVO           PaymentMethod = "OVO"
	PaymentMethodComplimentary PaymentMethod = "COMPLIMENTARY"
)

type PaymentStatus string

const (
	PaymentStatusUnpaid  PaymentStatus = "UNPAID"
	PaymentStatusPaid    PaymentStatus = "PAID"
	PaymentStatusPartial PaymentStatus = "PARTIAL"
)

type Order struct {
	ID             uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID      `gorm:"type:uuid;not null;index" json:"company_id"`
	BranchID       uuid.UUID      `gorm:"type:uuid;not null;index" json:"branch_id"`
	CustomerName   string         `gorm:"type:varchar(255)" json:"customer_name"`
	CustomerPhone  string         `gorm:"type:varchar(50)" json:"customer_phone"`
	TableNumber    string         `gorm:"type:varchar(50)" json:"table_number"`
	Notes          string         `gorm:"type:text" json:"notes"`
	ReferralCode   string         `gorm:"type:varchar(100)" json:"referral_code"`
	OrderMethod    OrderMethod    `gorm:"type:varchar(50);not null" json:"order_method"`
	PromoID        *uuid.UUID     `gorm:"type:uuid" json:"promo_id"`                           // ID promo yang digunakan
	PromoCode      string         `gorm:"type:varchar(100)" json:"promo_code"`                 // Kode promo
	DiscountAmount float64        `gorm:"type:decimal(15,2);default:0" json:"discount_amount"` // Jumlah diskon dari promo
	Status         OrderStatus    `gorm:"type:varchar(50);default:'PENDING'" json:"status"`
	SubtotalAmount float64        `gorm:"type:decimal(15,2);default:0" json:"subtotal_amount"` // Total item sebelum diskon & pajak
	TaxAmount      float64        `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`      // Total pajak
	TotalAmount    float64        `gorm:"type:decimal(15,2);default:0" json:"total_amount"`    // (Subtotal - Diskon) + Tax
	PaymentMethod  PaymentMethod  `gorm:"type:varchar(50)" json:"payment_method"`              // Metode pembayaran
	PaymentStatus  PaymentStatus  `gorm:"type:varchar(50);default:'UNPAID'" json:"payment_status"` // Status pembayaran
	PaidAmount     float64        `gorm:"type:decimal(15,2);default:0" json:"paid_amount"`     // Jumlah yang sudah dibayar
	ChangeAmount   float64        `gorm:"type:decimal(15,2);default:0" json:"change_amount"`   // Kembalian (untuk TUNAI)
	PaymentNote    string         `gorm:"type:text" json:"payment_note"`                       // Catatan pembayaran
	PaidAt         *time.Time     `json:"paid_at"`                                             // Waktu pembayaran
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Company    Company      `gorm:"foreignKey:CompanyID;constraint:OnDelete:CASCADE" json:"company,omitempty"`
	Branch     Branch       `gorm:"foreignKey:BranchID;constraint:OnDelete:CASCADE" json:"branch,omitempty"`
	Promo      *Promo       `gorm:"foreignKey:PromoID" json:"promo,omitempty"`
	OrderItems []OrderItem  `gorm:"foreignKey:OrderID" json:"order_items,omitempty"`
}

type OrderItem struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	OrderID   uuid.UUID      `gorm:"type:uuid;not null;index" json:"order_id"`
	ProductID uuid.UUID      `gorm:"type:uuid;not null;index" json:"product_id"`
	Quantity  int            `gorm:"not null" json:"quantity"`
	Price     float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Note      string         `gorm:"type:text" json:"note"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	Order   Order   `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE" json:"-"`
	Product Product `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE" json:"product,omitempty"`
}
