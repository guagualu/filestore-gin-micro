package router

import "github.com/gin-gonic/gin"

func Router(r *gin.Engine) {
	g := r.Group("/user")
	g.GET("/signUp")
}
