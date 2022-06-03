package service

import (
	"demo1/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

//  ---favourite---

func FavoriteAction(c *gin.Context) {
	var req UserFavoriteRequest
	var user repository.User
	var video repository.Video

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteAction should bind error",
			},
		})
	}

	if err := repository.FindUserByToken(req.Token, &user); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "User can't find error",
			},
		})
	}

	if err := repository.FindVideoById(req.VideoId, &video); err != nil {
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Video can't find error",
			},
		})
	}

	if req.ActionType == 1 {
		//点赞操作
		video.FavoriteCount++
		//将该视频加入用户的点赞列表
		if err := repository.ChangeFavorite(req.UserId, req.VideoId, 1); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Add into favorite list error",
				},
			})
		}
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Favorite successful",
			},
		})
	} else {
		//取消点赞
		video.FavoriteCount--
		//将该视频从用户的点赞列表移除
		if err := repository.ChangeFavorite(req.UserId, req.VideoId, 2); err != nil {
			c.JSON(http.StatusOK, UserFavoriteResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "Delete from favorite list error",
				},
			})
		}
		c.JSON(http.StatusOK, UserFavoriteResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "UnFavorite successful",
			},
		})
	}
}

//  ---favourite list---

func FavoriteList(c *gin.Context) {
	var req UserFavoriteListRequest
	var videoList []repository.Video

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "FavoriteList Action should bind error",
			},
			VideoList: nil,
		})
	}
	if err := repository.FindFavoriteVideoByUid(req.UserId, &videoList); err != nil {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "Can't find favorite video by uid",
			},
			VideoList: nil,
		})
	} else {
		c.JSON(http.StatusOK, UserFavoriteListResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "Get favorite list success",
			},
			VideoList: videoList,
		})
	}
}
