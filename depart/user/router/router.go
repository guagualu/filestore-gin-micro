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
	g := r.Group("/user")
	g.Use(middleware.Cors())
	{
		g.POST("/signUp", service.SignUp)
		g.POST("/signIn", service.SignIn)
	}
	g.Use(middleware.JWTMiddleware())
	{
		g.GET("/info", service.GetUserInfo)
	}

}
