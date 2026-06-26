package dto

type InvoiceResponse struct {
	ID            string  `json:"id"`
	OrderID       string  `json:"order_id"`
	OrderNumber   string  `json:"order_number"`
	InvoiceNumber string  `json:"invoice_number"`
	Subtotal      float64 `json:"subtotal"`
	Discount      float64 `json:"discount"`
	Tax           float64 `json:"tax"`
	TotalAmount   float64 `json:"total_amount"`
	Status        string  `json:"status"`
	DueDate       *string `json:"due_date"`
}

type CreateInvoiceRequest struct {
	OrderID  string  `json:"order_id" binding:"required"`
	Discount float64 `json:"discount"`
	Tax      float64 `json:"tax"`
	DueDate  string  `json:"due_date"`
}

type UpdateInvoiceStatusRequest struct {
	Status string `json:"status" binding:"required"`
}