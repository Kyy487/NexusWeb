package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/orderrequirement/dto"
	"nexusweb-market/backend/internal/modules/orderrequirement/service"

	"github.com/gin-gonic/gin"
)

type RequirementHandler struct {
	service service.RequirementService
}

func NewRequirementHandler(service service.RequirementService) *RequirementHandler {
	return &RequirementHandler{service: service}
}

func (h *RequirementHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")

	requirements, err := h.service.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get order requirements",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order requirements retrieved successfully",
		"data":    requirements,
	})
}

func (h *RequirementHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	requirement, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "requirement not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "requirement retrieved successfully",
		"data":    requirement,
	})
}

func (h *RequirementHandler) Create(c *gin.Context) {
	orderID := c.Param("orderId")

	var req dto.CreateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	requirement, err := h.service.Create(c.Request.Context(), orderID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create requirement",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "requirement created successfully",
		"data":    requirement,
	})
}

func (h *RequirementHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateRequirementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	requirement, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update requirement",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "requirement updated successfully",
		"data":    requirement,
	})
}

func (h *RequirementHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete requirement",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "requirement deleted successfully",
	})
}