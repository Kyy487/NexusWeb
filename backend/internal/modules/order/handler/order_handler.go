package handler

import (
	"context"
	"fmt"
	"net/http"

	activitylogService "nexusweb-market/backend/internal/modules/activitylog/service"
	"nexusweb-market/backend/internal/modules/order/dto"
	"nexusweb-market/backend/internal/modules/order/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.OrderService
	logger  activityLogger
}

type activityLogger interface {
	Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error
}

func NewOrderHandler(service service.OrderService, logger activitylogService.ActivityLogService) *OrderHandler {
	return &OrderHandler{service: service, logger: logger}
}

func (h *OrderHandler) GetAll(c *gin.Context) {
	orders, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get orders",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "orders retrieved successfully",
		"data":    orders,
	})
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	order, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "order not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order retrieved successfully",
		"data":    order,
	})
}

func (h *OrderHandler) Create(c *gin.Context) {
	var req dto.CreateOrderRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	order, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create order",
			"error":   err.Error(),
		})
		return
	}

	if h.logger != nil {
		userID := c.GetString("user_id")
		if userID == "" {
			userID = req.CustomerID
		}
		_ = h.logger.Log(c.Request.Context(), userID, "ORDER", "CREATE", fmt.Sprintf("Order %s created successfully.", order.OrderNumber), c.ClientIP())
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "order created successfully",
		"data":    order,
	})
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateOrderStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	order, err := h.service.UpdateStatus(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update order status",
			"error":   err.Error(),
		})
		return
	}

	if h.logger != nil {
		userID := c.GetString("user_id")
		if userID != "" {
			_ = h.logger.Log(c.Request.Context(), userID, "ORDER", "UPDATE_STATUS", fmt.Sprintf("Order %s status updated to %s.", order.OrderNumber, req.Status), c.ClientIP())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order status updated successfully",
		"data":    order,
	})
}
