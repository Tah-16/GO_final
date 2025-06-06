package model

import (
	"time"
)

type Cart struct {
	CartID     int       `gorm:"column:cart_id;primary_key;AUTO_INCREMENT"`
	CustomerID int       `gorm:"column:customer_id;NOT NULL"`
	CartName   string    `gorm:"column:cart_name"`
	CreatedAt  time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt  time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

func (m *Cart) TableName() string {
	return "cart"
}
