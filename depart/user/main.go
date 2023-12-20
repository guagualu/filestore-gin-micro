package main

import (
	"fileStore/depart/user/router"
	rpcServe "fileStore/depart/user/rpc/server"
	"fileStore/internel/domain"
	"fmt"
	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("hello,world")
	domain.ServiceName = "user"
	r := gin.Default()
	router.Router(r)
	go rpcServe.RpcServer()
	r.Run() // 启动 HTTP 服务器
}
