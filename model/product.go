// Code generated by sql2gorm. DO NOT EDIT.
package model

import (
	"time"
)

type Product struct {
	ProductID     int       `gorm:"column:product_id;primary_key;AUTO_INCREMENT"`
	ProductName   string    `gorm:"column:product_name;NOT NULL"`
	Description   string    `gorm:"column:description"`
	Price         float64     `gorm:"column:price;NOT NULL"`
	StockQuantity int       `gorm:"column:stock_quantity;NOT NULL"`
	CreatedAt     time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt     time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (m *Product) TableName() string {
	return "product"
}

