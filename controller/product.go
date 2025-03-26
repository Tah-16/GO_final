package controller

import (
	"go-basic/dto"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ProductController รับ router และฐานข้อมูล DB เพื่อจัดการเส้นทางที่เกี่ยวกับสินค้า
func ProductController(router *gin.Engine, db *gorm.DB) {
	routers := router.Group("/Product")
	{
		routers.GET("/searchproducts", func(c *gin.Context) { searchProducts(c, db) })
		routers.POST("/addproducttoCart", func(c *gin.Context) { addProductToCart(c, db) })
	}
}

// searchProducts ใช้ค้นหาสินค้าจากคำค้นหาและช่วงราคา
func searchProducts(c *gin.Context, p_db *gorm.DB) {
	if p_db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// รับข้อมูลคำค้นหาจากคำขอ (ใช้ DTO)
	var req dto.ProductSearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// ตรวจสอบค่าราคาเพื่อให้แน่ใจว่าเป็นค่าที่สมเหตุสมผล
	if req.MinPrice < 0 || req.MaxPrice < 0 || req.MinPrice > req.MaxPrice {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price range"})
		return
	}

	// สร้างคำสั่งค้นหาสินค้า
	var products []model.Product
	query := p_db.Model(&model.Product{})

	// ค้นหาด้วย keyword ในชื่อหรือคำอธิบาย
	if req.Keyword != "" {
		query = query.Where("product_name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// ค้นหาด้วยช่วงราคา
	if req.MinPrice > 0 {
		query = query.Where("price >= ?", req.MinPrice)
	}
	if req.MaxPrice > 0 {
		query = query.Where("price <= ?", req.MaxPrice)
	}

	// ดึงข้อมูลสินค้าที่ตรงกับเงื่อนไข
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to search products"})
		return
	}

	// ส่งผลลัพธ์สินค้ากลับไป
	c.JSON(http.StatusOK, gin.H{"products": products})
}
func addProductToCart(c *gin.Context, db *gorm.DB) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	// รับข้อมูลจากคำขอ
	var req dto.AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	// ค้นหารถเข็นที่มีชื่อที่กำหนดโดยลูกค้า
	var cart model.Cart
	if err := db.Where("customer_id = ? AND cart_name = ?", req.CustomerID, req.CartName).First(&cart).Error; err != nil {
		// หากไม่พบรถเข็น ให้สร้างรถเข็นใหม่
		cart = model.Cart{
			CustomerID: req.CustomerID,
			CartName:   req.CartName,
		}
		if err := db.Create(&cart).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create cart"})
			return
		}
	}

	// ค้นหาสินค้าในตะกร้าว่ามีสินค้าเดียวกันหรือไม่
	var cartItem model.CartItem
	if err := db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error; err == nil {
		// หากมีสินค้าแล้ว ให้เพิ่มจำนวนสินค้า
		cartItem.Quantity += req.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update cart item quantity"})
			return
		}
	} else if err := db.Where("cart_id = ?", cart.CartID).First(&cartItem).Error; err != nil {
		// หากไม่มีสินค้าในตะกร้า ให้เพิ่มสินค้าใหม่
		cartItem = model.CartItem{
			CartID:    cart.CartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		}
		if err := db.Create(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add product to cart"})
			return
		}
	}

	// ส่งผลลัพธ์สำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": "Product added to cart successfully"})
}
