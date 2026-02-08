package models

import (
	"time"

	"gorm.io/gorm"
)

// Order represents a purchase order
type Order struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Order information
	OrderNumber    string     `json:"order_number" gorm:"uniqueIndex;not null;size:50"`
	UserID         uint       `json:"user_id" gorm:"not null;index"`
	User           User       `json:"user" gorm:"foreignKey:UserID"`
	TotalAmount    float64    `json:"total_amount" gorm:"not null"`
	Status         string     `json:"status" gorm:"not null;size:20;default:'payment_pending'"` // payment_pending, request, cooking, done, cancelled
	PaymentMethod  string     `json:"payment_method" gorm:"size:20;default:'card'"`        // card, cash, qris
	StandID        uint       `json:"stand_id" gorm:"not null;index"` // Canteen stand ID
	Stand          User       `json:"stand" gorm:"foreignKey:StandID"` // Reference to stand admin
	OrderItems     []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`

	// Payment details
	CashAmount     float64 `json:"cash_amount,omitempty" gorm:"default:0"`     // Amount of cash provided by user (for cash payment)
	PaymentProofURL string `json:"payment_proof_url,omitempty" gorm:"type:text"` // URL to payment proof image (for QRIS payment)
}

// TableName specifies the table name for Order model
func (Order) TableName() string {
	return "orders"
}
