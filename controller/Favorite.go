package controller

import (
	"demo1/model"
	"demo1/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// FavoriteAction 登录用户对视频的点赞和取消点赞操作
func FavoriteAction(c *gin.Context) {
	var req model.UserFavoriteRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteAction should bind error",
			},
		})
	}

	req.UserID = c.GetUint("user_id")

	favoriteResponse, err := service.FavoriteAction(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("FavoriteAction ok")
	}

	c.JSON(http.StatusOK, favoriteResponse)
}

// FavoriteList 登录用户的所有点赞视频
func FavoriteList(c *gin.Context) {
	var req model.UserFavoriteListRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, model.UserFavoriteListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteList Action should bind error",
			},
			VideoList: nil,
		})
	}

	req.UserID = c.GetUint("user_id")

	favoriteResponse, err := service.FavoriteList(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("get favorite list ok")
	}

	c.JSON(http.StatusOK, favoriteResponse)
}
