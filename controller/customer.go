package controller

import (
	"go-basic/dto"
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var c_db *gorm.DB

func CustomerController(router *gin.Engine, db *gorm.DB) {
	c_db = db
	routers := router.Group("/auth")
	{
		routers.POST("/login", func(c *gin.Context) { postLogin(c, c_db) })
		routers.POST("/login", func(c *gin.Context) { postLogin(c, c_db) })
	}
}

func postLogin(c *gin.Context, c_db *gorm.DB) {
	if c_db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	var req dto.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var customer model.Customer
	// ตรวจสอบว่ามี email ในฐานข้อมูลหรือไม่
	if err := c_db.Where("email = ?", req.Email).First(&customer).Error; err != nil {
		// หากไม่พบ email, สร้างผู้ใช้ใหม่
		// ทำการ hash password ก่อนบันทึกลงฐานข้อมูล
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}

		// สร้างข้อมูลลูกค้าใหม่
		customer = model.Customer{
			FirstName:   "Tah",
			LastName:    "za",
			Email:       req.Email,
			Password:    string(hashedPassword), // บันทึก password ที่ถูก hash
			PhoneNumber: "0912345678",
			Address:     "Bangkok",
		}

		// บันทึกผู้ใช้ใหม่ลงในฐานข้อมูล
		if err := c_db.Create(&customer).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		// ส่งข้อมูลผู้ใช้กลับ
		c.JSON(http.StatusOK, gin.H{
			"user": gin.H{
				"id":           customer.CustomerID,
				"first_name":   customer.FirstName,
				"last_name":    customer.LastName,
				"email":        customer.Email,
				"phone_number": customer.PhoneNumber,
				"address":      customer.Address,
			},
		})
		return
	}

	// หากพบ email ในฐานข้อมูล, ตรวจสอบ password
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// ส่งข้อมูลผู้ใช้กลับ
	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":           customer.CustomerID,
			"first_name":   customer.FirstName,
			"last_name":    customer.LastName,
			"email":        customer.Email,
			"phone_number": customer.PhoneNumber,
			"address":      customer.Address,
		},
	})
}
func changePassword(c *gin.Context, c_db *gorm.DB) {
	if c_db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is not initialized"})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	var customer model.Customer
	if err := c_db.Where("email = ?", req.Email).First(&customer).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer not found"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Old password is incorrect"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
		return
	}

	customer.Password = string(hashedPassword)
	if err := c_db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
}
