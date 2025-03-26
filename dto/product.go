package dto

type ProductSearchRequest struct {
	Keyword  string  `json:"keyword" binding:"required"`   // คำค้นหาสินค้า
	MinPrice float64 `json:"min_price" binding:"required"` // ราคาต่ำสุด
	MaxPrice float64 `json:"max_price" binding:"required"` // ราคาสูงสุด
}
