package repository

import (
	"project-name/internal/entity"

	"gorm.io/gorm"
)

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (r *RoleRepository) FindAll() ([]entity.Role, error) {
	var roles []entity.Role
	err := r.DB.Where("is_active = ?", true).Order("type, name").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindByType(roleType string) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.DB.Where("type = ? AND is_active = ?", roleType, true).Order("name").Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) FindByTypeExcluding(roleType string, excludeNames []string) ([]entity.Role, error) {
	var roles []entity.Role
	err := r.DB.Where("type = ? AND is_active = ? AND name NOT IN ?", roleType, true, excludeNames).Order("name").Find(&roles).Error
	return roles, err
}
