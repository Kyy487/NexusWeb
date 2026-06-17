package model

import "time"

type User struct {
	ID        string
	RoleID    string
	RoleName  string
	Name      string
	Email     string
	Phone     string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}