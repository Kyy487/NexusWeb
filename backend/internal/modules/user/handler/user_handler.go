package handler

import (
	"net/http"

	"nexusweb-market/backend/internal/modules/user/dto"
	"nexusweb-market/backend/internal/modules/user/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "profile fetched successfully",
		"data":    user,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "failed to fetch users",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "users fetched successfully",
		"data":    users,
	})
}

func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	id := c.Param("id")

	var req dto.UpdateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request body",
		})
		return
	}

	if err := h.service.UpdateUserStatus(c.Request.Context(), id, req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "user status updated successfully",
	})
}