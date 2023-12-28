package router

import (
	"fileStore/internel/middleware"
	"fileStore/internel/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello, world!",
		})
	})
	g := r.Group("/file")
	// 使用gin插件支持跨域请求
	g.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"http://localhost:8080"}, // []string{"http://localhost:8080"},
		AllowMethods:  []string{"GET", "POST"},
		AllowHeaders:  []string{"Origin", "Range", "x-requested-with", "Content-Type"},
		ExposeHeaders: []string{"Content-Length", "Accept-Ranges", "Content-Range", "Content-Disposition"},
		// AllowCredentials: true,
	}))
	g.Use(middleware.JWTMiddleware())
	{
		g.POST("/upload", service.FileUpload)
		g.POST("/fast/upload", service.FileFastUpload)
		g.GET("/upload/mp/init", service.FileMpUploadInit)
		g.POST("/upload/mp", service.FileMpUpload)
		g.POST("/upload/completed", service.CompleteFileMpUpload)
		g.POST("/upload/retry/init", service.ReTryFileMpUploadInit)
		g.POST("/download", service.Download)
		g.GET("/pre/info", service.PreFileInfo)
	}

}
