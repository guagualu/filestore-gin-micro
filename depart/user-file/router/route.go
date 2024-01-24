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
		g.GET("/get/userfile", service.GetUserFile)
		g.POST("/create", service.CreateUserFile)
		g.POST("/list", service.ListUserFiles)
		g.POST("/rename", service.RenameUserFile)
		g.POST("/delete", service.DeletedUserFiles)
		g.POST("/clash/list", service.ListDeletedUserFiles)
		g.POST("/clash/delete", service.TrueDeletedUserFiles)
		g.POST("/clash/recover", service.RecoverDeletedUserFiles)
	}

}
