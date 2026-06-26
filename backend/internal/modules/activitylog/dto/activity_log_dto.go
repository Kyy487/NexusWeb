package dto

type CreateActivityLogRequest struct {
	UserID      string `json:"user_id"`
	Module      string `json:"module"`
	Action      string `json:"action"`
	Description string `json:"description"`
	IPAddress   string `json:"ip_address"`
}

type ActivityLogResponse struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Module      string `json:"module"`
	Action      string `json:"action"`
	Description string `json:"description"`
	IPAddress   *string `json:"ip_address"`
	CreatedAt   string `json:"created_at"`
}