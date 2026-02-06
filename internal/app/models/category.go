package models

import (
	"time"

	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Category information
	Name        string    `json:"name" gorm:"not null;size:50;uniqueIndex"`
	Description string    `json:"description" gorm:"size:255"`
	IsActive    bool      `json:"is_active" gorm:"default:true"`
	Products    []Product `json:"products" gorm:"foreignKey:CategoryID"`
}

// TableName specifies the table name for Category model
func (Category) TableName() string {
	return "categories"
}
