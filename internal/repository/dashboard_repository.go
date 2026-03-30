package repository

import (
	"project-name/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetTotalItemsByDate(companyID, branchID uuid.UUID, days int) ([]entity.DailyStats, error)
	GetRevenueByDate(companyID, branchID uuid.UUID, days int) ([]entity.DailyStats, error)
	GetBestSellingItems(companyID, branchID uuid.UUID, startDate, endDate time.Time, limit int) ([]entity.BestSellingItem, error)
	GetComplimentaryItems(companyID, branchID uuid.UUID, startDate, endDate time.Time) ([]entity.ComplimentaryItemStats, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db: db}
}

// GetTotalItemsByDate menghitung total item terjual per tanggal (N hari terakhir)
func (r *dashboardRepository) GetTotalItemsByDate(companyID, branchID uuid.UUID, days int) ([]entity.DailyStats, error) {
	var results []struct {
		Date  string
		Value float64
	}
	
	// Hitung tanggal mulai
	startDate := time.Now().AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	
	query := r.db.Model(&entity.OrderItem{}).
		Select("TO_CHAR(orders.created_at, 'YYYY-MM-DD') as date, COALESCE(SUM(order_items.quantity), 0) as value").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.company_id = ?", companyID).
		Where("orders.payment_status = ?", entity.PaymentStatusPaid).
		Where("DATE(orders.created_at) >= ?", startDate)
	
	if branchID != uuid.Nil {
		query = query.Where("orders.branch_id = ?", branchID)
	}
	
	err := query.Group("TO_CHAR(orders.created_at, 'YYYY-MM-DD')").Order("date ASC").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	// Convert to DailyStats
	stats := make([]entity.DailyStats, len(results))
	for i, r := range results {
		stats[i] = entity.DailyStats{
			Date:  r.Date,
			Value: r.Value,
		}
	}
	
	return stats, nil
}

// GetRevenueByDate menghitung total pendapatan per tanggal (N hari terakhir)
func (r *dashboardRepository) GetRevenueByDate(companyID, branchID uuid.UUID, days int) ([]entity.DailyStats, error) {
	var results []struct {
		Date  string
		Value float64
	}
	
	// Hitung tanggal mulai
	startDate := time.Now().AddDate(0, 0, -(days - 1)).Format("2006-01-02")
	
	query := r.db.Model(&entity.Order{}).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, COALESCE(SUM(total_amount), 0) as value").
		Where("company_id = ?", companyID).
		Where("payment_status = ?", entity.PaymentStatusPaid).
		Where("DATE(created_at) >= ?", startDate)
	
	if branchID != uuid.Nil {
		query = query.Where("branch_id = ?", branchID)
	}
	
	err := query.Group("TO_CHAR(created_at, 'YYYY-MM-DD')").Order("date ASC").Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	// Convert to DailyStats
	stats := make([]entity.DailyStats, len(results))
	for i, r := range results {
		stats[i] = entity.DailyStats{
			Date:  r.Date,
			Value: r.Value,
		}
	}
	
	return stats, nil
}

// GetBestSellingItems mendapatkan item terlaris dalam periode tertentu
func (r *dashboardRepository) GetBestSellingItems(companyID, branchID uuid.UUID, startDate, endDate time.Time, limit int) ([]entity.BestSellingItem, error) {
	var items []entity.BestSellingItem
	
	query := r.db.Model(&entity.OrderItem{}).
		Select("order_items.product_id, products.name as product_name, SUM(order_items.quantity) as total_qty, SUM(order_items.quantity * order_items.price) as total_amount").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("orders.company_id = ?", companyID).
		Where("orders.payment_status = ?", entity.PaymentStatusPaid).
		Where("orders.created_at BETWEEN ? AND ?", startDate, endDate)
	
	if branchID != uuid.Nil {
		query = query.Where("orders.branch_id = ?", branchID)
	}
	
	err := query.
		Group("order_items.product_id, products.name").
		Order("total_qty DESC").
		Limit(limit).
		Scan(&items).Error
	
	return items, err
}

// GetComplimentaryItems mendapatkan item yang dibayar dengan metode complimentary
func (r *dashboardRepository) GetComplimentaryItems(companyID, branchID uuid.UUID, startDate, endDate time.Time) ([]entity.ComplimentaryItemStats, error) {
	var items []entity.ComplimentaryItemStats
	
	query := r.db.Model(&entity.OrderItem{}).
		Select("order_items.product_id, products.name as product_name, SUM(order_items.quantity) as total_qty").
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Joins("JOIN products ON products.id = order_items.product_id").
		Where("orders.company_id = ?", companyID).
		Where("orders.payment_method = ?", entity.PaymentMethodComplimentary).
		Where("orders.payment_status = ?", entity.PaymentStatusPaid).
		Where("orders.created_at BETWEEN ? AND ?", startDate, endDate)
	
	if branchID != uuid.Nil {
		query = query.Where("orders.branch_id = ?", branchID)
	}
	
	err := query.
		Group("order_items.product_id, products.name").
		Order("total_qty DESC").
		Scan(&items).Error
	
	return items, err
}
