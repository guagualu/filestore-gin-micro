package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fileStore/conf"
	"fileStore/internel/biz"
	"fileStore/internel/domain"
	"fileStore/internel/middleware/mq"
	"fileStore/internel/middleware/oss"
	"fileStore/internel/pkg/code/errcode"
	"fileStore/internel/pkg/code/sucesscode"
	"fileStore/internel/pkg/encoding"
	"fileStore/internel/pkg/response"
	"fileStore/log"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type FileUploadReq struct {
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
}

type FileFastUploadReq struct {
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
	FileName string `form:"file_name"  json:"file_name" binding:"required"`
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

type DownloadRsp struct {
	Url       string `json:"url"`
	StoreType string `json:"store_type"`
}

type DownloadReq struct {
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
	UserUuid string `form:"user_uuid"  json:"user_uuid" binding:"required"`
	FileName string `form:"file_name"  json:"file_name" binding:"required"`
}

type PreFileInfoRes struct {
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
	FileSize int    `form:"file_size"  json:"file_size" binding:"required"`
}
type GetFileInfoReq struct {
	FileHash string `form:"file_hash"  json:"file_hash" binding:"required"`
}

type GetFileInfoRes struct {
	ID       uint   `json:"id"`
	FileHash string `json:"file_hash"`
	FileName string `json:"file_name"`
	FileSize int64  ` json:"file_size"`
	FileAddr string `json:"file_addr"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}

type ChatReq struct {
	SendUserUuid string `form:"send_user_uuid" json:"send_user_uuid" binding:"required"`
	ToUserUuid   string `form:"to_user_uuid" json:"to_user_uuid" binding:"required"`
	SessionUuid  string `form:"session_uuid" json:"session_uuid" binding:"required"`
}

type CreateSessionReq struct {
	UserAUuid string `form:"user_a_uuid" json:"user_a_uuid" binding:"required"`
	UserBUuid string `form:"user_b_uuid" json:"user_b_uuid" binding:"required"`
}

type GetUserAllSessionReq struct {
	UserUuid string `form:"user_uuid" json:"user_uuid" binding:"required"`
}

type GetUserAllSessionStructRes struct {
	Id          uint        `json:"id"`           // 物理主键
	SessionUuid string      `json:"session_uuid"` // 业务主键
	ChatUser    domain.User `json:"chat_user"`    // 接受者uuid
	CreatedAt   string      `json:"created_at"`
	UpdatedAt   string      `json:"updated_at"`
}

type GetUserAllSessionRes struct {
	List []GetUserAllSessionStructRes `json:"list"`
}

type GetUserASessionInfoReq struct {
	SessionUuid string `form:"session_uuid" json:"session_uuid" binding:"required"`
}
type GetUserASessionInfoStructRes struct {
	Id             uint        `json:"id"`
	SessionUuid    string      `json:"session_uuid"`
	SendUser       domain.User `json:"send_user"`    // 发送者
	ToUser         domain.User `json:"to_user"`      // 接受者
	MessageType    int         `json:"message_type"` //0为普通文本消息，1为文件消息 content存的是文件hash
	MessageContent string      `json:"message_content"`
	CreatedAt      string      `json:"created_at"`
	UpdatedAt      string      `json:"updated_at"`
}
type GetUserASessionInfoRes struct {
	ContentList []GetUserASessionInfoStructRes `json:"content_list"`
}

type IsSessionExitReq struct {
	UserAUuid string `form:"user_a_uuid" json:"user_a_uuid" binding:"required"`
	UserBUuid string `form:"user_b_uuid" json:"user_b_uuid" binding:"required"`
}

func GetFileInfo(c *gin.Context) {
	var req GetFileInfoReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	info, err := biz.GetFileInfo(context.Background(), req.FileHash)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.FileGetFail, "文件信息获取错误", nil))
		return
	}
	res := GetFileInfoRes{
		ID:       info.ID,
		FileHash: info.FileHash,
		FileName: info.FileName,
		FileSize: info.FileSize,
		FileAddr: info.FileAddr,
		CreateAt: info.CreateAt.String(),
		UpdateAt: info.UpdateAt.String(),
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件信息获取成功", res))
	return
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
	////进行转存 转存完成后的file表的更新
	//fileInfo := mq.MqFileInfo{
	//	FileHash:    fileMeta.FileHash,
	//	FileName:    fileMeta.FileName,
	//	CurLocateAt: fileMeta.FileAddr,
	//}
	//err = biz.StoreFileOss(buf.Bytes(), fileInfo)
	//if err != nil {
	//	c.JSON(400, response.NewRespone(errcode.FileStoreFail, "文件转存失败", nil))
	//	return
	//}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件存储成功", fileMeta))
}

func FileFastUpload(c *gin.Context) {
	var req FileFastUploadReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	//文件快传
	err = biz.FastUpload(context.Background(), req.FileHash, req.FileName, req.UserUuid)
	if err != nil {
		c.JSON(500, response.NewRespone(errcode.FileFastUploadFail, "文件快传错误", err))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件快传成功", nil))
	return
}

func PreFileInfo(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", nil))
		return
	}
	file, err := fileHeader.Open()
	//将文件转为[]byte
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", nil))
		return
	}
	fileHash := encoding.Sha1(buf.Bytes())
	fileSize := fileHeader.Size
	c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "文件hash获取成功", PreFileInfoRes{
		FileHash: fileHash,
		FileSize: int(fileSize),
	}))
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
	//var req MultipartUploadReq
	//err := c.ShouldBind(&req)
	//if err != nil {
	//	c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", err))
	//	return
	//}
	uploadId := c.Request.Header.Get("upload_id")
	chunkIndexStr := c.Request.Header.Get("chunk_index")
	chunkIndex, _ := strconv.Atoi(chunkIndexStr)
	mpInfo := domain.MultipartUploadInfo{
		UploadID:   uploadId,
		ChunkIndex: chunkIndex,
	}
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.FileGetFail, "文件获取错误", err.Error()))
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
		c.JSON(200, response.NewRespone(errcode.FileStoreFail, "正在分块上传中", MultipartUploadCompleteRsp{
			Completed: false,
			Progress:  int((float32(num) / float32(req.ChunkCount)) * 100),
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
			log.Logger.Error("本地文件打开失败 err:", err)
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
		c.JSON(500, response.NewRespone(errcode.RetryErr, "文件分块上传重试初始化失败", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "文件分块上传重试初始化成功", ReTryFileMpUploadInitRsp{chunkArray}))
	return
}

// DownloadHandler : 文件下载接口
func Download(c *gin.Context) {
	var req DownloadReq
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	// TODO: 处理异常情况
	//获取file信息以及查询是否属于当前用户
	fileInfo, ferr := biz.GetFileInfo(context.Background(), req.FileHash)
	userFile, uferr := biz.GetUserFileInfo(context.Background(), req.FileHash, req.UserUuid, req.FileName)
	if ferr != nil || uferr != nil || fileInfo == nil || userFile == nil {
		c.JSON(400, response.NewRespone(errcode.DownloadFileNotValid, "所要下载的文件不存在或者不属于当前用户", nil))
		return
	}
	if strings.HasPrefix(fileInfo.FileAddr, conf.GetConfig().LocalStore) {
		// 本地文件， 直接下载
		//具体来说，它将文件从指定的路径读取并发送给客户端，以便用户下载或查看文件。其中，参数c是一个gin.Context类型的指针，它表示当前的HTTP请求上下文。uniqFile.FileAddr.String表示文件的路径，userFile.FileName表示文件的文件名。
		c.FileAttachment(fileInfo.FileAddr, userFile.FileName)
	} else if strings.HasPrefix(fileInfo.FileAddr, "oss") {
		// oss中的文件
		signedURL := oss.DownloadURL(fileInfo.FileAddr)
		c.JSON(http.StatusOK, response.NewRespone(sucesscode.Success, "oss下载url", DownloadRsp{
			Url:       signedURL,
			StoreType: "oss",
		}))
	}
}

func Chat(c *gin.Context) {
	req := ChatReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	//校验token
	tokenStruct := struct {
		Token string `json:"token"`
	}{}
	conn.ReadJSON(&tokenStruct)
	fmt.Println(tokenStruct.Token)
	//创建channel
	sendChannel := make(chan domain.ImSendMsg, 100)
	toChnnel := make(chan domain.ImSendMsg, 100)
	if domain.GetKey(req.SessionUuid+req.SendUserUuid) == nil {
		domain.StoreKey(req.SessionUuid+req.SendUserUuid, sendChannel)
	}
	if domain.GetKey(req.SessionUuid+req.ToUserUuid) == nil {
		domain.StoreKey(req.SessionUuid+req.ToUserUuid, toChnnel)
	}
	go Read(conn, req.SessionUuid+req.SendUserUuid, req.SessionUuid+req.ToUserUuid)
	go Write(conn, req.SessionUuid+req.SendUserUuid)
}

type CreateSessionRes struct {
	SessionUuid string `json:"session_uuid"`
}

func CreateSession(c *gin.Context) {
	req := CreateSessionReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	imsession, err := biz.CreateSession(context.Background(), req.UserAUuid, req.UserBUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.CreateSessionErr, "创建session错误", nil))
		return
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "创建session成功", CreateSessionRes{SessionUuid: imsession.SessionUuid}))
}

func GetUserAllSession(c *gin.Context) {
	req := GetUserAllSessionReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	list, err := biz.GetUserAllSession(context.Background(), req.UserUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.GetSessionAllErr, "获取用户session错误", nil))
		return
	}
	res := GetUserAllSessionRes{}
	res.List = make([]GetUserAllSessionStructRes, 0)
	for _, v := range list {
		chatToUuid := ""
		//如果不是那么就是对方
		if v.UserAUuid != req.UserUuid {
			chatToUuid = v.UserAUuid
		} else {
			chatToUuid = v.UserBUuid
		}
		chatToUser, err := biz.FileServerGetUserInfo(context.Background(), chatToUuid)
		if err != nil {
			c.JSON(400, response.NewRespone(errcode.GetSessionAllErr, "获取用户session错误", nil))
			return
		}
		res.List = append(res.List, GetUserAllSessionStructRes{
			Id:          v.Id,
			SessionUuid: v.SessionUuid,
			ChatUser:    *chatToUser,
			CreatedAt:   v.CreatedAt.String(),
			UpdatedAt:   v.UpdatedAt.String(),
		})
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "获取用户session成功", res))
}

func GetUserASessionInfo(c *gin.Context) {
	req := GetUserASessionInfoReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	list, err := biz.GetUserASession(context.Background(), req.SessionUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.GetSessionInfoErr, "获取session信息错误", nil))
		return
	}
	res := GetUserASessionInfoRes{}
	res.ContentList = make([]GetUserASessionInfoStructRes, 0)
	for _, v := range list {
		sendUser, err := biz.FileServerGetUserInfo(context.Background(), v.SendUserUuid)
		if err != nil {
			c.JSON(500, response.NewRespone(errcode.GetSessionAllErr, "获取用户session错误1", err))
			return
		}
		ToUser, err := biz.FileServerGetUserInfo(context.Background(), v.ToUserUuid)
		if err != nil {
			c.JSON(500, response.NewRespone(errcode.GetSessionAllErr, "获取用户session错误2", err))
			return
		}
		res.ContentList = append(res.ContentList, GetUserASessionInfoStructRes{
			Id:             v.Id,
			SessionUuid:    v.SessionUuid,
			SendUser:       *sendUser,
			ToUser:         *ToUser,
			MessageType:    v.MessageType,
			MessageContent: v.MessageContent,
			CreatedAt:      v.CreatedAt.String(),
			UpdatedAt:      v.UpdatedAt.String(),
		})
	}
	c.JSON(200, response.NewRespone(sucesscode.Success, "获取session信息成功", res))
}

func Read(conn *websocket.Conn, sendKey, toKey string) {
	for {
		conn.PongHandler() //心跳？
		sendMsg := domain.ImSendMsg{}
		// _,msg,_:=c.Socket.ReadMessage()
		err := conn.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		fmt.Println(sendMsg)
		if err != nil {
			log.Logger.Error("im read err1:", err)
			break
		}
		//如果是文件 那么内容是hash，那么存储一个结构体包含发送者
		if sendMsg.MessageType == 1 {
			fileInfo := domain.ImSessionFileContent{}
			json.Unmarshal([]byte(sendMsg.Message), &fileInfo)
			fileContent := domain.ImSessionFileContent{
				FileHash:   fileInfo.FileHash,
				FileName:   fileInfo.FileName,
				FileSender: sendMsg.SendUserUuid,
			}
			fileContentJson, _ := json.Marshal(fileContent)
			sendMsg.Message = string(fileContentJson)
		}
		//保存到mysql
		msg := domain.ImSessionContent{
			SessionUuid:    sendMsg.SessionUuid,
			SendUserUuid:   sendMsg.SendUserUuid,
			ToUserUuid:     sendMsg.ToUserUuid,
			MessageType:    sendMsg.MessageType,
			MessageContent: sendMsg.Message,
		}
		SessionContentId, err := biz.SaveASessionContent(context.Background(), msg)
		if err != nil {
			log.Logger.Error("im read err2:", err)
			break
		}
		//刷新session的更新事件
		err = biz.UpdateSessionUpdateTime(context.Background(), msg.SessionUuid)
		if err != nil {
			log.Logger.Error("im read err3:", err)
			break
		}
		sendMsg.SessionContentId = SessionContentId
		//写入channel进行通知
		toChan := domain.GetKey(sendMsg.SessionUuid + sendMsg.ToUserUuid)
		if toChan != nil { //因为可能此时对方连接已经关闭 给一个nil channel发送数据会造成永久的阻塞
			toChan <- sendMsg
		}

		fmt.Println("gg", sendMsg)
	}
	//删除channel 只删除自己的channel 否则可能对方还在就把对方的channel删了
	domain.DeleteKey(sendKey)
	_ = conn.Close()
}

func Write(conn *websocket.Conn, reciveKey string) {
	defer func() {
		_ = conn.Close()
	}()
	for {
		reciveChannel := domain.GetKey(reciveKey)
		message, ok := <-reciveChannel
		if !ok {
			_ = conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}
		fmt.Println(message.ToUserUuid, "接受消息:", message.Message)
		res := response.NewRespone(sucesscode.Success, "聊天写入成功", message)
		msg, _ := json.Marshal(res)
		_ = conn.WriteMessage(websocket.TextMessage, msg)
	}

}

type IsSessionExitRes struct {
	IsExist     bool   `json:"is_exist"`
	SessionUuid string `json:"session_uuid"`
}

func IsSessionExit(c *gin.Context) {
	req := IsSessionExitReq{}
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.ValidationFaild, "参数错误", nil))
		return
	}
	session, exist, err := biz.IsSeesionExist(context.Background(), req.UserAUuid, req.UserBUuid)
	if err != nil {
		c.JSON(400, response.NewRespone(errcode.GetSessionInfoErr, "获取session信息错误", nil))
		return
	}
	res := IsSessionExitRes{IsExist: exist, SessionUuid: session.SessionUuid}
	c.JSON(200, response.NewRespone(sucesscode.Success, "获取session信息成功", res))
}
