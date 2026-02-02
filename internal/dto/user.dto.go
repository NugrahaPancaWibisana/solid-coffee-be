package dto

import "time"

type User struct {
	ID        int        `json:"id" example:"1"`
	Fullname  string     `json:"fullname" example:"John Doe"`
	Email     string     `json:"email" example:"user@example.com"`
	Photo     string     `json:"photo" example:"avatar.png"`
	Phone     string     `json:"phone" example:"081234567890"`
	Address   string     `json:"address" example:"Jakarta"`
	Role      string     `json:"role,omitempty" example:"user"`
	LastLogin *time.Time `json:"last_login,omitempty" example:"2025-01-01T10:00:00Z"`
	CreatedAt time.Time  `json:"created_at" example:"2025-01-01T09:00:00Z"`
}
