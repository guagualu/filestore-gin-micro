package main

import (
	"context"
	"fileStore/conf"
	userPb "fileStore/depart/user/proto"
	etcdClient "fileStore/internel/middleware/etcd/client"
	"fmt"
)

func main() {
	userConn, err := etcdClient.NewConnection(conf.GetConfig().EtcdAddr, 10000, "user_server")
	if err != nil {
		fmt.Println("userConn err:", err)
		return
	}
	userClient := userPb.NewUserServiceClient(userConn)
	res, err := userClient.UserInfo(context.Background(), &userPb.ReqUserInfo{Uuid: "asdasd"})
	if err != nil {
		fmt.Println("grpc req err :", err)
		return
	}
	fmt.Println(res)

}
