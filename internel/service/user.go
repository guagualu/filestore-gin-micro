package service

import (
	"context"
	"fileStore/internel/biz"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/encoding"
	"fileStore/internel/pkg/jwt"
	"fileStore/internel/pkg/response"
	"fileStore/log"
	"github.com/gin-gonic/gin"
	"net/http"
)

var salt = "filestore" //盐值 用于加密

type signUpReq struct {
	UserName     string `form:"user_name" uri:"user_name" json:"user_name" binding:"required"`
	UserPassword string `form:"user_password" uri:"user_password" json:"user_password" binding:"required"`
	Email        string `form:"email" uri:"email" json:"email" binding:"required"`
	Mobile       string `form:"mobile" uri:"mobile" json:"mobile" binding:"required"`
}

type signInReq struct {
	Mobile       string `form:"mobile" uri:"mobile" json:"mobile" binding:"required"`
	UserPassword string `form:"user_password" uri:"user_password" json:"user_password" binding:"required"`
}
type userInfoReq struct {
	Uuid string `form:"uuid" uri:"uuid" json:"uuid" binding:"required"`
}

func SignUp(c *gin.Context) {
	//1、获取客户端字段 并进行有效性验证
	var req signUpReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数校验失败", nil))
		log.Logger.Error("参数校验失败:", err)
		return
	}

	//2、密码加密
	password := encoding.Sha1([]byte(req.UserPassword + salt))
	//3、进行db操作 insert操作
	err = biz.SignupUserinfo(context.Background(), req.UserName, password, req.Email, req.Mobile)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.Faild, "注册失败", nil))
		return
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "注册成功", nil))
	return
}

func SignIn(c *gin.Context) {
	//1、获取客户端字段
	var req signInReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数校验失败", nil))
		log.Logger.Error("参数校验失败:", err)
		return
	}
	//2、数据库获取验证
	password := encoding.Sha1([]byte(req.UserPassword + salt))
	u, err := biz.SignInUserinfo(context.Background(), password, req.Mobile)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.Faild, "登录失败", nil))
		return
	}
	//3、生成token
	token, err := jwt.GenerateToken(u.Uuid)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.Faild, "Token生成失败", nil))
		return
	}
	//4、上传token 和成功信息(可以前端实现)
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "登录成功", token))
	return

}

func GetUserInfo(c *gin.Context) {
	//1、解析请求参数
	var req userInfoReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数校验失败", nil))
		log.Logger.Error("参数校验失败:", err)
		return
	}
	user, err := biz.GetUserInfo(context.Background(), req.Uuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.Faild, "获取失败", nil))
		return
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "登录成功", user))
	return
}
