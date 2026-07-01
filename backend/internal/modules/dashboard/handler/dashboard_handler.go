package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/dashboard/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service service.DashboardService
}

func NewDashboardHandler(service service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: service}
}

// GetStats returns admin-wide stats for ADMIN/SUPER_ADMIN,
// or customer-specific stats for CUSTOMER role.
func (h *DashboardHandler) GetStats(c *gin.Context) {
	userID := c.GetString("user_id")
	role := c.GetString("role")

	if role == "CUSTOMER" {
		if userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"success": false,
				"message": "user not authenticated",
			})
			return
		}

		stats, err := h.service.GetCustomerStats(c.Request.Context(), userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": "failed to get dashboard stats",
				"error":   err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "customer dashboard stats retrieved successfully",
			"data":    stats,
		})
		return
	}

	// Admin / Super Admin — return global stats
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get dashboard stats",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "dashboard stats retrieved successfully",
		"data":    stats,
	})
}