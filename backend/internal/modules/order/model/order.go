package model

import "time"

type Order struct {
	ID          string
	CustomerID  string
	ServiceID   string
	PackageID   string
	OrderNumber string
	Title       string
	Description string
	Deadline    *time.Time
	TotalPrice  float64
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt *time.Time
	CancelledAt *time.Time

	CustomerName string
	ServiceName  string
	PackageName  string
}