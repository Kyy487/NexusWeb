package dto

type ProgressResponse struct {
	ID                 string `json:"id"`
	OrderID            string `json:"order_id"`
	Title              string `json:"title"`
	Description        string `json:"description"`
	ProgressPercentage int    `json:"progress_percentage"`
	CreatedBy          string `json:"created_by"`
}

type CreateProgressRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ProgressPercentage int    `json:"progress_percentage" binding:"required"`
	CreatedBy  string `json:"created_by" binding:"required"`
}

type UpdateProgressRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	ProgressPercentage int    `json:"progress_percentage" binding:"required"`
}