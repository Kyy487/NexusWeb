package dto

type UserResponse struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Role   string `json:"role"`
	Status string `json:"status"`
}

type UpdateUserStatusRequest struct {
	Status string `json:"status" binding:"required"`
}