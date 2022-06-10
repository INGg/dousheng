package main

import (
	"demo1/middleware"
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	// 初始化
	err := middleware.InitLogger()
	defer middleware.CloseLogger()
	if err != nil {
		fmt.Printf("init logger failed, err:%v\n", err)
		return
	}

	//gin.SetMode("release") // release模式，在生产环境下使用
	gin.SetMode("debug") // debug模式，在debug环境下使用

	// 连接数据库
	repository.InitDb()
	// TODO 检查static文件夹是否存在，不存在就创建

	// 准备启动gin引擎
	r := gin.Default()

	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 处理静态文件
	r.StaticFS("/static", http.Dir("./static"))

	// 启动路由
	initRouter(r)

	r.Run(":9090")
}
