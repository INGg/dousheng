package controller

import (
	"demo1/model"
	"demo1/service"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

func Feed(c *gin.Context) {
	var req model.FeedRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.FeedResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "feed should bind error",
			},
			VideoList: nil,
			NextTime:  0,
		})
	}

	if req.LatestTime == 0 {
		req.LatestTime = time.Now().Unix()
	}

	req.FromUserID = c.GetUint("user_id")

	// 调用服务
	feedResponse, err := service.Feed(&req)
	if err != nil {
		zap.L().Error(err.Error())
	} else {
		zap.L().Info("get feed successfully")
	}
	c.JSON(http.StatusOK, &feedResponse)
}
