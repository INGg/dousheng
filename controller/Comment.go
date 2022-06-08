package controller

import (
	"demo1/model"
	"demo1/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func CommentAction(c *gin.Context) {
	var req model.CommentActionRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("comment action should bind error")
		c.JSON(http.StatusBadRequest, model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "comment action should bind error",
			},
			Comment: model.Comment{},
		})
	}

	// 从解析的token中取出username
	req.UserName = c.GetString("username") //  不过好像没用
	req.UserID = c.GetUint("user_id")

	// 调用服务
	var commentResponse *model.CommentActionResponse
	var err error

	if req.ActionType == 1 { // 发布评论
		commentResponse, err = service.AddComment(&req)
	} else if req.ActionType == 2 { // 删除评论
		commentResponse, err = service.DeleteComment(&req)
	} else {
		zap.L().Error("comment action action type error")
		c.JSON(http.StatusBadRequest, model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "comment action action type error",
			},
			Comment: model.Comment{},
		})
		return
	}

	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("comment action ok")
	}
	c.JSON(http.StatusOK, &commentResponse)
}

func CommentList(c *gin.Context) {
	var req model.CommentListRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("comment list should bind error")
		c.JSON(http.StatusBadRequest, model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "comment list should bind error",
			},
			Comment: model.Comment{},
		})
		return
	}

	req.UserName = c.GetString("username")
	req.UserID = c.GetUint("user_id")

	// 调用服务
	listResponse, err := service.CommentList(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info(fmt.Sprintf("video id %v get comment list successfully", req.VideoID))
	}
	c.JSON(http.StatusOK, &listResponse)
}
