package router

import (
	"fileStore/internel/middleware"
	"fileStore/internel/service"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	g := r.Group("/file")
	g.Use(middleware.JWTMiddleware())
	{
		g.POST("/upload", service.FileUpload)
	}

}
