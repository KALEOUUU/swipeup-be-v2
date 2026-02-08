package models

import (
	"time"

	"gorm.io/gorm"
)

// Cart represents a shopping cart for a student
type Cart struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Cart information
	UserID     uint       `json:"user_id" gorm:"not null;unique;index"` // One cart per user
	User       User       `json:"user" gorm:"foreignKey:UserID"`
	CartItems  []CartItem `json:"cart_items" gorm:"foreignKey:CartID"`
	TotalItems int        `json:"total_items" gorm:"default:0"`
	TotalPrice float64    `json:"total_price" gorm:"default:0"`
}

// CartItem represents an item in the shopping cart
type CartItem struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Cart item information
	CartID    uint    `json:"cart_id" gorm:"not null;index"`
	Cart      Cart    `json:"cart" gorm:"foreignKey:CartID"`
	ProductID uint    `json:"product_id" gorm:"not null;index"`
	Product   Product `json:"product" gorm:"foreignKey:ProductID"`
	Quantity  int     `json:"quantity" gorm:"not null;default:1"`
	Price     float64 `json:"price" gorm:"not null"`     // Price at time of adding to cart
	Subtotal  float64 `json:"subtotal" gorm:"not null"`  // Quantity * Price
	StandID   uint    `json:"stand_id" gorm:"not null"`  // For grouping by stand during checkout
}

// TableName specifies the table name for Cart model
func (Cart) TableName() string {
	return "carts"
}

// TableName specifies the table name for CartItem model
func (CartItem) TableName() string {
	return "cart_items"
}