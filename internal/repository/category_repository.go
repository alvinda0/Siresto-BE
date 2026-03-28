package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	Create(category *entity.Category) error
	Update(category *entity.Category) error
	Delete(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.Category, error)
	FindByCompanyID(companyID uuid.UUID, branchID *uuid.UUID) ([]entity.Category, error)
	GetMaxPosition(companyID uuid.UUID, branchID *uuid.UUID) (int, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(category *entity.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *entity.Category) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&entity.Category{}, id).Error
}

func (r *categoryRepository) FindByID(id uuid.UUID) (*entity.Category, error) {
	var category entity.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) FindByCompanyID(companyID uuid.UUID, branchID *uuid.UUID) ([]entity.Category, error) {
	var categories []entity.Category
	query := r.db.Where("company_id = ?", companyID)
	
	// Filter by branch
	if branchID == nil {
		query = query.Where("branch_id IS NULL")
	} else {
		query = query.Where("branch_id = ?", *branchID)
	}
	
	err := query.Order("position ASC").Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) GetMaxPosition(companyID uuid.UUID, branchID *uuid.UUID) (int, error) {
	var maxPosition int
	query := r.db.Model(&entity.Category{}).Where("company_id = ?", companyID)
	
	// Filter by branch
	if branchID == nil {
		query = query.Where("branch_id IS NULL")
	} else {
		query = query.Where("branch_id = ?", *branchID)
	}
	
	err := query.Select("COALESCE(MAX(position), 0)").Scan(&maxPosition).Error
	return maxPosition, err
}
