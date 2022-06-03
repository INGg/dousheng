package repository

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"
)

type Favorite struct {
	ID      uint `gorm:"primaryKey; not null; auto_increment" json:"id"`
	UserID  uint `gorm:"not null" json:"user_id"`
	VideoID uint `gorm:"not null" json:"video_id"`
	Deleted gorm.DeletedAt
}

type FavoriteDAO struct {
}

var (
	favoriteDAO  *FavoriteDAO
	favoriteOnce sync.Once
)

func NewFavoriteDAO() *FavoriteDAO {
	favoriteOnce.Do(func() {
		favoriteDAO = &FavoriteDAO{}
	})
	return favoriteDAO
}

func (f *FavoriteDAO) ChangeFavorite(uid uint, vid uint, action_type int8) error {
	if action_type == 1 {
		//存入点赞记录
		res := db.Create(&Favorite{
			UserID:  uid,
			VideoID: vid,
		})
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Print("create video error")
			return res.Error
		}

		fmt.Println("insert favorite info ok")
		return nil
	} else {
		//删除点赞记录
		res := db.Where(&Favorite{UserID: uid, VideoID: vid}).Delete(&Favorite{})
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			log.Print("delete video error")
			return res.Error
		}
	}
	return nil
}

func (f *FavoriteDAO) FindFavoriteVideoByUid(uid uint, videoList *[]Video) error {
	var vids []uint
	var video Video
	//将扫描到的vid存到vids切片里
	res := db.Select("vid").Where("user_id = ?", uid).Find(&vids)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("vid can't find error")
		return res.Error
	}
	videoDAO := NewVideoDAO()
	for i := range vids {
		videoDAO.FindVideoById(vids[i], &video)
		*videoList = append(*videoList, video)
	}
	return nil
}
