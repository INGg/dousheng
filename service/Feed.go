package service

import (
	"demo1/middleware"
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// ---Feed---

//type repository.Video struct {
//	repository.Video  repository.repository.Video
//	Author repository.User
//}

func Feed(c *gin.Context) {
	// 把请求数据拿出来
	var req FeedRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, FeedResponse{
			StatusCode: 1,
			StatusMsg:  "feed should bind error",
			VideoList:  nil,
			NextTime:   0,
		})
	}

	if req.LatestTime == 0 {
		req.LatestTime = time.Now().Unix()
	}

	// 解析token
	if req.Token != "" {
		_, err := middleware.ParseToken(req.Token)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				StatusCode: 1,
				StatusMsg:  "token error",
				VideoList:  nil,
				NextTime:   0,
			})
			return
		}
	}

	fmt.Printf("%+v\n", req)

	// 获取10条Video列表
	var videoList = make([]repository.Video, 32)
	err := repository.GetVideoList(&videoList, 32, req.LatestTime)

	var resList = make([]Video, len(videoList))

	// 给获取到的video加上作者信息和是否对这个视频点赞了
	for i, video := range videoList {
		// 加上作者信息
		if err := repository.FindUserById(video.AuthorID, &videoList[i].Author); err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				StatusCode: 1,
				StatusMsg:  "get author error",
				VideoList:  nil,
				NextTime:   0,
			})
			return
		}

		resList[i].Video = video

		resList[i].IsFavorite = repository.CheckIsFavorite(videoList[i].AuthorID, video.ID)

		fmt.Printf("%+v\n", resList[i])
	}

	if err != nil {
		panic("get video error")
	}

	// 返回结果
	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: 0,
		StatusMsg:  "ok",
		VideoList:  resList,
		NextTime:   videoList[0].PublishTime,
	})
}
