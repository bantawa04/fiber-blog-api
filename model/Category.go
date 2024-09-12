package model

import (
	"time"
)

type Category struct {
	ID        uint   `gorm:"primarykey"`
	Name      string `json:"name" validate:"required,min=3,max=100"`
	Posts     []Post `json:"posts" gorm:"foreignKey:CategoryID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
