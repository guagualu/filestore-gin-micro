package middleware

import (
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/jwt"
	"fileStore/internel/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

//gin.H{
//"message": "Missing Authorization Header",
//}

// 定义一个中间件函数，用于验证 JWT
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求头中获取 JWT
		tokenString := c.GetHeader("Authorization")
		if len(tokenString) == 0 {
			c.JSON(http.StatusUnauthorized, response.NewRespone(errcode.TokenIsErr, "Missing Authorization Header", nil))
			c.Abort()
			return
		}

		// 验证 JWT
		userUuid, err := jwt.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, response.NewRespone(errcode.TokenIsErr, "Invalid token", nil))
			c.Abort()
			return
		}

		// 将用户 ID 注入到上下文中
		c.Set("userUuid", userUuid)
		c.Next()
	}
}
