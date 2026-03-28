package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *entity.Product) error
	FindAll(companyID, branchID uuid.UUID, search string, page, limit int) ([]entity.Product, int64, error)
	FindByID(id, companyID, branchID uuid.UUID) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(id, companyID, branchID uuid.UUID) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) Create(product *entity.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll(companyID, branchID uuid.UUID, search string, page, limit int) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.db.Model(&entity.Product{}).
		Preload("Company", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Branch", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Where("company_id = ? AND branch_id = ?", companyID, branchID)

	if search != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	if err := query.Offset(offset).Limit(limit).Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(id, companyID, branchID uuid.UUID) (*entity.Product, error) {
	var product entity.Product
	err := r.db.Preload("Company", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Branch", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Preload("Category", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Where("id = ? AND company_id = ? AND branch_id = ?", id, companyID, branchID).
		First(&product).Error
	return &product, err
}

func (r *productRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id, companyID, branchID uuid.UUID) error {
	return r.db.Where("id = ? AND company_id = ? AND branch_id = ?", id, companyID, branchID).
		Delete(&entity.Product{}).Error
}
