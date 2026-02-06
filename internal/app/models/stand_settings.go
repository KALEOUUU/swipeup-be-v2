package models

import (
	"time"

	"gorm.io/gorm"
)

// StandSettings represents canteen stand settings
type StandSettings struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Stand settings
	StandID   uint   `json:"stand_id" gorm:"not null;uniqueIndex"`
	Stand     User   `json:"stand" gorm:"foreignKey:StandID"`
	StoreName  string `json:"store_name" gorm:"not null;size:100"` // Canteen stand/store name
	QRIS       string `json:"qris" gorm:"size:255"` // QRIS code for payment
	IsActive   bool   `json:"is_active" gorm:"default:true"`
}

// TableName specifies the table name for StandSettings model
func (StandSettings) TableName() string {
	return "stand_settings"
}
