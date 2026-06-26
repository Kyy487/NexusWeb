package model

import "time"

type Invoice struct {
	ID            string
	OrderID       string
	OrderNumber   string
	InvoiceNumber string
	Subtotal      float64
	Discount      float64
	Tax           float64
	TotalAmount   float64
	Status        string
	DueDate       *time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}