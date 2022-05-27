package service

import (
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
	err := c.BindJSON(&req)
	if err != nil {
		req.LatestTime = time.Now().Unix()
	}

	// 获取10条Video列表
	var videoList = make([]repository.Video, 32)
	err = repository.GetVideoList(&videoList, 32, req.LatestTime)

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

		resList[i].IsFavorite = repository.CheckIsFavorite(req.Token, video.ID)

		fmt.Println(resList[i])
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
