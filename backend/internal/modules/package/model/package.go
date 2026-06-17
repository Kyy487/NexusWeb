package model

import "time"

type Package struct {
	ID            string
	ServiceID     string
	ServiceName   string
	Name          string
	Description   string
	Price         float64
	RevisionCount int
	DeliveryDays  int
	Features      []string
	Status        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}