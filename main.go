package main

import (
	"demo1/repository"
	"github.com/gin-gonic/gin"
)

func main() {
	// 连接数据库
	repository.InitDb()

	// 启动gin引擎
	r := gin.Default()

	// 启动路由
	initRouter(r)

	r.Run(":9090")
}
