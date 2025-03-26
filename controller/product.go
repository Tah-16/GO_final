package controller

import (
	"go-basic/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ProductController(router *gin.Engine) {
	routers := router.Group("/Product")
	{
		routers.GET("/", product)
	}
}
func product(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)
	if db == nil {
		panic("DB Connection not found")
	}

	var products []model.Product

	if err := db.Find(&products).Error; err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
