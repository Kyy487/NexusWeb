package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"nexusweb-market/backend/internal/modules/orderfile/service"

	"github.com/gin-gonic/gin"
)

type OrderFileHandler struct {
	service service.OrderFileService
}

func NewOrderFileHandler(service service.OrderFileService) *OrderFileHandler {
	return &OrderFileHandler{service: service}
}

func (h *OrderFileHandler) Upload(c *gin.Context) {
	orderID := c.PostForm("order_id")
	uploadedBy := c.PostForm("uploaded_by")
	fileType := c.PostForm("file_type")

	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": "invalid multipart form",
		"error":   err.Error(),
	})
	return
}
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "file is required",
			"error":   err.Error(),
		})
		return
	}

	uploadDir := "uploads/order-files"

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create upload directory",
			"error":   err.Error(),
		})
		return
	}

		fileName := fmt.Sprintf("%d-%s", time.Now().UnixNano(), filepath.Base(file.Filename))
		filePath := filepath.Join(uploadDir, fileName)
		fileURL := "uploads/order-files/" + fileName

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to save uploaded file",
			"error":   err.Error(),
		})
		return
	}

	result, err := h.service.SaveFile(
		c.Request.Context(),
		orderID,
		uploadedBy,
		fileType,
		file,
		fileURL,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to save file data",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "file uploaded successfully",
		"data":    result,
	})
}

func (h *OrderFileHandler) GetByOrderID(c *gin.Context) {
	files, err := h.service.GetByOrderID(c.Request.Context(), c.Param("orderId"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get order files",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "order files retrieved successfully",
		"data":    files,
	})
}

func (h *OrderFileHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c.Request.Context(), c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete file",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "file deleted successfully",
	})
}