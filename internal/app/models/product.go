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
	ImageURL    string  `json:"image_url" gorm:"type:text"` // URL to the product image
	Discount     float64 `json:"discount" gorm:"default:0"` // Discount percentage (0-100)
	DiscountedPrice float64 `json:"discounted_price" gorm:"-"` // Calculated field, not stored in DB
	IsActive    bool    `json:"is_active" gorm:"default:true"`
	StandID      uint    `json:"stand_id" gorm:"not null;index"` // Canteen stand ID
}

// TableName specifies the table name for Product model
func (Product) TableName() string {
	return "products"
}

// AfterFind hook to calculate discounted price
func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	if p.Discount > 0 {
		p.DiscountedPrice = p.Price * (1 - p.Discount/100)
	} else {
		p.DiscountedPrice = p.Price
	}
	return
}
