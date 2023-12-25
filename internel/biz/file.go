package biz

import (
	"context"
	"encoding/json"
	"fileStore/conf"
	"fileStore/depart/file/client"
	filePb "fileStore/depart/file/proto"
	userFilePb "fileStore/depart/user-file/proto"
	"fileStore/internel/data"
	"fileStore/internel/domain"
	"fileStore/internel/middleware/mq"
	"fileStore/internel/pkg/encoding"
	"fileStore/internel/pkg/util"
	"fileStore/internel/pkg/uuid"
	"fileStore/log"
	"fmt"
	dtmgrpc "github.com/dtm-labs/client/dtmgrpc"
	"os"
	"strconv"
)

func StoreFileLocal(fileData []byte, fileName string, fileSize int64, userUuid string) (*domain.File, error) {
	//计算hash
	fileHash := encoding.Sha1(fileData)
	// 存入本地
	locatedAt := conf.GetConfig().LocalStore + "/" + fileHash
	localFile, err := os.Create(locatedAt)
	if err != nil {
		return nil, err
	}
	_, err = localFile.Write(fileData)
	if err != nil {
		return nil, err
	}
	res := domain.File{
		FileHash: fileHash,
		FileName: fileName,
		FileSize: fileSize,
		FileAddr: locatedAt,
	}
	//写入file表 以及 user-file表 使用分布式事务 dtm saga模式
	DtmServer := conf.GetConfig().DtmServer
	fileReq := filePb.FileReq{
		FileHash:  res.FileHash,
		FileName:  res.FileName,
		FileSize:  res.FileSize,
		LocatedAt: res.FileAddr,
	}
	userFileReq := userFilePb.UserFileReq{
		FileHash: res.FileHash,
		FileName: res.FileName,
		UserUuid: userUuid,
	}
	saga := dtmgrpc.NewSagaGrpc(DtmServer, uuid.NewUuid()).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 补偿操作为url: qsBusi+"/TransOutCompensate"
		Add("127.0.0.1:8083"+"/fileProto.FileService/SaveFile", "127.0.0.1:8083"+"/fileProto.FileService/DeleteFile", &fileReq).
		Add("127.0.0.1:8085"+"/userFileProto.UserFileService/SaveUserFile", "127.0.0.1:8085"+"/userFileProto.UserFileService/DeleteUserFile", &userFileReq)
	//// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransIn"， 补偿操作为url: qsBusi+"/TransInCompensate"
	//Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	saga.WaitResult = true
	saga.RequestTimeout = 1000 * 100
	saga.TimeoutToFail = 1000 * 1000
	saga.RetryLimit = 0
	err = saga.Submit()
	if err != nil {
		fmt.Println("分布式事务执行失败,err:", err)
		return nil, err
	}
	return &res, nil

}

func StoreFileOss(fileData []byte, fileInfo mq.MqFileInfo) error {
	//通过mq publish
	msg, _ := json.Marshal(fileInfo)
	mq.Rabpublish("oss", string(msg))
	return nil
}

func FastUpload(ctx context.Context, fileData []byte, fileName string, userUuid string) error {
	//计算hash
	fileHash := encoding.Sha1(fileData)
	//查找file表
	_, err := data.GetFileByFileHash(ctx, fileHash)
	if err != nil {
		return err
	}
	//存入user-file中
	userFileClient := client.GetUserFileClient()
	_, err = userFileClient.SaveUserFile(ctx, &userFilePb.UserFileReq{
		FileHash: fileHash,
		FileName: fileName,
		UserUuid: userUuid,
	})
	if err != nil {
		return err
	}
	return nil
}

func FileMploadInit(mpFileInfo domain.MultipartUploadInfo) error {
	return data.SaveFileUploadInfo(mpFileInfo)
}

func FileMploadLocal(fileData []byte, uploadId string, chunkIndex int) error {
	// 存入本地
	locatedAt := conf.GetConfig().LocalMpStore + "/" + uploadId + "/" + strconv.Itoa(chunkIndex)
	localFile, err := os.Create(locatedAt)
	if err != nil {
		return err
	}
	_, err = localFile.Write(fileData)
	if err != nil {
		return err
	}
	//redis记录
	info := domain.MultipartUploadInfo{
		UploadID:   uploadId,
		ChunkIndex: chunkIndex,
	}
	err = data.SaveFileMpUpload(info)
	if err != nil {
		return err
	}
	return nil
}

func FileMpUploadCheck(uploadId string, chunkCount int) (int, error) {
	//查找存储的分块文件数目
	info := domain.MultipartUploadInfo{
		UploadID:   uploadId,
		ChunkCount: chunkCount,
	}
	sum, err := data.GetFileMpUploadSum(info)
	if err != nil {
		return 0, err
	}
	return sum, nil
}

// 分块文件进行合并
func FileMpUploadStore(uploadId, fileHash, fileName string, userUuid string, fileSize int64) (string, error) {
	//先进行合并
	locatedAddr, err := FileMpUploadMerge(uploadId, fileHash)
	if err != nil {
		return "", err
	}
	fileReq := filePb.FileReq{
		FileHash:  fileHash,
		FileName:  fileName,
		FileSize:  fileSize,
		LocatedAt: locatedAddr,
	}
	userFileReq := userFilePb.UserFileReq{
		FileHash: fileHash,
		FileName: fileName,
		UserUuid: userUuid,
	}
	DtmServer := conf.GetConfig().DtmServer
	saga := dtmgrpc.NewSagaGrpc(DtmServer, uuid.NewUuid()).
		// 添加一个TransOut的子事务，正向操作为url: qsBusi+"/TransOut"， 补偿操作为url: qsBusi+"/TransOutCompensate"
		Add("127.0.0.1:8083"+"/fileProto.FileService/SaveFile", "127.0.0.1:8083"+"/fileProto.FileService/DeleteFile", &fileReq).
		Add("127.0.0.1:8085"+"/userFileProto.UserFileService/SaveUserFile", "127.0.0.1:8085"+"/userFileProto.UserFileService/DeleteUserFile", &userFileReq)
	//// 添加一个TransIn的子事务，正向操作为url: qsBusi+"/TransIn"， 补偿操作为url: qsBusi+"/TransInCompensate"
	//Add(qsBusi+"/TransIn", qsBusi+"/TransInCompensate", req)
	// 提交saga事务，dtm会完成所有的子事务/回滚所有的子事务
	saga.WaitResult = true
	saga.RequestTimeout = 1000 * 100
	saga.TimeoutToFail = 1000 * 1000
	saga.RetryLimit = 0
	err = saga.Submit()
	if err != nil {
		log.Logger.Error("分布式事务执行失败,err:", err)
		return "", err
	}
	return locatedAddr, nil
}

// 分块文件进行合并
func FileMpUploadMerge(uploadId, fileHash string) (string, error) {
	srcPath := conf.GetConfig().LocalMpStore + "/" + uploadId
	destPath := conf.GetConfig().LocalStore + "/" + fileHash
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	_, err := util.ExecLinuxShell(cmd)
	if err != nil {
		log.Logger.Error("分块文件合并失败")
		return "", err
	}
	return destPath, nil
}

// 失败的分块文件检查
func CheckFailedMpUploadFile(uploadId string, chunkCount int) ([]int, error) {
	//检查失败的chunkIndex
	chunkArray, err := data.GetFailedChunk(uploadId, chunkCount)
	if err != nil {
		return nil, err
	}
	return chunkArray, nil
}
