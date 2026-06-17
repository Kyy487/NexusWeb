package model

import "time"

type Category struct {
	ID          string
	Name        string
	Slug        string
	Description string
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}