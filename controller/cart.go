package controller

import (
	"fmt"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CartController(router *gin.Engine, db *gorm.DB) {
	routers := router.Group("/cart")
	{
		routers.GET("/getcarts", func(c *gin.Context) { getAllCarts(c, db) })
	}
}

func getAllCarts(c *gin.Context, db *gorm.DB) {
	customerID := c.Query("customer_id")
	fmt.Println("Customer ID:", customerID)

	var carts []model.Cart
	if err := db.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Carts not found"})
		return
	}

	var cartDetails []gin.H
	for _, cart := range carts {
		var cartItems []model.CartItem

		if err := db.Where("cart_id = ?", cart.CartID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch cart items"})
			return
		}

		var cartItemDetails []gin.H
		for _, item := range cartItems {
			var product model.Product
			if err := db.Where("product_id = ?", item.ProductID).First(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch product details"})
				return
			}

			totalPrice := float64(item.Quantity) * product.Price

			cartItemDetails = append(cartItemDetails, gin.H{
				"product_name": product.ProductName,
				"quantity":     item.Quantity,
				"price":        product.Price,
				"total_price":  totalPrice,
			})
		}

		cartDetails = append(cartDetails, gin.H{
			"cart_name": cart.CartName,
			"items":     cartItemDetails,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"carts": cartDetails,
	})
}
