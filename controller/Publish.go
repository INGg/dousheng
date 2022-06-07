package controller

import (
	"demo1/model"
	"demo1/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"path/filepath"
)

func Publish(c *gin.Context) {
	var req model.PublishActionRequest

	// 绑定结构体
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("publish should bind error")
		fmt.Println(err.Error())
		c.JSON(http.StatusOK, model.PublishActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "publish should bind error",
			},
		})
	}

	req.UserName = c.GetString("username")
	req.UserID = c.GetUint("user_id")

	//fmt.Printf("%+v\n", req)

	// 保存文件到文件夹
	//finalName := filepath.Base(req.Data.Filename)
	////saveFile := filepath.Join("./static/", finalName)
	//saveFile := "./static/" + finalName

	file := req.Data
	suffix := filepath.Ext(file.Filename) //得到后缀
	//path := "./static/" + file.Filename
	fileName := service.NewFileName(req.UserID)
	path := "./static/" + fileName + suffix

	//fmt.Println(finalName, saveFile)
	//fmt.Println("上传文件成功")

	if err := c.SaveUploadedFile(file, path); err != nil { // 保存失败返回失败信息
		fmt.Println("video save error")
		c.JSON(http.StatusOK, model.PublishActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "video save error",
			},
		})
		return
	}

	actionResponse, err := service.Publish(&req, file, fileName)
	if err != nil {
		zap.L().Error("publish video error")
	} else {
		zap.L().Info("publish successfully")
	}
	c.JSON(http.StatusOK, &actionResponse)
}

func PublishList(c *gin.Context) {
	var req model.PublishListRequest
	if err := c.ShouldBind(&req); err != nil {

		fmt.Println("get published list should bind error")
		c.JSON(http.StatusOK, model.PublishListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get published list should bind error",
			},
			VideoList: nil,
		})
	}

	req.UserName = c.GetString("username")

	zap.L().Info(fmt.Sprintf("PublishListRequest user id : %+v\n", req.UserId))

	// 调用服务
	listResponse, err := service.PublishList(&req)
	if err != nil {
		zap.L().Error("get published list error")
	} else {
		zap.L().Info("get published list successfully")
	}
	c.JSON(http.StatusOK, &listResponse)
}
