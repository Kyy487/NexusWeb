package model

import "time"

type Requirement struct {
	ID        string
	OrderID   string
	Question  string
	Answer    string
	CreatedAt time.Time
	UpdatedAt time.Time
}