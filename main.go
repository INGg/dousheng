package main

import (
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

var IP string

func inputIp() string {
	fmt.Scan(&IP)
	return IP
}

func main() {
	// 连接数据库
	repository.InitDb()

	// 启动gin引擎
	r := gin.Default()

	// 处理静态文件
	//dir, _ := os.Getwd()
	r.StaticFS("/static", http.Dir("./static"))

	// 启动路由
	initRouter(r)

	r.Run(":9090")
}
