package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/activitylog/dto"
	"nexusweb-market/backend/internal/modules/activitylog/service"

	"github.com/gin-gonic/gin"
)

type ActivityLogHandler struct {
	service service.ActivityLogService
}

func NewActivityLogHandler(service service.ActivityLogService) *ActivityLogHandler {
	return &ActivityLogHandler{service: service}
}

func (h *ActivityLogHandler) Create(c *gin.Context) {
	var req dto.CreateActivityLogRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	if req.IPAddress == "" {
		req.IPAddress = c.ClientIP()
	}

	result, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create activity log",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "activity log created successfully",
		"data":    result,
	})
}

func (h *ActivityLogHandler) GetAll(c *gin.Context) {
	logs, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get activity logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "activity logs retrieved successfully",
		"data":    logs,
	})
}

func (h *ActivityLogHandler) GetByUserID(c *gin.Context) {
	logs, err := h.service.GetByUserID(c.Request.Context(), c.Param("userId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get user activity logs",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user activity logs retrieved successfully",
		"data":    logs,
	})
}