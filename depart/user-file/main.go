package main

import (
	"fileStore/depart/user-file/router"
	rpcServe "fileStore/depart/user-file/rpc/server"
	"fileStore/internel/domain"
	"github.com/gin-gonic/gin"
)

func main() {
	domain.ServiceName = "user-file"
	r := gin.Default()
	router.Router(r)
	go rpcServe.RpcServer()
	r.Run(":8084") // 启动 HTTP 服务器
}
