package service

import (
	"context"
	"fileStore/internel/biz"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/encoding"
	"fileStore/internel/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

var salt = "filestore" //盐值 用于加密

type signUpReq struct {
	userName     string `uri:"user_name"  binding:"required,gt=10"`
	userPassword string `uri:"user_password" binding:"required,gt=20"`
	Email        string `uri:"email"  binding:"required,gt=10"`
	Mobile       string `uri:"mobile" binding:"required,,eq=11"`
}

func SignUp(c *gin.Context) {
	//1、获取客户端字段 并进行有效性验证
	req := signUpReq{}
	c.ShouldBindUri(&req)

	//2、密码加密
	password := encoding.Sha1([]byte(req.userPassword + salt))
	//3、进行db操作 insert操作
	err := biz.SignupUserinfo(context.Background(), req.userName, password, req.userPassword, req.Mobile)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.Faild, "注册失败", nil))
		return
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "注册成功", nil))
	return
}
