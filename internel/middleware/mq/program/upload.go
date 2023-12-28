package program

import (
	"context"
	"encoding/json"
	"fileStore/internel/data"
	"fileStore/internel/middleware/mq"
	"fileStore/internel/middleware/oss"
	"fileStore/log"
	"os"
)

func StoreOssPragram(message []byte) error {
	//1、将message 从json转为golang类型 struct
	fileinfo := mq.MqFileInfo{}
	err := json.Unmarshal(message, &fileinfo)
	if err != nil {
		return err
	}
	//2、找到file的临时存储地址，读取文件
	file, err := os.Open(fileinfo.CurLocateAt)
	if err != nil {
		return err
	}
	//filebyte, err := io.ReadAll(file)
	//3、获取ceph连接 并且将文件存储进去
	// 文件写入OSS存储
	ossPath := "oss/" + fileinfo.FileHash + "/" + fileinfo.FileName
	// 判断写入OSS为同步还是异步

	// TODO: 设置oss中的文件名，方便指定文件名下载
	err = oss.Bucket().PutObject(ossPath, file)
	if err != nil {
		log.Logger.Error("oss err:=", err)
		return err
	}
	//4、oss存入成功后，更新file表的locateAt
	err = data.UpdataFileLocated(context.Background(), fileinfo.FileHash, ossPath)
	if err != nil {
		log.Logger.Error("file reStore err")
		return err
	}
	return nil
}
