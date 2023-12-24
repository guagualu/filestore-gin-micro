package router

import (
	"fileStore/internel/middleware"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	g := r.Group("/user-file")
	g.Use(middleware.JWTMiddleware())
	{
		//g.POST("/upload", service.FileUpload)
	}

}
