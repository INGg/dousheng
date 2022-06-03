package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
	"time"
)

type Video struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	PublishTime   int64  `gorm:"index; timestamp" json:"publish_time"`
	Author        User   `gorm:"foreignKey:AuthorID" json:"author"`
	AuthorID      uint   //`json:"author_id"`
	PlayUrl       string `gorm:"type:varchar(128)" json:"play_url"`  // 播放地址
	CoverUrl      string `gorm:"type:varchar(128)" json:"cover_url"` // 封面地址
	FavoriteCount int64  `gorm:"not_null; default:0" json:"favorite_count"`
	CommentCount  int64  `gorm:"not_null; default:0" json:"comment_count"`
	Title         string ` json:"title"`
}

var VideoCount int64

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
func (v *VideoDAO) GetVideoList(videoList *[]Video, lim int, ReqTime int64) error {

	res := db.Limit(lim).Order("publish_time desc").Where("publish_time <= ?", ReqTime).Find(videoList)

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func (v *VideoDAO) FindVideoByPathAndUid(path string, uid int64, video *Video) error {
	if res := db.Model(Video{}).Where("play_url = ? AND author_id = ?", path, uid).First(video); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// InsertVideo 将视频信息写入数据库，返回错误信息和错误
func (v *VideoDAO) InsertVideo(uid uint, filepath string, title string) error {

	// TODO 判断是否已经存在了这个视频

	// 生成url
	// 视频url
	VideoCount++
	playUrl := "http://" + IP + HOST + filepath
	// 封面url
	coverUrl := "https://cdn.pixabay.com/photo/2016/03/27/18/10/bear-1283347_1280.jpg"
	// 构造video

	// 根据token获取视频上传者
	fmt.Println(uid, filepath, title)

	// 存入数据库
	res := db.Create(&Video{
		PublishTime:   time.Now().Unix(),
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

// CheckIsFavorite 判断uid这个人是不是给这个视频点赞了
func (v *VideoDAO) CheckIsFavorite(uid uint, videoId uint) bool {
	// TODO 这个是不是应该写在Favorite里面
	return false
}

func (v *VideoDAO) VideoCount() int64 {
	var count int64
	db.Model(&Video{}).Count(&count)
	return count
}

// FindAllVideoByUid 通过uid找到这个人发布的所有视频
func (v *VideoDAO) FindAllVideoByUid(uid uint, VideoList *[]Video) error {
	res := db.Model(&Video{}).Where("author_id = ?", uid).Find(VideoList)
	if res.Error != nil {
		return res.Error
	}
	for _, video := range *VideoList {
		fmt.Printf("%+v\n", video)
	}
	return nil
}

// FindVideoById 通过id找到Video
func (v *VideoDAO) FindVideoById(uid uint, video *Video) error {
	if res := db.Model(User{}).Where("author_id = ?", uid).First(video); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find user error")
		return res.Error
	}
	return nil
}
