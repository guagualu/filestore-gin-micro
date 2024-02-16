package server

import (
	"fileStore/conf"
	pb "fileStore/depart/user/proto"
	"fileStore/depart/user/rpc/server/handle"
	etcdServe "fileStore/internel/middleware/etcd/server"
	"fileStore/log"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
	"net"
	"time"
)

func RpcServer() {
	Serve := grpc.NewServer()
	rg := &handle.UserRpcServiceStruct{}
	pb.RegisterUserServiceServer(
		Serve, rg,
	)
	//etcd 服务发现
	registryConf := etcdServe.RegisterConfig{
		Config: clientv3.Config{ //etcd服务器相关配置
			Endpoints:            []string{conf.GetConfig().EtcdAddr},
			DialTimeout:          time.Duration(3) * time.Second,
			DialKeepAliveTime:    time.Duration(4) * time.Second,
			DialKeepAliveTimeout: time.Duration(5) * time.Second,
			Username:             "root",
			Password:             "meixi253",
			// Logger:               logger,
		},
		ServerName: "user_server",    //微服务的服务名 可以后面使用config
		Address:    "127.0.0.1:8081", //服务集群的真实地址
		Lease:      15,
	}
	reg, err := etcdServe.NewServiceRegister(registryConf)
	if err != nil {
		log.Logger.Error("etcd err:", err)
	}
	defer reg.Close()
	//程序还在运行 就需要一直续约
	go reg.ListenLeaseRespChan()
	lis, err := net.Listen("tcp", ":8081")
	if err != nil {
		fmt.Println("listen err:=", err)
		return
	}
	err = Serve.Serve(lis)
	if err != nil {
		fmt.Println("serve err:=", err)
		return
	}
}
