package model

import "time"

type Payment struct {
	ID              string
	InvoiceID       string
	Amount          float64
	PaymentMethod   *string
	PaymentStatus   string
	PaymentProofURL *string
	PaidAt          *time.Time
	VerifiedBy      *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}