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

func (h *DashboardHandler) GetStats(c *gin.Context) {
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