package handler

import (
	"context"
	"fmt"
	"net/http"

	activitylogService "nexusweb-market/backend/internal/modules/activitylog/service"
	"nexusweb-market/backend/internal/modules/payment/dto"
	"nexusweb-market/backend/internal/modules/payment/service"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	service service.PaymentService
	logger  activityLogger
}

type activityLogger interface {
	Log(ctx context.Context, userID string, module string, action string, description string, ipAddress string) error
}

func NewPaymentHandler(service service.PaymentService, logger activitylogService.ActivityLogService) *PaymentHandler {
	return &PaymentHandler{service: service, logger: logger}
}

func (h *PaymentHandler) GetAll(c *gin.Context) {
	payments, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payments",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payments retrieved successfully",
		"data":    payments,
	})
}

func (h *PaymentHandler) GetByID(c *gin.Context) {
	payment, err := h.service.GetByID(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "payment not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payment retrieved successfully",
		"data":    payment,
	})
}

func (h *PaymentHandler) GetByInvoiceID(c *gin.Context) {
	payments, err := h.service.GetByInvoiceID(c.Request.Context(), c.Param("invoiceId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payments by invoice",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payments retrieved successfully",
		"data":    payments,
	})
}

func (h *PaymentHandler) GetByCustomerID(c *gin.Context) {
	customerID := c.GetString("user_id")
	if customerID == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "user not authenticated",
		})
		return
	}

	payments, err := h.service.GetByCustomerID(c.Request.Context(), customerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get payments",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payments retrieved successfully",
		"data":    payments,
	})
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req dto.CreatePaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	payment, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create payment",
			"error":   err.Error(),
		})
		return
	}

	if h.logger != nil {
		userID := c.GetString("user_id")
		if userID != "" {
			_ = h.logger.Log(c.Request.Context(), userID, "PAYMENT", "CREATE", fmt.Sprintf("Payment created for invoice %s.", payment.InvoiceID), c.ClientIP())
		}
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "payment created successfully",
		"data":    payment,
	})
}

func (h *PaymentHandler) UpdateStatus(c *gin.Context) {
	var req dto.UpdatePaymentStatusRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	payment, err := h.service.UpdateStatus(c.Request.Context(), c.Param("id"), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update payment status",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "payment status updated successfully",
		"data":    payment,
	})
}
func (h *PaymentHandler) GetWhatsAppLink(c *gin.Context) {
	result, err := h.service.GetWhatsAppLink(c.Request.Context(), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get whatsapp payment link",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "whatsapp payment link retrieved successfully",
		"data":    result,
	})
}
