package model

import "time"

type Progress struct {
	ID                 string
	OrderID            string
	Title              string
	Description        string
	ProgressPercentage int
	CreatedBy          string
	CreatedAt          time.Time
}