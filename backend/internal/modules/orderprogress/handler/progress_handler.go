package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/orderprogress/dto"
	"nexusweb-market/backend/internal/modules/orderprogress/service"

	"github.com/gin-gonic/gin"
)

type ProgressHandler struct {
	service service.ProgressService
}

func NewProgressHandler(service service.ProgressService) *ProgressHandler {
	return &ProgressHandler{service: service}
}

func (h *ProgressHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")

	progressList, err := h.service.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get order progress",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order progress retrieved successfully",
		"data":    progressList,
	})
}

func (h *ProgressHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	progress, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "progress not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "progress retrieved successfully",
		"data":    progress,
	})
}

func (h *ProgressHandler) Create(c *gin.Context) {
	orderID := c.Param("orderId")

	var req dto.CreateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	progress, err := h.service.Create(c.Request.Context(), orderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create progress",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "progress created successfully",
		"data":    progress,
	})
}

func (h *ProgressHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateProgressRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	progress, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update progress",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "progress updated successfully",
		"data":    progress,
	})
}

func (h *ProgressHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete progress",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "progress deleted successfully",
	})
}