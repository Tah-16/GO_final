package dto

// AddToCartRequest DTO สำหรับการเพิ่มสินค้าในรถเข็น
type AddToCartRequest struct {
	CustomerID int    `json:"customer_id"` // รหัสลูกค้า
	CartName   string `json:"cart_name"`   // ชื่อรถเข็น
	ProductID  int    `json:"product_id"`  // รหัสสินค้า
	Quantity   int    `json:"quantity"`    // จำนวนสินค้าที่ต้องการเพิ่ม
}

type CartItemDetail struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
}

type CartDetail struct {
	CartName string           `json:"cart_name"`
	Items    []CartItemDetail `json:"items"`
}
type CartItem struct {
	CartItemID int `json:"cart_item_id" gorm:"primaryKey"`
	CartID     int `json:"cart_id"`
	ProductID  int `json:"product_id"`
	Quantity   int `json:"quantity"`
}
