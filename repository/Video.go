package repository

import (
	"demo1/model/entity"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
	"time"
)

type VideoDAO struct {
}

var (
	videoDAO  *VideoDAO
	videoOnce sync.Once
)

func NewVideoDAO() *VideoDAO {
	videoOnce.Do(func() {
		videoDAO = &VideoDAO{}
	})
	return videoDAO
}

// GetVideoList 获取视频列表给Feed
func (v *VideoDAO) GetVideoList(videoList *[]entity.Video, lim int, ReqTime int64) error {

	res := db.Limit(lim).Order("publish_time desc").Where("publish_time <= ?", ReqTime).Find(videoList)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (v *VideoDAO) FindVideoByPathAndUid(path string, uid int64, video *entity.Video) error {
	if res := db.Model(entity.Video{}).Where("play_url = ? AND author_id = ?", path, uid).First(video); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// InsertVideo 将视频信息写入数据库，返回错误信息和错误
func (v *VideoDAO) InsertVideo(uid uint, playUrl string, coverUrl string, title string) error {
	// 构造video
	var author entity.User
	if err := NewUserDAO().FindUserById(uid, &author); err != nil { // 找到作者
		return err
	}

	// 根据token获取视频上传者
	fmt.Println(uid, playUrl, title)

	// 存入数据库
	res := db.Create(&entity.Video{
		PublishTime:   time.Now().Unix(),
		Author:        entity.User(author),
		AuthorID:      uid,
		PlayUrl:       playUrl,
		CoverUrl:      coverUrl,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("create video error")
		return res.Error
	}

	fmt.Println("insert video info ok")
	return nil
}

//func (v *VideoDAO) VideoCount() int64 {
//	var count int64
//	db.Model(&entity.Video{}).Count(&count)
//	return count
//}
//
//func (v *VideoDAO) VideoCount() int64 {
//	var count int64
//	db.Model(&entity.Video{}).Count(&count)
//	return count
//}

// FindAllVideoByUid 通过uid找到这个人发布的所有视频
func (v *VideoDAO) FindAllVideoByUid(uid uint, VideoList *[]entity.Video) error {
	res := db.Model(&entity.Video{}).Where("author_id = ?", uid).Find(VideoList)
	if res.Error != nil {
		return res.Error
	}
	for _, video := range *VideoList {
		fmt.Printf("%+v\n", video)
	}
	return nil
}

// FindVideoById 通过id找到Video
func (v *VideoDAO) FindVideoById(vid uint, video *entity.Video) error {
	if res := db.Model(&entity.Video{}).Where("id = ?", vid).First(video); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find video error")
		return res.Error
	}
	return nil
}

func (v *VideoDAO) QueryVideoCountByUid(uid uint, count *int64) error {
	if res := db.Model(&entity.Video{}).Where("author_id =  ?", uid).Count(count); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("find video error")
		return res.Error
	}
	return nil
}
