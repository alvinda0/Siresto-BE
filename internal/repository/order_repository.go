package repository

import (
	"project-name/internal/entity"
	"project-name/pkg"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	CreateOrderItems(items []entity.OrderItem) error
	UpdateOrderItem(item *entity.OrderItem) error
	Update(order *entity.Order) error
	UpdateStatus(id uuid.UUID, status entity.OrderStatus) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Order, error)
	FindAll(companyID, branchID *uuid.UUID, status, method, customer, orderID string, pagination pkg.PaginationParams) ([]entity.Order, int64, error)
	DeleteOrderItems(orderID uuid.UUID) error
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Create(order *entity.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItems(items []entity.OrderItem) error {
	if len(items) == 0 {
		return nil
	}
	return r.db.Create(&items).Error
}

func (r *orderRepository) UpdateOrderItem(item *entity.OrderItem) error {
	return r.db.Model(item).Updates(map[string]interface{}{
		"quantity": item.Quantity,
		"note":     item.Note,
	}).Error
}

func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Model(order).Updates(map[string]interface{}{
		"customer_name":   order.CustomerName,
		"customer_phone":  order.CustomerPhone,
		"table_number":    order.TableNumber,
		"notes":           order.Notes,
		"order_method":    order.OrderMethod,
		"promo_id":        order.PromoID,
		"promo_code":      order.PromoCode,
		"discount_amount": order.DiscountAmount,
		"status":          order.Status,
		"subtotal_amount": order.SubtotalAmount,
		"tax_amount":      order.TaxAmount,
		"total_amount":    order.TotalAmount,
		"payment_method":  order.PaymentMethod,
		"payment_status":  order.PaymentStatus,
		"paid_amount":     order.PaidAmount,
		"change_amount":   order.ChangeAmount,
		"payment_note":    order.PaymentNote,
		"paid_at":         order.PaidAt,
	}).Error
}

func (r *orderRepository) UpdateStatus(id uuid.UUID, status entity.OrderStatus) error {
	return r.db.Model(&entity.Order{}).Where("id = ?", id).Update("status", status).Error
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Order{}, "id = ?", id).Error
}

func (r *orderRepository) FindByID(id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Preload("OrderItems.Product").
		Preload("Company").
		Preload("Branch").
		Preload("Promo").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAll(companyID, branchID *uuid.UUID, status, method, customer, orderID string, pagination pkg.PaginationParams) ([]entity.Order, int64, error) {
	var orders []entity.Order
	var total int64

	query := r.db.Model(&entity.Order{})

	// Filter by company
	if companyID != nil {
		query = query.Where("company_id = ?", *companyID)
	}

	// Filter by branch
	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}

	// Filter by status
	if status != "" {
		query = query.Where("status = ?", status)
	}

	// Filter by order method
	if method != "" {
		query = query.Where("order_method = ?", method)
	}

	// Search by customer name (case-insensitive, partial match)
	if customer != "" {
		query = query.Where("LOWER(customer_name) LIKE ?", "%"+strings.ToLower(customer)+"%")
	}

	// Search by order ID (partial match)
	if orderID != "" {
		query = query.Where("CAST(id AS TEXT) LIKE ?", "%"+orderID+"%")
	}

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Preload("OrderItems.Product").
		Preload("Company").
		Preload("Branch").
		Preload("Promo").
		Order("created_at DESC").
		Limit(pagination.Limit).
		Offset(offset).
		Find(&orders).Error

	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *orderRepository) DeleteOrderItems(orderID uuid.UUID) error {
	return r.db.Where("order_id = ?", orderID).Delete(&entity.OrderItem{}).Error
}
