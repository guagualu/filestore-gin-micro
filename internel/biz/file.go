package biz

import (
	"context"
	"encoding/json"
	"fileStore/conf"
	"fileStore/depart/file/client"
	filePb "fileStore/depart/file/proto"
	userFilePb "fileStore/depart/user-file/proto"
	userPb "fileStore/depart/user/proto"
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
	"time"
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
	saga.RequestTimeout = 1000 * 1000
	saga.TimeoutToFail = 1000 * 10000
	saga.Concurrent = false
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

func FastUpload(ctx context.Context, fileHash string, fileName string, userUuid string) error {
	//查找file表
	_, err := data.GetFileByFileHash(ctx, fileHash)
	if err != nil {
		return err
	}
	userFileClient := client.GetUserFileClient()
	//查找关系 如果已经存在就不再重复存入user-file
	file, _ := userFileClient.GetUserFile(ctx, &userFilePb.UserFileReq{
		FileHash: fileHash,
		FileName: fileName,
		UserUuid: userUuid,
	})
	if file != nil {
		return nil
	}
	//存入user-file中
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
	//if chunkIndex == 1 || chunkIndex == 3 {
	//	return errors.New("测试")
	//}
	// 存入本地 注意 要先存入本地 再写redis 因为有个completed检查可能会在响应未完成时就来请求redis
	locatedAt := conf.GetConfig().LocalMpStore + "/" + uploadId + "/" + "chunk_" + strconv.Itoa(chunkIndex)
	os.MkdirAll(conf.GetConfig().LocalMpStore+"/"+uploadId, 0666)
	localFile, err := os.Create(locatedAt)
	if err != nil {
		return err
	}
	defer localFile.Close()
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
	saga.Concurrent = false
	err = saga.Submit()
	if err != nil {
		log.Logger.Error("分布式事务执行失败,err:", err)
		return "", err
	}
	//针对后端程序的地址 上面的是针对脚本的相对地址
	backAddr := "./static/tmp/" + fileHash
	return backAddr, nil
}

// 分块文件进行合并
func FileMpUploadMerge(uploadId, fileHash string) (string, error) {
	srcPath := conf.GetConfig().LocalMpStore + "/" + uploadId
	destPath := "../../tmp/" + fileHash
	cmd := fmt.Sprintf("cd %s && ls | sort -n | xargs cat > %s", srcPath, destPath)
	//_, err := util.ExecLinuxShell(cmd)  linux使用
	_, _, err := util.ExecWinShell("merge.sh", cmd) //windows系统使用
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

func GetFileInfo(ctx context.Context, fileHash string) (*domain.File, error) {
	return data.GetFileByFileHash(ctx, fileHash)
}

func GetFileInfoList(ctx context.Context, fileHashs []string) ([]*domain.File, error) {
	return data.ListFileInfo(ctx, fileHashs)
}

func GetUserFileInfo(ctx context.Context, fileHash, userUuid, fileName string) (*domain.UserFile, error) {
	res, err := client.GetUserFileClient().GetUserFile(ctx, &userFilePb.UserFileReq{
		FileHash: fileHash,
		FileName: fileName,
		UserUuid: userUuid,
	})
	if err != nil {
		return nil, err
	}
	createdAt, _ := time.Parse("2006-01-02 15:04:05", res.CreatedAt)
	updatedAt, _ := time.Parse("2006-01-02 15:04:05", res.UpdatedAt)
	return &domain.UserFile{
		ID:       uint(res.Id),
		FileHash: res.FileHash,
		UserUuid: res.UserUuid,
		FileName: res.FileName,
		CreateAt: createdAt,
		UpdateAt: updatedAt,
	}, nil
}

func FileServerGetUserInfo(ctx context.Context, userUuid string) (*domain.User, error) {
	info, err := client.GetUserClient().UserInfo(ctx, &userPb.ReqUserInfo{
		Uuid: userUuid,
	})
	if err != nil {
		return nil, err
	}
	res := domain.User{
		Id:       uint(info.Id),
		Uuid:     info.Uuid,
		NickName: info.NickName,
		Email:    info.Email,
		Mobile:   info.Mobile,
	}
	return &res, nil
}
