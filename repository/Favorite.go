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

func (f *FavoriteDAO) Favorite(uid uint, vid uint) error {
	//存入点赞记录
	res := db.Create(&Favorite{
		UserID:  uid,
		VideoID: vid,
	})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("create favorite error")
		return res.Error
	}
	fmt.Println("insert favorite info ok")
	return nil 
}

func (f *FavoriteDAO) UnFavorite(uid uint, vid uint) error {
	//删除点赞记录
	res := db.Where(&Favorite{UserID: uid, VideoID: vid}).Delete(&Favorite{})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("delete favorite error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) FindFavoriteVideoByUid(uid uint, videoList *[]Video) error {
	var vids []uint
	//将扫描到的vid存到vids切片里
	res := db.Model(Favorite{}).Select("video_id").Where("user_id = ?", uid).Find(&vids)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("vid can't find error")
		return res.Error
	}
	videoDAO := NewVideoDAO()
	for i := range vids {
		var video Video
		videoDAO.FindVideoById(vids[i], &video)
		*videoList = append(*videoList, video)
	}
	return nil
}

func (f *FavoriteDAO) AddFavoriteCount(vid uint) error {
	res := db.Model(&Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Add favorite count error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) ReduceFavoriteCount(vid uint) error {
	res := db.Model(&Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Reduce favorite count error")
		return res.Error
	}
	return nil
}
