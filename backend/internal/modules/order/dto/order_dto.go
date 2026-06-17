package dto

type OrderResponse struct {
	ID           string  `json:"id"`
	CustomerID   string  `json:"customer_id"`
	CustomerName string  `json:"customer_name"`
	ServiceID    string  `json:"service_id"`
	ServiceName  string  `json:"service_name"`
	PackageID    string  `json:"package_id"`
	PackageName  string  `json:"package_name"`
	OrderNumber  string  `json:"order_number"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Deadline     *string `json:"deadline"`
	TotalPrice   float64 `json:"total_price"`
	Status       string  `json:"status"`
}

type CreateOrderRequest struct {
	CustomerID  string `json:"customer_id" binding:"required"`
	ServiceID   string `json:"service_id" binding:"required"`
	PackageID   string `json:"package_id" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Deadline    string `json:"deadline"`
}

type UpdateOrderStatusRequest struct {
	Status string `json:"status" binding:"required"`
}