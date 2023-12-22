package service

import (
	"bytes"
	"fileStore/internel/biz"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/response"
	"github.com/gin-gonic/gin"
	"io"
)

type FileUploadReq struct {
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
}

func FileUpload(c *gin.Context) {
	var req FileUploadReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", nil))
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", nil))
		return
	}
	//将文件转为[]byte
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", nil))
		return
	}
	//存入本地
	fileMeta, err := biz.StoreFileLocal(buf.Bytes(), fileHeader.Filename, fileHeader.Size)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileStoreFail, "文件获取错误", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件存储成功", fileMeta))
	//进行转存 转存完成后的file表的更新
}
