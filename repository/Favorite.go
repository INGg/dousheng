package repository

import (
	"demo1/model/entity"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"log"
	"sync"

	"gorm.io/gorm"
)

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
	res := db.Create(&entity.Favorite{
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
	res := db.Where(&entity.Favorite{UserID: uid, VideoID: vid}).Delete(&entity.Favorite{})
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("delete favorite error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) FindFavoriteVideoByUid(uid uint, videoList *[]entity.Video) error {
	res := db.Model(entity.Favorite{}).Select("video_id").Where("user_id = ?", uid).Find(&videoList)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		zap.L().Error("vid can't find error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) AddFavoriteCount(vid uint) error {
	res := db.Model(&entity.Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count + ?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Add favorite count error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) ReduceFavoriteCount(vid uint) error {
	res := db.Model(&entity.Video{}).Where("id = ?", vid).Update("favorite_count", gorm.Expr("favorite_count - ?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Reduce favorite count error")
		return res.Error
	}
	return nil
}

func (f *FavoriteDAO) CheckIsFavorite(uid uint, vid uint) bool {
	res := db.Where(&entity.Favorite{UserID: uid, VideoID: vid}).First(&uid)
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
