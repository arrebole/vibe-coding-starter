package entity

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Title     string         `gorm:"size:200;not null;index" json:"title"`
	Completed bool           `gorm:"not null;default:false" json:"completed"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (Todo) TableName() string {
	return "todos"
}
