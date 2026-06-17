package model

import "time"

type Service struct {
	ID            string
	CategoryID    string
	CategoryName  string
	Name          string
	Slug          string
	Description   string
	BasePrice     float64
	EstimatedDays int
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}