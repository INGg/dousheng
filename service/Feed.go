package service

import (
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"fmt"
	"go.uber.org/zap"
)

// ---Feed---

//type repository.Video struct {
//	repository.Video  repository.repository.Video
//	Author repository.User
//}

func Feed(req *model.FeedRequest) (*model.FeedResponse, error) {

	//fmt.Printf("%+v\n", req)

	// 创建单例
	userDAO := repository.NewUserDAO()
	videoDAO := repository.NewVideoDAO()
	favoriteDAO := repository.NewFavoriteDAO()
	relationDAO := repository.NewRelationDAO()

	// 获取10条Video列表
	var videoList = make([]entity.Video, 32)
	err := videoDAO.GetVideoList(&videoList, 30, req.LatestTime)

	if len(videoList) == 0 {
		return &model.FeedResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, but feed list is nil",
			},
			VideoList: nil,
			NextTime:  0,
		}, nil
	}

	var resList = make([]model.Video, len(videoList))

	// 给获取到的video加上作者信息和是否对这个视频点赞了
	for i, video := range videoList {
		// 加上作者信息
		if err := userDAO.FindUserById(video.AuthorID, &videoList[i].Author); err != nil {
			return &model.FeedResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "get author error",
				},
				VideoList: nil,
				NextTime:  0,
			}, err
		}

		resList[i].Video = video
		resList[i].Author = videoList[i].Author

		// 查询是否已经登陆
		if req.FromUserID != 0 {
			// 查询发起请求用户是否关注了这个人
			resList[i].Author.IsFollow = relationDAO.QueryAFollowB(req.FromUserID, resList[i].AuthorID)

			// 查询发起请求用户是否给这个视频点赞了
			resList[i].Video.IsFavorite = favoriteDAO.CheckIsFavorite(req.FromUserID, video.ID)
		}

		fmt.Printf("%+v\n", resList[i])
	}

	if err != nil {
		zap.L().Error("get video error")
		return &model.FeedResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get video error",
			},
			VideoList: nil,
			NextTime:  0,
		}, err
	}

	// 返回结果
	return &model.FeedResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "ok",
		},
		VideoList: &resList,
		NextTime:  videoList[0].PublishTime,
	}, nil
}
