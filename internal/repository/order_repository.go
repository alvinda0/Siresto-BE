package repository

import (
	"project-name/internal/entity"
	"project-name/pkg"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrderRepository interface {
	Create(order *entity.Order) error
	CreateOrderItems(items []entity.OrderItem) error
	Update(order *entity.Order) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Order, error)
	FindAll(companyID, branchID *uuid.UUID, pagination pkg.PaginationParams) ([]entity.Order, int64, error)
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

func (r *orderRepository) Update(order *entity.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Order{}, "id = ?", id).Error
}

func (r *orderRepository) FindByID(id uuid.UUID) (*entity.Order, error) {
	var order entity.Order
	err := r.db.Preload("OrderItems.Product").
		Preload("Company").
		Preload("Branch").
		First(&order, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &order, nil
}

func (r *orderRepository) FindAll(companyID, branchID *uuid.UUID, pagination pkg.PaginationParams) ([]entity.Order, int64, error) {
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

	// Count total
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination
	offset := (pagination.Page - 1) * pagination.Limit
	err := query.Preload("OrderItems.Product").
		Preload("Company").
		Preload("Branch").
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
