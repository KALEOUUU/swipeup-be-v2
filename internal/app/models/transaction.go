package models

import (
	"time"

	"gorm.io/gorm"
)

// Transaction represents a balance transaction (top-up or deduction)
type Transaction struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Transaction information
	TransactionNumber string  `json:"transaction_number" gorm:"uniqueIndex;not null;size:50"`
	UserID            uint    `json:"user_id" gorm:"not null;index"`
	User              User    `json:"user" gorm:"foreignKey:UserID"`
	Type              string  `json:"type" gorm:"not null;size:20"` // top_up, purchase, refund
	Amount            float64 `json:"amount" gorm:"not null"`
	BalanceBefore     float64 `json:"balance_before" gorm:"not null"`
	BalanceAfter      float64 `json:"balance_after" gorm:"not null"`
	Description       string  `json:"description" gorm:"size:255"`
	OrderID           *uint   `json:"order_id" gorm:"index"` // nullable, reference to order if applicable
	Order             *Order  `json:"order" gorm:"foreignKey:OrderID"`
}

// TableName specifies the table name for Transaction model
func (Transaction) TableName() string {
	return "transactions"
}
