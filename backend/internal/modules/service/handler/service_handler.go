package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/service/dto"
	"nexusweb-market/backend/internal/modules/service/service"

	"github.com/gin-gonic/gin"
)

type ServiceHandler struct {
	service service.ServiceService
}

func NewServiceHandler(service service.ServiceService) *ServiceHandler {
	return &ServiceHandler{service: service}
}

func (h *ServiceHandler) GetAll(c *gin.Context) {
	services, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch services",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "services fetched successfully",
		"data":    services,
	})
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	service, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "service fetched successfully",
		"data":    service,
	})
}

func (h *ServiceHandler) Create(c *gin.Context) {
	var req dto.CreateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	service, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create service",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "service created successfully",
		"data":    service,
	})
}

func (h *ServiceHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateServiceRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	if err := h.service.Update(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "service updated successfully",
	})
}

func (h *ServiceHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "service deleted successfully",
	})
}