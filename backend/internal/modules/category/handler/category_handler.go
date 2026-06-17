package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/category/dto"
	"nexusweb-market/backend/internal/modules/category/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: service}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch categories",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "categories fetched successfully",
		"data":    categories,
	})
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	category, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "category fetched successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	category, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to create category",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "category created successfully",
		"data":    category,
	})
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateCategoryRequest

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
		"message": "category updated successfully",
	})
}

func (h *CategoryHandler) Delete(c *gin.Context) {
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
		"message": "category deleted successfully",
	})
}