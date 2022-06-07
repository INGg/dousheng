package service

import (
	"demo1/model"
	"demo1/repository"
	"errors"
	"fmt"
	"go.uber.org/zap"
)

// FavoriteAction 登录用户对视频的点赞和取消点赞操作
func FavoriteAction(req *model.UserFavoriteRequest) (*model.UserFavoriteResponse, error) {
	var user repository.User
	var video repository.Video

	// 创建单例
	userDAO := repository.NewUserDAO()
	videoDAO := repository.NewVideoDAO()

	if err := userDAO.FindUserById(req.UserID, &user); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "User can't find error",
			},
		}, errors.New("FavoriteAction User can't find error")
	}

	if err := videoDAO.FindVideoById(req.VideoID, &video); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Video can't find error",
			},
		}, errors.New("FavoriteAction Video can't find error")
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
	//将该视频加入用户的点赞列表
	if err := favoriteDAO.Favorite(req.UserID, req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Add into favorite list error",
			},
		}, errors.New("add into favorite list error")
	}
	//video.FavoriteCount++
	if err := favoriteDAO.AddFavoriteCount(req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Add favorite count error",
			},
		}, errors.New("add favorite count error")
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
	//将该视频从用户的点赞列表移除
	if err := favoriteDAO.UnFavorite(req.UserID, req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Delete from favorite list error",
			},
		}, errors.New("delete from favorite list error")
	}
	// 视频的点赞数量减少
	if err := favoriteDAO.ReduceFavoriteCount(req.VideoID); err != nil {
		return &model.UserFavoriteResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Reduce favorite count error",
			},
		}, errors.New("reduce favorite count error")
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

	// 结构数组
	var videoList []repository.Video

	if err := favoriteDAO.FindFavoriteVideoByUid(req.UserID, &videoList); err != nil {
		return &model.UserFavoriteListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "Can't find favorite video list by uid",
			},
			VideoList: nil,
		}, errors.New("can't find favorite video list by uid")
	}

	for i := 0; i < len(videoList); i++ {
		// 找一下视频信息
		if err := videoDAO.FindVideoById(videoList[i].ID, &videoList[i]); err != nil {
			zap.L().Error(fmt.Sprintf("vid:%v can't find", videoList[i].ID))
		}

		// 找一下作者信息
		if err := userDAO.FindUserById(videoList[i].AuthorID, (*repository.User)(&videoList[i].Author)); err != nil {
			zap.L().Error(fmt.Sprintf("uid:%v can't find", videoList[i].AuthorID))
		}
	}

	return &model.UserFavoriteListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "Get favorite list success",
		},
		VideoList: videoList,
	}, nil
}
