package repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
)

type Video struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	PublishTime   int64  `gorm:"index" json:"publish_time"`
	Author        User   `gorm:"foreignKey:AuthorID" json:"author"`
	AuthorID      uint   `json:"author_id"`
	PlayUrl       string `gorm:"type:varchar(128)" json:"play_url"`  // 播放地址
	CoverUrl      string `gorm:"type:varchar(128)" json:"cover_url"` // 封面地址
	FavoriteCount int64  `gorm:"not_null; default:0" json:"favorite_count"`
	CommentCount  int64  `gorm:"not_null; default:0" json:"comment_count"`
	Title         string ` json:"title"`
}

var VideoCount int64

func GetVideoList(video_list *[]Video, lim int, ReqTime int64) error {

	res := db.Limit(lim).Order("publish_time desc").Where("publish_time <= ?", ReqTime).Find(video_list)

	//for id, v := range *video_list {
	//	fmt.Println(id, v)
	//}

	if res.Error != nil {
		return res.Error
	}

	return nil
}

func FindVideoByPathAndUid(path string, uid int64, video *Video) error {
	if res := db.Model(Video{}).Where("play_url = ? AND author_id = ?", path, uid).First(video); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// InsertVideo 将视频信息写入数据库，返回错误信息和错误
func InsertVideo(token string, filepath string, title string) error {

	// 判断是否已经存在了这个视频

	// 生成url

	// 视频url
	VideoCount++
	play_url := filepath
	// 封面url
	cover_url := ""
	// 构造video

	// 根据token获取视频上传者
	var author User
	err := FindUserByToken(token, &author)
	if err != nil {
		log.Fatal("user not exists")
		return err
	}

	fmt.Println(token, filepath, title)

	// 存入数据库
	res := db.Create(&Video{
		PublishTime:   time.Now().Unix(),
		Author:        author,
		AuthorID:      author.ID,
		PlayUrl:       play_url,
		CoverUrl:      cover_url,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
	})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Fatal("create video error")
		return res.Error
	}

	fmt.Println("insert video info ok")
	return nil
}

// CheckIsFavorite 判断uid这个人是不是给这个视频点赞了
func CheckIsFavorite(uid uint, videoId uint) bool {
	return false
}

func FindAllVideoByUid(uid uint, VideoList *[]Video) error {

	//var list []Video
	res := db.Model(&Video{}).Where("author_id = ?", uid).Find(VideoList)
	//fmt.Println(len(*VideoList))
	//for i, video := range *VideoList {
	//	fmt.Println(i, video)
	//}
	if res.Error != nil {
		return res.Error
	}
	return nil
}
