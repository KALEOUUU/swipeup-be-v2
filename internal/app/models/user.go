package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a system user (student, teacher, or admin)
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// User information
	StudentId  string  `json:"student_id" gorm:"size:50;default:''"`
	Name       string  `json:"name" gorm:"not null;size:100"`
	Email      string  `json:"email" gorm:"size:100"`
	Phone      string  `json:"phone" gorm:"size:20"`
	Role       string  `json:"role" gorm:"not null;size:20;default:'student'"` // student, teacher, admin, stand
	Class      string  `json:"class" gorm:"size:50"`
	Balance    float64 `json:"balance" gorm:"default:0"`
	IsActive   bool    `json:"is_active" gorm:"default:true"`
	RFIDCard   string  `json:"rfid_card" gorm:"column:rf_id_card;size:50"`
	Password   string  `json:"-" gorm:"column:password;size:255"` // hashed password, never expose in JSON
}

// TableName specifies the table name for User model
func (User) TableName() string {
	return "users"
}
