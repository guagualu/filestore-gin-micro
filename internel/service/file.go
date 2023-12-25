package service

import (
	"bytes"
	"context"
	"fileStore/internel/biz"
	"fileStore/internel/domain"
	"fileStore/internel/middleware/mq"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/response"
	"fileStore/log"
	"github.com/gin-gonic/gin"
	"io"
	"os"
)

type FileUploadReq struct {
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
}

// MultipartUploadInfo : 初始化信息
type MultipartUploadInitReq struct {
	FileHash   string `form:"file_hash"  json:"file_hash" binding:"required"`
	FileSize   int    `form:"file_size"  json:"file_size" binding:"required"`
	UploadID   string `form:"upload_id"  json:"upload_id" binding:"required"`
	ChunkSize  int    `form:"chunk_size"  json:"chunk_size" binding:"required"`
	ChunkCount int    `form:"chunk_count"  json:"chunk_count" binding:"required"`
}

type MultipartUploadReq struct {
	UploadID   string `form:"upload_id"  json:"upload_id" binding:"required"`
	ChunkIndex int    `form:"chunk_index"  json:"chunk_index" binding:"required"`
}

type MultipartUploadCompleteReq struct {
	FileHash   string `form:"file_hash"  json:"file_hash" binding:"required"`
	FileName   string `form:"file_name"  json:"file_name" binding:"required"`
	FileSize   int    `form:"file_size"  json:"file_size" binding:"required"`
	UploadID   string `form:"upload_id"  json:"upload_id" binding:"required"`
	ChunkCount int    `form:"chunk_count"  json:"chunk_count" binding:"required"`
	UserUuid   string `form:"user_uuid"  json:"user_uuid" binding:"required"`
}
type MultipartUploadCompleteRsp struct {
	Completed bool `json:"completed"`
	Progress  int  `json:"progress"`
}

type ReTryFileMpUploadInitReq struct {
	UploadID   string `form:"upload_id"  json:"upload_id" binding:"required"`
	ChunkCount int    `form:"chunk_count"  json:"chunk_count" binding:"required"`
}

type ReTryFileMpUploadInitRsp struct {
	ChunkIndexArray []int ` json:"chunk_index_array"`
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
	fileMeta, err := biz.StoreFileLocal(buf.Bytes(), fileHeader.Filename, fileHeader.Size, req.UserUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileStoreFail, "文件获取错误", nil))
		return
	}
	//进行转存 转存完成后的file表的更新
	fileInfo := mq.MqFileInfo{
		FileHash:    fileMeta.FileHash,
		FileName:    fileMeta.FileName,
		CurLocateAt: fileMeta.FileAddr,
	}
	err = biz.StoreFileOss(buf.Bytes(), fileInfo)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileStoreFail, "文件转存失败", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件存储成功", fileMeta))
}

func FileFastUpload(c *gin.Context) {
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
	//文件快传
	err = biz.FastUpload(context.Background(), buf.Bytes(), fileHeader.Filename, req.UserUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileFastUploadFail, "文件快传错误", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件快传成功", nil))
	return
}

func FileMpUploadInit(c *gin.Context) {
	var req MultipartUploadInitReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	mpInfo := domain.MultipartUploadInfo{
		FileHash:   req.FileHash,
		FileSize:   req.FileSize,
		UploadID:   req.UploadID,
		ChunkSize:  req.ChunkSize,
		ChunkCount: req.ChunkCount,
	}
	err = biz.FileMploadInit(mpInfo)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileFastUploadFail, "分块上传初始化错误", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件分块上传初始化成功", nil))
	return
}

func FileMpUpload(c *gin.Context) {
	var req MultipartUploadReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	mpInfo := domain.MultipartUploadInfo{
		UploadID:   req.UploadID,
		ChunkIndex: req.ChunkIndex,
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
	//存入本地并完成redis里面的记录更新
	err = biz.FileMploadLocal(buf.Bytes(), mpInfo.UploadID, mpInfo.ChunkIndex)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileStoreFail, "文件存储错误", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件分块上传成功", nil))
	return
}

// 如果没有完成 通知前端进度 如果完成分块上传 完成分块文件的合成和转储 通知前端
func CompleteFileMpUpload(c *gin.Context) {
	var req MultipartUploadCompleteReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	//检查进度
	num, err := biz.FileMpUploadCheck(req.UploadID, req.ChunkCount)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileMpCheckFail, "文件分块检查错误", nil))
		return
	}
	if num != req.ChunkCount {
		c.JSON(400, response.NewRespone(errcode.FileStoreFail, "正在分块上传中", MultipartUploadCompleteRsp{
			Completed: false,
			Progress:  (num / req.ChunkCount) * 100,
		}))
	} else {
		//分块上传完毕 完成分块文件组装、数据库更新
		locatedAddr, err := biz.FileMpUploadStore(req.UploadID, req.FileHash, req.FileName, req.UserUuid, int64(req.FileSize))
		if err != nil {
			c.JSON(400, response.NewRespone(errcode.FileMpCheckFail, "文件分块存储错误", nil))
			return
		}
		c.JSON(200, response.NewRespone(sucesscode.Success, "文件完成分块存储", MultipartUploadCompleteRsp{
			Completed: true,
			Progress:  100,
		}))
		//进行oss异步转存
		//进行转存 转存完成后的file表的更新
		fileInfo := mq.MqFileInfo{
			FileHash:    req.FileHash,
			FileName:    req.FileName,
			CurLocateAt: locatedAddr,
		}
		//将文件转为[]byte
		file, err := os.Open(locatedAddr)
		if err != nil {
			log.Logger.Error("本地文件打开失败")
			return
		}
		buf := bytes.NewBuffer(nil)
		_, err = io.Copy(buf, file)
		err = biz.StoreFileOss(buf.Bytes(), fileInfo)
		if err != nil {
			log.Logger.Error("转存oss失败")
			return
		}
	}
	return
}

// 重试
func ReTryFileMpUploadInit(c *gin.Context) {
	var req ReTryFileMpUploadInitReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	chunkArray, err := biz.CheckFailedMpUploadFile(req.UploadID, req.ChunkCount)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.RetryErr, "文件分块上传重试初始化失败", nil))
		return
	}
	c.JSON(400, response.NewRespone(sucesscode.Success, "文件分块上传重试初始化成功", ReTryFileMpUploadInitRsp{chunkArray}))
	return
}
