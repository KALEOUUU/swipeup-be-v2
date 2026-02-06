package models

import (
	"time"

	"gorm.io/gorm"
)

// GlobalSettings represents global system settings
type GlobalSettings struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Global settings
	Key        string  `json:"key" gorm:"uniqueIndex;not null;size:50"`
	Value      string  `json:"value" gorm:"size:255"`
	IsActive   bool    `json:"is_active" gorm:"default:true"`
}

// TableName specifies the table name for GlobalSettings model
func (GlobalSettings) TableName() string {
	return "global_settings"
}
