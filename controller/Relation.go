package controller

import (
	"demo1/model"
	"demo1/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func RelationAction(c *gin.Context) {
	var req model.FollowActionRequest

	if err := c.ShouldBind(&req); err != nil {
		zap.L().Error("ActionRelation should bind error")
		c.JSON(http.StatusOK, &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "actionRelation should bind error",
			},
		})
		return
	}

	// 从token解析的信息中赋值
	req.UserID = c.GetUint("user_id")

	var relationResponse *model.FollowActionResponse
	var err error

	// 调用服务
	if req.ActionType == 1 { // 关注
		relationResponse, err = service.AddRelation(&req)
	} else if req.ActionType == 2 { // 取消关注
		relationResponse, err = service.CancelRelation(&req)
	} else { // 参数错误
		zap.L().Error("relation action action type error")
		c.JSON(http.StatusOK, &model.FollowActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "relation action action type error",
			},
		})
		return
	}

	// 记录日志 返回结果
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("relation action ok")
	}
	c.JSON(http.StatusOK, &relationResponse)
}

func FollowList(c *gin.Context) {
	var req model.UserFollowListRequest

	if err := c.ShouldBind(&req); err != nil {
		// 无法将参数赋值到req中
		c.JSON(200, model.UserFollowListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "无法获取关注列表",
			},
			UserList: nil,
		})
	}

	req.UserID = c.GetUint("user_id")

	var followListResponse *model.UserFollowListResponse
	var err error

	// 调用服务
	followListResponse, err = service.FollowList(&req)

	// 记录日志 返回结果
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("get follow list ok")
	}
	c.JSON(http.StatusOK, &followListResponse)
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	var req model.UserFollowerListRequest

	// 参数错误.无法赋值到req中
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, model.UserFollowerListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "参数错误",
			},
			UserList: nil,
		})
	}

	var followerListResponse *model.UserFollowerListResponse
	var err error

	// 调用服务
	followerListResponse, err = service.FollowerList(&req)

	// 记录日志 返回结果
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("get follow list ok")
	}
	c.JSON(http.StatusOK, &followerListResponse)
}
