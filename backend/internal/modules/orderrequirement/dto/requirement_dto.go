package dto

type RequirementResponse struct {
	ID       string `json:"id"`
	OrderID  string `json:"order_id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type CreateRequirementRequest struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}

type UpdateRequirementRequest struct {
	Question string `json:"question" binding:"required"`
	Answer   string `json:"answer" binding:"required"`
}