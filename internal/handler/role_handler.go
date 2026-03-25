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
// External users: hanya bisa lihat external role
func (h *RoleHandler) GetAllRoles(c *gin.Context) {
	internalRole, _ := c.Get("internalRole")
	externalRole, _ := c.Get("externalRole")
	
	var roles []entity.Role
	var err error
	
	// Jika internal user, tampilkan semua role
	if internalRole != nil && internalRole != "" {
		roles, err = h.roleService.GetAllRoles()
	} else if externalRole != nil && externalRole != "" {
		// Jika external user, hanya tampilkan external role
		roles, err = h.roleService.GetRolesByType("EXTERNAL")
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

