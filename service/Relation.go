package service

import (
	"demo1/middleware"
	"demo1/repository"
	"github.com/gin-gonic/gin"
	"net/http"
)

//单例模式
var relationDao = repository.NewRelationDAO()
var userDao = repository.NewUserDAO()

// RelationAction 关注操作
func RelationAction(c *gin.Context) {
	var req FollowActionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, Response{
			StatusCode: 1,
			StatusMsg:  "参数解析失败",
		})

	} else {
		// 判断Token
		if req.Token != "" {
			_, err := middleware.ParseToken(req.Token)
			if err != nil {
				c.JSON(http.StatusOK, Response{
					StatusCode: 1,
					StatusMsg:  "token error",
				})
				return
			}
		}
		if err := relationDao.AddRelation(req.UserId, req.ToUserId); err != nil {
			c.JSON(200, Response{
				StatusCode: 1,
				StatusMsg:  "关注失败",
			})
		}
		c.JSON(200, Response{
			StatusCode: 0,
			StatusMsg:  "Success",
		})
	}

	// TODO
}

// FollowList 获取关注列表
func FollowList(c *gin.Context) {
	var req UserFollowerListRequest
	var FollowList []repository.User
	UserD := repository.NewUserDAO()
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, UserFollowListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "无法获取关注列表",
			},
			userList: nil,
		})
	} else {
		// 判断token
		if req.Token != "" {
			_, err := middleware.ParseToken(req.Token)
			if err != nil {
				c.JSON(http.StatusOK, FeedResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "",
					},
					VideoList: nil,
					NextTime:  0,
				})
				return
			}
		}
		FollowId := relationDao.QueryFollowIdByUserId(req.UserId)
		if FollowId == nil {
			c.JSON(200, UserFollowListResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "你还没有关注呢",
				},
				userList: nil,
			})
		} else {
			err = UserD.FindMUserByIdList(FollowId, &FollowList)
			if err != nil {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "用户查询出错",
					},
					userList: nil,
				})
			} else {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "Success",
					},
					userList: FollowList,
				})
			}

		}
	}
	// TODO
}

// FollowerList 获取粉丝列表
func FollowerList(c *gin.Context) {
	var req UserFollowerListRequest
	var FollowerList []repository.User
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(200, UserFollowerListResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "参数错误",
			},
			userList: nil,
		})
	} else {
		if req.Token != "" {
			_, err := middleware.ParseToken(req.Token)
			if err != nil {
				c.JSON(http.StatusOK, UserFollowListResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "Token is err",
					},
					userList: nil,
				})
				return
			}
		}
		FollowerId := relationDao.QueryFollowIdByUserId(req.UserId)
		if FollowerId == nil {
			c.JSON(200, UserFollowerListResponse{
				Response: Response{
					StatusCode: 0,
					StatusMsg:  "你还没有粉丝呢",
				},
				userList: nil,
			})

		} else {
			if err := userDao.FindMUserByIdList(FollowerId, &FollowerList); err != nil {
				c.JSON(200, UserFollowListResponse{
					Response: Response{
						StatusCode: 1,
						StatusMsg:  "用户查询出错",
					},
					userList: nil,
				})
			} else {
				c.JSON(200, UserFollowerListResponse{
					Response: Response{
						StatusCode: 0,
						StatusMsg:  "success",
					},
					userList: FollowerList,
				})
			}

		}
	}

	// TODO
}
