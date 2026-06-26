package model

import "time"

type OrderFile struct {
	ID         string
	OrderID    string
	UploadedBy string
	FileName   string
	FileURL    string
	FileType   *string
	FileSize   int64
	CreatedAt  time.Time
}