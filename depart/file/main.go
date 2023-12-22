package main

import (
	"fileStore/depart/file/router"
	rpcServe "fileStore/depart/file/rpc/server"
	"fileStore/internel/domain"
	"github.com/gin-gonic/gin"
)

func main() {
	domain.ServiceName = "file"
	r := gin.Default()
	router.Router(r)
	go rpcServe.RpcServer()
	r.Run(":8082") // 启动 HTTP 服务器
}
