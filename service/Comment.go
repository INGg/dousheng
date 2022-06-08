package service

import (
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"time"
)

// AddComment 登录用户对视频进行评论
func AddComment(req *model.CommentActionRequest) (*model.CommentActionResponse, error) {
	// 创建单例
	commentDao := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()
	relationDAO := repository.NewRelationDAO()

	commentId, err := commentDao.CreateComment(req.UserID, req.VideoID, &req.CommentText)
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "create comment error",
			},
			Comment: model.Comment{},
		}, err
	}

	var author entity.User
	if err := userDao.FindUserById(req.UserID, &author); err != nil { // 找作者信息
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "comment author not exists",
			},
			Comment: model.Comment{},
		}, err
	}

	// 请求的人是否关注了评论的人
	author.IsFollow = relationDAO.QueryAFollowB(req.UserID, author.ID)

	return &model.CommentActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		Comment: model.Comment{
			ID:         commentId,
			User:       author,
			Content:    req.CommentText,
			CreateDate: time.Now().Format("01-02"),
		},
	}, nil
}

func DeleteComment(req *model.CommentActionRequest) (*model.CommentActionResponse, error) {
	// 创建单例
	commentDao := repository.NewCommentDAO()

	err := commentDao.DeleteCommentById(req.CommentID)
	if err != nil {
		return &model.CommentActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "delete comment error",
			},
			Comment: model.Comment{},
		}, err
	}

	return &model.CommentActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		Comment: model.Comment{},
	}, err
}

// CommentList 查看视频的所有评论，按发布时间倒序
func CommentList(req *model.CommentListRequest) (*model.CommentListResponse, error) {

	// 创建单例
	commentDAO := repository.NewCommentDAO()
	userDao := repository.NewUserDAO()
	relationDAO := repository.NewRelationDAO()

	var commentList []entity.Comment // 猜想，如果评论量特别大的话，是不是可以做成分段查询的，how，是不是需要前端来请求
	if err := commentDAO.GetAllComment(&commentList, req.VideoID); err != nil {
		return &model.CommentListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get comment list error",
			},
			CommentList: nil,
		}, err
	}

	if len(commentList) == 0 {
		return &model.CommentListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, but comment list is nil",
			},
			CommentList: nil,
		}, nil
	}

	// 构造结果
	var author entity.User
	resList := make([]model.Comment, len(commentList))
	for i, comment := range commentList {
		// 找到评论作者信息
		if err := userDao.FindUserById(comment.AuthorID, &author); err != nil {
			return &model.CommentListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "comment author not exists",
				},
				CommentList: nil,
			}, err
		}

		// 请求的人是否关注了评论的人
		if req.FromUserID != 0 {
			author.IsFollow = relationDAO.QueryAFollowB(req.FromUserID, author.ID)
		}

		resList[i] = model.Comment{
			ID:         comment.ID,
			User:       author, // 如果这个能直接成功拿到的话，前面有一个通过id来找人的逻辑就不用写了
			Content:    comment.Content,
			CreateDate: time.Unix(comment.CommentPublishTime, 0).Format("01-02"),
		}
	}

	//for _, comment := range resList {
	//	fmt.Printf("%+v\n", comment)
	//}

	// 返回结果
	return &model.CommentListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		CommentList: &resList,
	}, nil
}
