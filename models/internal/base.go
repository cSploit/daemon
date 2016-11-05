package internal

import "time"

// the same as gorm.Model, but without the DeletedAt
type Base struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
