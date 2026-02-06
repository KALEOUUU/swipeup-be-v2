package models

import (
	"time"

	"gorm.io/gorm"
)

// Product represents a canteen product/item
type Product struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"index"`

	// Product information
	Name        string  `json:"name" gorm:"not null;size:100"`
	Description string  `json:"description" gorm:"size:255"`
	CategoryID  uint    `json:"category_id" gorm:"not null"`
	Category    Category `json:"category" gorm:"foreignKey:CategoryID"`
	Price       float64 `json:"price" gorm:"not null"`
	Stock       int     `json:"stock" gorm:"default:0"`
	ImageURL    string  `json:"image_url" gorm:"size:255"`
	ImageBase64 string  `json:"image_base64" gorm:"type:text"` // Base64 encoded image
	QRISBase64  string  `json:"qris_base64" gorm:"type:text"` // Base64 encoded QRIS photo
	Discount     float64 `json:"discount" gorm:"default:0"` // Discount percentage (0-100)
	DiscountedPrice float64 `json:"discounted_price" gorm:"-"` // Calculated field, not stored in DB
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	StandID      uint    `json:"stand_id" gorm:"not null;index"` // Canteen stand ID
	Stand        User    `json:"stand" gorm:"foreignKey:StandID"` // Reference to stand admin
	Status       string  `json:"status" gorm:"not null;size:20;default:'request'"` // payment_pending, request, cooking, done
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}
