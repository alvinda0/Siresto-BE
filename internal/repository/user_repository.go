package repository

import (
	"project-name/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uuid.UUID) (*entity.User, error)
	FindByEmail(email string) (*entity.User, error)
	FindAll() ([]entity.User, error)
	FindByCompanyID(companyID uuid.UUID, limit, offset int) ([]entity.User, int64, error)
	FindByBranchID(branchID uuid.UUID, limit, offset int) ([]entity.User, int64, error)
	FindInternalUsers(limit, offset int) ([]entity.User, int64, error)
	FindExternalUsers(limit, offset int) ([]entity.User, int64, error)
	Update(user *entity.User) error
	Delete(id uuid.UUID) error
	DB() *gorm.DB
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).
		Preload("Role").
		Preload("Company.Owner.Role").
		Preload("Company.Owner").
		Preload("Company").
		Preload("Branch").
		First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).
		Preload("Role").
		Preload("Company").
		Preload("Branch").
		First(&user).Error
	return &user, err
}

func (r *userRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	return users, err
}

func (r *userRepository) Update(user *entity.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&entity.User{}).Error
}

func (r *userRepository) FindByCompanyID(companyID uuid.UUID, limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	
	query := r.db.Where("company_id = ?", companyID)
	
	// Count total
	query.Model(&entity.User{}).Count(&total)
	
	// Get paginated data
	err := query.Preload("Role").Preload("Branch").
		Limit(limit).Offset(offset).
		Find(&users).Error
	
	return users, total, err
}

func (r *userRepository) FindByBranchID(branchID uuid.UUID, limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	
	query := r.db.Where("branch_id = ?", branchID)
	
	// Count total
	query.Model(&entity.User{}).Count(&total)
	
	// Get paginated data
	err := query.Preload("Role").Preload("Branch").
		Limit(limit).Offset(offset).
		Find(&users).Error
	
	return users, total, err
}

func (r *userRepository) FindInternalUsers(limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	
	query := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.type = ?", "INTERNAL")
	
	// Count total
	query.Model(&entity.User{}).Count(&total)
	
	// Get paginated data
	err := query.Preload("Role").
		Limit(limit).Offset(offset).
		Find(&users).Error
	
	return users, total, err
}

func (r *userRepository) FindExternalUsers(limit, offset int) ([]entity.User, int64, error) {
	var users []entity.User
	var total int64
	
	query := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.type = ?", "EXTERNAL")
	
	// Count total
	query.Model(&entity.User{}).Count(&total)
	
	// Get paginated data
	err := query.Preload("Role").
		Preload("Company").
		Preload("Branch").
		Limit(limit).Offset(offset).
		Find(&users).Error
	
	return users, total, err
}

func (r *userRepository) DB() *gorm.DB {
	return r.db
}
