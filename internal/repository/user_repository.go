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
	FindByCompanyID(companyID uuid.UUID) ([]entity.User, error)
	FindByBranchID(branchID uuid.UUID) ([]entity.User, error)
	FindInternalUsers() ([]entity.User, error)
	FindExternalUsers() ([]entity.User, error)
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
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

func (r *userRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
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

func (r *userRepository) FindByCompanyID(companyID uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("company_id = ?", companyID).Preload("Role").Preload("Branch").Find(&users).Error
	return users, err
}

func (r *userRepository) FindByBranchID(branchID uuid.UUID) ([]entity.User, error) {
	var users []entity.User
	err := r.db.Where("branch_id = ?", branchID).Find(&users).Error
	return users, err
}

func (r *userRepository) FindInternalUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.type = ?", "INTERNAL").
		Preload("Role").
		Find(&users).Error
	return users, err
}

func (r *userRepository) FindExternalUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.type = ?", "EXTERNAL").
		Preload("Role").
		Preload("Company").
		Preload("Branch").
		Find(&users).Error
	return users, err
}

func (r *userRepository) DB() *gorm.DB {
	return r.db
}
