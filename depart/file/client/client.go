package client

import (
	filePb "fileStore/depart/file/proto"
	userFilePb "fileStore/depart/user-file/proto"
	userPb "fileStore/depart/user/proto"
	etcdClient "fileStore/internel/middleware/etcd/client"
	"fmt"
	"sync"
)

var (
	once           sync.Once
	userClient     userPb.UserServiceClient
	fileClient     filePb.FileServiceClient
	userFileClient userFilePb.UserFileServiceClient
)

func GetUserClient() userPb.UserServiceClient {
	once.Do(func() {
		userConn, err := etcdClient.NewConnection("47.109.159.227:2379", 10000, "user_server")
		if err != nil {
			fmt.Println("userConn err:", err)
			return
		}
		userClient = userPb.NewUserServiceClient(userConn)
	})
	return userClient
}

func GetFileClient() filePb.FileServiceClient {
	once.Do(func() {
		fileConn, err := etcdClient.NewConnection("47.109.159.227:2379", 10000, "file_server")
		if err != nil {
			fmt.Println("fileConn err:", err)
			return
		}
		fileClient = filePb.NewFileServiceClient(fileConn)
	})
	return fileClient
}

func GetUserFileClient() userFilePb.UserFileServiceClient {
	once.Do(func() {
		userFileConn, err := etcdClient.NewConnection("47.109.159.227:2379", 10000, "user_file_server")
		if err != nil {
			fmt.Println("fileConn err:", err)
			return
		}
		userFileClient = userFilePb.NewUserFileServiceClient(userFileConn)
	})
	return userFileClient
}
