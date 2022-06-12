package service

import (
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"demo1/util"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"mime/multipart"
)

// ---Publish---

func Publish(req *model.PublishActionRequest, file *multipart.FileHeader, fileName string) (*model.PublishActionResponse, error) {

	// 创建单例
	videoDAO := repository.NewVideoDAO()

	// 构造文件信息:
	// 生成url
	// 视频url
	playUrl := util.GetFileUrl(fileName + util.GetDefaultVideoSuffix())

	// 封面url
	coverUrl := util.GetFileUrl(fileName + util.GetDefaultImageSuffix())
	// 生成并存储封面
	if err := util.SaveImageFromVideo(fileName, true); err != nil {
		// 没有生成封面使用默认的
		coverUrl = "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	}

	// 文件信息写入数据库
	if err := videoDAO.InsertVideo(req.UserID, playUrl, coverUrl, req.Title); err != nil {
		return &model.PublishActionResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "video info insert database error",
			},
		}, errors.New("video info insert database error")
	}

	// 成功返回
	zap.L().Info(file.Filename + " uploaded successfully")
	return &model.PublishActionResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  file.Filename + " uploaded successfully",
		},
	}, nil

}

// ---PublishList---

func PublishList(req *model.PublishListRequest) (*model.PublishListResponse, error) {
	// 创建单例
	userDAO := repository.NewUserDAO()
	videoDAO := repository.NewVideoDAO()
	favoriteDAO := repository.NewFavoriteDAO()
	relationDAO := repository.NewRelationDAO()

	var videoList []entity.Video

	if err := videoDAO.FindAllVideoByUid(req.UserID, &videoList); err != nil {
		return &model.PublishListResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get published list error",
			},
			VideoList: nil,
		}, err
	}

	if len(videoList) == 0 { // 该用户没发过视频
		return &model.PublishListResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "ok, publish list is nil",
			},
			VideoList: nil,
		}, nil
	}

	resList := make([]model.Video, len(videoList))

	for i, video := range videoList {
		resList[i].Video = video
		if err := userDAO.FindUserById(video.AuthorID, &resList[i].Video.Author); err != nil {
			return &model.PublishListResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "get published list error",
				},
				VideoList: nil,
			}, nil
		}
		if req.FromUserID != 0 {
			resList[i].Video.Author.IsFollow = relationDAO.QueryAFollowB(req.FromUserID, resList[i].AuthorID)

			resList[i].Video.IsFavorite = favoriteDAO.CheckIsFavorite(req.FromUserID, video.ID)
		}
	}

	//for i, video := range resList {
	//	fmt.Printf("resList[%v]: %+v\n", i, video)
	//}

	return &model.PublishListResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "get published list successfully",
		},
		VideoList: &resList,
	}, nil

}

// NewFileName 根据userId+用户发布的视频数量连接成独一无二的文件名
func NewFileName(userId uint) string {
	var count int64

	err := repository.NewVideoDAO().QueryVideoCountByUid(userId, &count)
	if err != nil {
		log.Println(err)
	}
	return fmt.Sprintf("%d-%d", userId, count)
}
