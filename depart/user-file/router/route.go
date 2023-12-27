package router

import (
	"fileStore/internel/middleware"
	"fileStore/internel/service"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	g := r.Group("/user-file")
	g.Use(middleware.JWTMiddleware())
	{
		g.POST("/list", service.ListUserFiles)
		g.POST("/rename", service.RenameUserFile)
		g.POST("/delete", service.DeletedUserFiles)
	}

}
