package dto

type ServiceResponse struct {
	ID            string  `json:"id"`
	CategoryID    string  `json:"category_id"`
	CategoryName  string  `json:"category_name"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	Description   string  `json:"description"`
	BasePrice     float64 `json:"base_price"`
	EstimatedDays int     `json:"estimated_days"`
	Status        string  `json:"status"`
}

type CreateServiceRequest struct {
	CategoryID    string  `json:"category_id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Slug          string  `json:"slug" binding:"required"`
	Description   string  `json:"description"`
	BasePrice     float64 `json:"base_price"`
	EstimatedDays int     `json:"estimated_days"`
	Status        string  `json:"status"`
}

type UpdateServiceRequest struct {
	CategoryID    string  `json:"category_id" binding:"required"`
	Name          string  `json:"name" binding:"required"`
	Slug          string  `json:"slug" binding:"required"`
	Description   string  `json:"description"`
	BasePrice     float64 `json:"base_price"`
	EstimatedDays int     `json:"estimated_days"`
	Status        string  `json:"status"`
}