package service

import (
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"fmt"
	"go.uber.org/zap"
)

// FavoriteAction 登录用户对视频的点赞和取消点赞操作
func FavoriteAction(req *model.UserFavoriteRequest) (*model.UserFavoriteResponse, error) {
	var user entity.User
	var video entity.Video

	// 创建单例
	userDAO := repository.NewUserDAO()
	videoDAO := repository.NewVideoDAO()

	if err := userDAO.FindUserById(req.UserID, &user); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "User can't find error",
			},
		}, err
	}

	if err := videoDAO.FindVideoById(req.VideoID, &video); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Video can't find error",
			},
		}, err
	}

	// 判断操作类型
	if req.ActionType == 1 { // 点赞
		return addFavorite(req)
	} else if req.ActionType == 2 { // 取消点赞
		return cancelFavorite(req)
	} else {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Favorite action error",
			},
		}, nil
	}
}

func addFavorite(req *model.UserFavoriteRequest) (*model.UserFavoriteResponse, error) {

	// 创建单例
	favoriteDAO := repository.NewFavoriteDAO()

	//点赞操作

	// 是否已点赞
	if favoriteDAO.CheckIsFavorite(req.UserID, req.VideoID) == true {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "has favorite this video",
			},
		}, nil
	}

	//将该视频加入用户的点赞列表
	if err := favoriteDAO.Favorite(req.UserID, req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Add into favorite list error",
			},
		}, err
	}
	//video.FavoriteCount++
	if err := favoriteDAO.AddFavoriteCount(req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Add favorite count error",
			},
		}, err
	}

	return &model.UserFavoriteResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "Favorite successful",
		},
	}, nil
}

func cancelFavorite(req *model.UserFavoriteRequest) (*model.UserFavoriteResponse, error) {
	// 创建单例
	favoriteDAO := repository.NewFavoriteDAO()

	//取消点赞
	// 检查是否已点赞
	if favoriteDAO.CheckIsFavorite(req.UserID, req.VideoID) == false {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "has not favorite this video",
			},
		}, nil
	}

	//将该视频从用户的点赞列表移除
	if err := favoriteDAO.UnFavorite(req.UserID, req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Delete from favorite list error",
			},
		}, err
	}
	// 视频的点赞数量减少
	if err := favoriteDAO.ReduceFavoriteCount(req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Reduce favorite count error",
			},
		}, err
	}

	// 成功
	return &model.UserFavoriteResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "UnFavorite successful",
		},
	}, nil
}

// FavoriteList 登录用户的所有点赞视频
func FavoriteList(req *model.UserFavoriteListRequest) (*model.UserFavoriteListResponse, error) {
	// 创建单例
	favoriteDAO := repository.NewFavoriteDAO()
	videoDAO := repository.NewVideoDAO()
	userDAO := repository.NewUserDAO()
	relationDAO := repository.NewRelationDAO()

	// 结构数组
	var vids []uint

	if err := favoriteDAO.FindFavoriteVideoByUid(req.UserID, &vids); err != nil {
		return &model.UserFavoriteListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Can't find favorite video list by uid",
			},
			VideoList: nil,
		}, err
	}

	if len(vids) == 0 {
		return &model.UserFavoriteListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, but list is nil",
			},
			VideoList: nil,
		}, nil
	}

	videoList := make([]entity.Video, len(vids))

	for i := 0; i < len(videoList); i++ {
		// 找一下视频信息
		if err := videoDAO.FindVideoById(vids[i], &videoList[i]); err != nil {
			zap.L().Error(fmt.Sprintf("vid:%v can't find", videoList[i].ID))
		}

		// 找一下作者信息
		if err := userDAO.FindUserById(videoList[i].AuthorID, &videoList[i].Author); err != nil {
			zap.L().Error(fmt.Sprintf("uid:%v can't find", videoList[i].AuthorID))
		}

		// 请求的人是否关注了作者
		if req.FromUserID != 0 {
			videoList[i].Author.IsFollow = relationDAO.QueryAFollowB(req.FromUserID, videoList[i].AuthorID)

			// 请求的人有没有给视频点赞
			videoList[i].IsFavorite = favoriteDAO.CheckIsFavorite(req.FromUserID, videoList[i].ID)
		}
	}

	return &model.UserFavoriteListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "Get favorite list success",
		},
		VideoList: &videoList,
	}, nil
}
