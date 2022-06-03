package service

import (
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// CommentAction 登录用户对视频进行评论
func CommentAction(c *gin.Context) {
	var req CommentActionRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("comment action should bind error")
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "comment action should bind error",
			},
			Comment: Comment{},
		})
	}

	// 从解析的token中取出username
	req.UserName = c.GetString("username") //  不过好像没用

	commentDao := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()

	if req.ActionType == 1 { // 发布评论
		commentId, err := commentDao.CreateComment(req.UserId, req.VideoId, &req.CommentText)
		if err != nil {
			c.JSON(http.StatusBadRequest, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "create comment error",
				},
				Comment: Comment{},
			})
			return
		}
		var author User
		if err := userDao.FindUserById(req.UserId, (*repository.User)(&author)); err != nil { // 找作者信息
			c.JSON(http.StatusBadRequest, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "comment author not exists",
				},
				Comment: Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "ok",
			},
			Comment: Comment{
				ID:         commentId,
				User:       author,
				Content:    req.CommentText,
				CreateDate: time.Now().Format("01-02"),
			},
		})
	} else if req.ActionType == 2 { // 删除评论
		err := commentDao.DeleteCommentById(req.CommentId)
		if err != nil {
			c.JSON(http.StatusBadRequest, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "delete comment error",
				},
				Comment: Comment{},
			})
			return
		}
		c.JSON(http.StatusOK, CommentActionResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "ok",
			},
			Comment: Comment{},
		})
	} else {
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "comment action action type error",
			},
			Comment: Comment{},
		})
	}
}

// CommentList 查看视频的所有评论，按发布时间倒序
func CommentList(c *gin.Context) {
	var req CommentListRequest
	if err := c.ShouldBind(&req); err != nil {
		fmt.Println("comment list should bind error")
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "comment list should bind error",
			},
			Comment: Comment{},
		})
		return
	}

	// 创建单例
	commentDAO := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()

	var commentList []repository.Comment // 猜想，如果评论量特别大的话，是不是可以做成分段查询的，how，是不是需要前端来请求
	if err := commentDAO.GetAllComment(&commentList, req.VideoId); err != nil {
		c.JSON(http.StatusBadRequest, CommentActionResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "get comment list error",
			},
			Comment: Comment{},
		})
		return
	}

	// 构造结果
	var author User
	resList := make([]Comment, len(commentList))
	for i, comment := range commentList {
		// 找到评论作者信息
		if err := userDao.FindUserById(comment.AuthorID, (*repository.User)(&author)); err != nil {
			c.JSON(http.StatusBadRequest, CommentActionResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "comment author not exists",
				},
				Comment: Comment{},
			})
		}
		resList[i] = Comment{
			ID:         comment.ID,
			User:       author, // 如果这个能直接成功拿到的话，前面有一个通过id来找人的逻辑就不用写了
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CommentPublishTime, 0).Format("01-02"),
		}
	}

	// 返回结果
	c.JSON(http.StatusOK, CommentListResponse{
		Response: Response{
			StatusCode: 1,
			StatusMsg:  "ok",
		},
		CommentList: &resList,
	})
}
