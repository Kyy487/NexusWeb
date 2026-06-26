package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	RoleSuperAdmin = "SUPER_ADMIN"
	RoleAdmin      = "ADMIN"
	RoleCustomer   = "CUSTOMER"
)

func RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")

		if role == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": "role is required",
			})
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "you do not have permission to access this resource",
		})
		c.Abort()
	}
}

func AdminOnly() gin.HandlerFunc {
	return RequireRoles(RoleSuperAdmin, RoleAdmin)
}

func SuperAdminOnly() gin.HandlerFunc {
	return RequireRoles(RoleSuperAdmin)
}

func CustomerOnly() gin.HandlerFunc {
	return RequireRoles(RoleCustomer)
}

func AdminOrCustomer() gin.HandlerFunc {
	return RequireRoles(RoleSuperAdmin, RoleAdmin, RoleCustomer)
}
