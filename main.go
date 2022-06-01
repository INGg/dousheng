package main

import (
	"demo1/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化
	// 连接数据库
	repository.InitDb()
	// TODO 检查static文件夹是否存在，不存在就创建

	// 准备启动gin引擎
	r := gin.Default()

	// 处理静态文件
	r.StaticFS("/static", http.Dir("./static"))

	// 启动路由
	initRouter(r)

	r.Run(":9090")
}
