package models

import (
	"time"

	"gorm.io/gorm"
)

// OrderItem represents an item in an order
type OrderItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// OrderItem information
	OrderID   uint    `json:"order_id" gorm:"not null;index"`
	Order     Order   `json:"order" gorm:"foreignKey:OrderID"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
	Subtotal  float64 `json:"subtotal" gorm:"not null"`
}

// TableName specifies the table name for OrderItem model
func (OrderItem) TableName() string {
	return "order_items"
}
