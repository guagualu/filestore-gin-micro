package main

import (
	"context"
	userPb "fileStore/depart/user/proto"
	etcdClient "fileStore/internel/middleware/etcd/client"
	"fmt"
)

func main() {
	userConn, err := etcdClient.NewConnection("47.109.159.227:2379", 10000, "user_server")
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
