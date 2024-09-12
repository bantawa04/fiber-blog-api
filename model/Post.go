package model

import (
	"time"
)

type Post struct {
	ID         uint      `gorm:"primarykey"`
	Title      string    `json:"title" validate:"required,min=3,max=100"`
	Body       string    `json:"body" validate:"required,min=10"`
	CategoryID *uint     `json:"category_id" validate:"required"`
	Category   *Category `json:"category" gorm:"constraint:OnDelete:SET NULL"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
