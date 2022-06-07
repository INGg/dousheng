package controller

import (
	"demo1/model"
	"demo1/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func Register(c *gin.Context) {
	// 1.拿去请求信息
	var req model.UserRegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.UserResisterResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "login should bind error",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	// 调用服务
	resisterResponse, err := service.Register(&req)
	if err != nil {
		zap.L().Error("register error")
	} else {
		zap.L().Info(fmt.Sprintf("%v register successfully", req.UserName))
	}

	c.JSON(http.StatusOK, &resisterResponse)
}

func Login(c *gin.Context) {

	// 1.获取请求参数
	var req model.UserLoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.UserLoginResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "login should bind error",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	// 调用服务
	loginResponse, err := service.Login(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info(fmt.Sprintf("%v login successfully", req.UserName))
	}
	c.JSON(http.StatusOK, &loginResponse)
}

func UserInfo(c *gin.Context) {
	var req model.UserInfoRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.UserInfoResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "UserInfo should bind error",
			},
			User: model.User{},
		})
		return
	}
	req.UserName = c.GetString("username")

	// 调用服务
	infoResponse, err := service.UserInfo(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info(fmt.Sprintf("get uid : %v user info", req.UserId))
	}
	c.JSON(http.StatusOK, &infoResponse)
}
