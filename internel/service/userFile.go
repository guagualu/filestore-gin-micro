package service

import (
	"context"
	"fileStore/internel/biz"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ListUserFilesReq struct {
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
	Page     int    `form:"page"  json:"page"`
	PageSize int    `form:"page_size"  json:"page_size"`
}
type DeletedUserFilesReq struct {
	FileHashs []string `form:"file_hashs"  json:"file_hashs" binding:"required"`
	UserUuid  string   `form:"user_uuid"  json:"user_uuid" binding:"required"`
}

type RenameUserFileReq struct {
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
	FileName string `form:"file_name" json:"file_name" binding:"required"`
}

type ListUserFilesFileinfo struct {
	FileHash  string `json:"file_hash"`
	FileName  string `json:"file_name"`
	CreatedAt string `json:"created_at"`
	UpdateAt  string `json:"update_at"`
}

type ListUserFilesRes struct {
	FileList []ListUserFilesFileinfo `json:"file_list"`
	Count    int                     `json:"count"`
}

func ListUserFiles(c *gin.Context) {
	var req ListUserFilesReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	list, count, err := biz.GetUserFileList(context.Background(), req.UserUuid, req.Page, req.PageSize)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.ListUserFileErr, "获取用户文件列表失败", nil))
		return
	}
	listRes := make([]ListUserFilesFileinfo, 0)
	for _, v := range list {
		listRes = append(listRes, ListUserFilesFileinfo{
			FileHash:  v.FileHash,
			FileName:  v.FileName,
			CreatedAt: v.CreateAt.String(),
			UpdateAt:  v.UpdateAt.String(),
		})
	}

	res := ListUserFilesRes{
		FileList: listRes,
		Count:    int(count),
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "获取用户文件列表成功", res))
	return

}

func DeletedUserFiles(c *gin.Context) {
	var req DeletedUserFilesReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	err = biz.DeletedUserFileList(context.Background(), req.UserUuid, req.FileHashs)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.DeleteUserFilesErr, "删除用户文件失败", nil))
		return
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "删除用户文件列表成功", nil))
	return
}

func RenameUserFile(c *gin.Context) {
	var req RenameUserFileReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	err = biz.RenameUserFile(context.Background(), req.UserUuid, req.FileHash, req.FileName)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.RenameUserFileErr, "重命名用户文件失败", nil))
		return
	}
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "重命名用户文件成功", nil))
	return
}
