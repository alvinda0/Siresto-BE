package handler

import (
	"project-name/internal/entity"
	"project-name/internal/service"
	"project-name/pkg"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleService *service.RoleService
}

func NewRoleHandler(roleService *service.RoleService) *RoleHandler {
	return &RoleHandler{roleService: roleService}
}

// GetAllRoles untuk mendapatkan semua roles
// Internal users: bisa lihat semua role
// External users: filter berdasarkan role mereka
//   - OWNER: bisa lihat semua external role (OWNER, ADMIN, CASHIER, KITCHEN, WAITER)
//   - ADMIN: hanya bisa lihat role di bawahnya (ADMIN, CASHIER, KITCHEN, WAITER)
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	internalRole, _ := c.Get("internalRole")
	externalRole, _ := c.Get("externalRole")
	
	var roles []entity.Role
	var err error
	
	// Jika internal user, tampilkan semua role
	if internalRole != nil && internalRole != "" {
		roles, err = h.roleService.GetAllRoles()
	} else if externalRole != nil && externalRole != "" {
		// Jika external user, filter berdasarkan role
		externalRoleStr := externalRole.(string)
		
		if externalRoleStr == "OWNER" {
			// OWNER bisa lihat semua external role
			roles, err = h.roleService.GetRolesByType("EXTERNAL")
		} else if externalRoleStr == "ADMIN" {
			// ADMIN hanya bisa lihat role di bawahnya (tidak termasuk OWNER)
			roles, err = h.roleService.GetRolesExcluding("EXTERNAL", []string{"OWNER"})
		} else {
			// Role lain (CASHIER, KITCHEN, WAITER) tidak boleh akses
			pkg.ErrorResponse(c, 403, "Access denied", "You don't have permission to view roles")
			return
		}
	} else {
		pkg.ErrorResponse(c, 403, "Access denied", "Invalid user type")
		return
	}
	
	if err != nil {
		pkg.ErrorResponse(c, 500, "Failed to retrieve roles", err.Error())
		return
	}

	pkg.SuccessResponse(c, 200, "Roles retrieved successfully", roles)
}

