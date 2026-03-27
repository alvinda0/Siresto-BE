package service

import (
	"project-name/internal/entity"
	"project-name/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetAllRoles() ([]entity.Role, error) {
	return s.roleRepo.FindAll()
}

func (s *RoleService) GetRolesByType(roleType string) ([]entity.Role, error) {
	return s.roleRepo.FindByType(roleType)
}

func (s *RoleService) GetRolesExcluding(roleType string, excludeNames []string) ([]entity.Role, error) {
	return s.roleRepo.FindByTypeExcluding(roleType, excludeNames)
}
