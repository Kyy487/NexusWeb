package model

import "time"

type ActivityLog struct {
	ID          string    `db:"id"`
	UserID      string    `db:"user_id"`
	Module      string    `db:"module"`
	Action      string    `db:"action"`
	Description string    `db:"description"`
	IPAddress   *string   `db:"ip_address"`
	CreatedAt   time.Time `db:"created_at"`
}