package biz

import (
	"encoding/json"
	"fileStore/conf"
	filePb "fileStore/depart/file/proto"
	userFilePb "fileStore/depart/user-file/proto"
	"fileStore/internel/domain"
	"fileStore/internel/middleware/mq"
	"fileStore/internel/pkg/encoding"
	"fileStore/internel/pkg/uuid"
	"fmt"
	dtmgrpc "github.com/dtm-labs/client/dtmgrpc"
	"os"
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
	saga.RequestTimeout = 1000 * 10
	saga.TimeoutToFail = 1000 * 100
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
