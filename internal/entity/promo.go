package entity

import (
	"time"

	"github.com/google/uuid"
)

type Promo struct {
	ID             uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CompanyID      uuid.UUID  `gorm:"type:uuid;not null" json:"company_id"`
	BranchID       *uuid.UUID `gorm:"type:uuid" json:"branch_id"` // nullable untuk company-level promo
	Name           string     `gorm:"type:varchar(100);not null" json:"name"`
	Code           string     `gorm:"type:varchar(50);not null" json:"code"`
	PromoCategory  string     `gorm:"type:varchar(20);not null;default:'normal'" json:"promo_category"` // normal, product, bundle
	Type           string     `gorm:"type:varchar(20);not null" json:"type"`                            // percentage atau fixed
	Value          float64    `gorm:"type:decimal(15,2);not null" json:"value"`
	MaxDiscount    *float64   `gorm:"type:decimal(15,2)" json:"max_discount"`    // untuk percentage type
	MinTransaction *float64   `gorm:"type:decimal(15,2)" json:"min_transaction"` // minimum transaksi
	Quota          *int       `gorm:"type:int" json:"quota"`                     // jumlah maksimal penggunaan
	UsedCount      int        `gorm:"type:int;default:0" json:"used_count"`      // jumlah sudah digunakan
	StartDate      time.Time  `gorm:"type:date;not null" json:"start_date"`
	EndDate        time.Time  `gorm:"type:date;not null" json:"end_date"`
	IsActive       bool       `gorm:"type:boolean;default:true" json:"is_active"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`

	// Relations
	Company       Company        `gorm:"foreignKey:CompanyID" json:"company,omitempty"`
	Branch        *Branch        `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	PromoProducts []PromoProduct `gorm:"foreignKey:PromoID" json:"promo_products,omitempty"`
	PromoBundles  []PromoBundle  `gorm:"foreignKey:PromoID" json:"promo_bundles,omitempty"`
}

// PromoProduct untuk promo yang berlaku untuk produk tertentu
type PromoProduct struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PromoID   uuid.UUID `gorm:"type:uuid;not null" json:"promo_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Promo   Promo   `gorm:"foreignKey:PromoID" json:"promo,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

// PromoBundle untuk promo bundle (beli produk A + B dapat diskon)
type PromoBundle struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	PromoID   uuid.UUID `gorm:"type:uuid;not null" json:"promo_id"`
	ProductID uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity  int       `gorm:"type:int;not null;default:1" json:"quantity"` // jumlah produk yang harus dibeli
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Promo   Promo   `gorm:"foreignKey:PromoID" json:"promo,omitempty"`
	Product Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
}

func (PromoProduct) TableName() string {
	return "promo_products"
}

func (PromoBundle) TableName() string {
	return "promo_bundles"
}

func (Promo) TableName() string {
	return "promos"
}
