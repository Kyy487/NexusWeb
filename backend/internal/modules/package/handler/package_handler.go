package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/package/dto"
	"nexusweb-market/backend/internal/modules/package/service"

	"github.com/gin-gonic/gin"
)

type PackageHandler struct {
	service service.PackageService
}

func NewPackageHandler(service service.PackageService) *PackageHandler {
	return &PackageHandler{service: service}
}

func (h *PackageHandler) GetAll(c *gin.Context) {
	packages, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to get packages",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "packages retrieved successfully",
		"data":    packages,
	})
}

func (h *PackageHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	pkg, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "package not found",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "package retrieved successfully",
		"data":    pkg,
	})
}

func (h *PackageHandler) Create(c *gin.Context) {
	var req dto.CreatePackageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	pkg, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create package",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "package created successfully",
		"data":    pkg,
	})
}

func (h *PackageHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdatePackageRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}

	pkg, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to update package",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "package updated successfully",
		"data":    pkg,
	})
}

func (h *PackageHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	err := h.service.Delete(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to delete package",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "package deleted successfully",
	})
}