package service

import (
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

type PublishActionRequest struct {
	Token string                `json:"token" form:"token"`
	Data  *multipart.FileHeader `json:"data" form:"data"`
	Title string                `json:"title" form:"title"`
}

type PublishActionResponse struct {
	Response
}

func Publish(c *gin.Context) {
	var req PublishActionRequest
	// 绑定结构体
	if err := c.ShouldBind(&req); err == nil {
		if err != nil {
			c.JSON(http.StatusOK, PublishActionResponse{Response{
				StatusCode: 1,
				StatusMsg:  "get video file error",
			}})
			return
		}
		// 保存文件到文件夹
		finalName := filepath.Base(req.Data.Filename)
		saveFile := filepath.Join("./static/", finalName)

		fmt.Println("上传文件成功")

		if err := c.SaveUploadedFile(req.Data, saveFile); err != nil { // 保存失败返回失败信息
			c.JSON(http.StatusOK, PublishActionResponse{
				Response{
					StatusCode: 1,
					StatusMsg:  "video save error",
				},
			})
			return
		}

		// 文件信息写入数据库
		if err := repository.InsertVideo(req.Token, saveFile, req.Title); err != nil {
			c.JSON(http.StatusOK, PublishActionResponse{
				Response{
					StatusCode: 1,
					StatusMsg:  "video info insert database error",
				},
			})
			return
		}

		// 成功返回
		c.JSON(http.StatusOK, PublishActionResponse{
			Response{
				StatusCode: 0,
				StatusMsg:  finalName + " uploaded successfully",
			},
		})

	} else {
		c.JSON(http.StatusOK, PublishActionResponse{
			Response{
				StatusCode: 1,
				StatusMsg:  "publish should bind error",
			},
		})
	}
}

// ---PublishList---

type Video struct {
	repository.Video
	IsFavorite bool
}

type PublishListRequest struct {
	Token  string `json:"token" form:"token"`
	UserId uint   `json:"user_id" form:"user_id"`
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

func PublishList(c *gin.Context) {
	var req PublishListRequest
	if err := c.ShouldBind(&req); err == nil {

		var videoList []repository.Video

		if err := repository.FindAllVideoByUid(req.UserId, &videoList); err != nil {
			c.JSON(http.StatusOK, PublishListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "get published list error",
				},
				VideoList: nil,
			})
		}

		resList := make([]Video, len(videoList))

		for i, video := range videoList {
			fmt.Println(i, video)
		}

		for i, video := range videoList {
			resList[i].Video = video
			repository.FindUserById(video.AuthorID, &resList[i].Video.Author)
			resList[i].IsFavorite = repository.CheckIsFavorite(req.UserId, video.ID)
		}

		for i, video := range resList {
			fmt.Println(i, video)
		}

		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "get published list successfully",
			},
			VideoList: resList,
		})

	} else {
		c.JSON(http.StatusOK, PublishListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "get published list should bind error",
			},
			VideoList: nil,
		})
	}
}
