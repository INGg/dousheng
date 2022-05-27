package service

import (
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedRequest struct {
	LatestTime int64  `json:"latest_time,omitempty"`
	Token      string `json:"token,omitempty"`
}

type FeedResponse struct {
	StatusCode int32              `json:"status_code,omitempty"`
	StatusMsg  string             `json:"status_msg,omitempty"`
	VideoList  []repository.Video `json:"video_list,omitempty"`
	NextTime   int64              `json:"next_time,omitempty"`
}

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
	var video_list = make([]repository.Video, 32)
	err = repository.GetVideoList(&video_list, 32, req.LatestTime)

	// 给获取到的video加上作者信息
	for i, video := range video_list {
		err := repository.FindUserById(video.AuthorID, &video_list[i].Author)
		fmt.Println(video_list[i].Author)
		if err != nil {
			c.JSON(http.StatusOK, FeedResponse{
				StatusCode: 1,
				StatusMsg:  "get author error",
				VideoList:  nil,
				NextTime:   0,
			})
			return
		}
	}

	if err != nil {
		panic("get video error")
	}

	//sort.Slice(video_list, func(i, j int) bool {
	//	if video_list[i].PublishTime < video_list[j].PublishTime {
	//		return true
	//	}
	//	return false
	//})

	// 返回结果
	c.JSON(http.StatusOK, FeedResponse{
		StatusCode: 0,
		StatusMsg:  "ok",
		VideoList:  video_list,
		NextTime:   video_list[0].PublishTime,
	})
}
