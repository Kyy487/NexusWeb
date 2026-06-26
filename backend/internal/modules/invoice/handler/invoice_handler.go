package handler

import (
	"context"
	"fmt"
	"net/http"

	activitylogService "nexusweb-market/backend/internal/modules/activitylog/service"
	"nexusweb-market/backend/internal/modules/invoice/dto"
	"nexusweb-market/backend/internal/modules/invoice/service"

	"github.com/gin-gonic/gin"
)

type InvoiceHandler struct {
	service service.InvoiceService
	logger  activityLogger
}

type activityLogger interface {
	Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error
}

func NewInvoiceHandler(service service.InvoiceService, logger activitylogService.ActivityLogService) *InvoiceHandler {
	return &InvoiceHandler{service: service, logger: logger}
}

func (h *InvoiceHandler) GetAll(c *gin.Context) {
	invoices, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get invoices",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "invoices retrieved successfully",
		"data":    invoices,
	})
}

func (h *InvoiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	invoice, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "invoice not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "invoice retrieved successfully",
		"data":    invoice,
	})
}

func (h *InvoiceHandler) GetByOrderID(c *gin.Context) {
	orderID := c.Param("orderId")

	invoice, err := h.service.GetByOrderID(c.Request.Context(), orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "invoice not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "invoice retrieved successfully",
		"data":    invoice,
	})
}

func (h *InvoiceHandler) Create(c *gin.Context) {
	var req dto.CreateInvoiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	invoice, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create invoice",
			"error":   err.Error(),
		})
		return
	}

	if h.logger != nil {
		userID := c.GetString("user_id")
		if userID != "" {
			_ = h.logger.Log(c.Request.Context(), userID, "INVOICE", "CREATE", fmt.Sprintf("Invoice %s created for order %s.", invoice.InvoiceNumber, invoice.OrderID), c.ClientIP())
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "invoice created successfully",
		"data":    invoice,
	})
}

func (h *InvoiceHandler) UpdateStatus(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateInvoiceStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	invoice, err := h.service.UpdateStatus(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update invoice status",
			"error":   err.Error(),
		})
		return
	}

	if h.logger != nil {
		userID := c.GetString("user_id")
		if userID != "" {
			_ = h.logger.Log(c.Request.Context(), userID, "INVOICE", "UPDATE_STATUS", fmt.Sprintf("Invoice %s status updated to %s.", invoice.InvoiceNumber, req.Status), c.ClientIP())
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "invoice status updated successfully",
		"data":    invoice,
	})
}
