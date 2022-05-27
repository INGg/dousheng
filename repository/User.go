package repository

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var UserCount int64

type User struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	Name          string `gorm:"type:varchar(64); not null;" json:"name"`
	Token         string `gorm:"type:varchar(256); not null;" json:"token"`
	FollowCount   int64  `gorm:"not null; default:0" json:"follow_count"`   // 关注的人的数量
	FollowerCount int64  `gorm:"not null; default:0" json:"follower_count"` // 粉丝总数
	//IsFollow    bool   `json:"is_follow"`
}

func TableName() string {
	return "user"
}

// MakeToken 构造token
func MakeToken(username string, password string) string {
	return username + password
}

// FindUserByName 通过名字判断用户是否存在
func FindUserByName(username string, user *User) error {
	if res := db.Model(User{}).Where("name = ?", username).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// FindUserByToken 通过token判断用户是否存在
func FindUserByToken(token string, user *User) error {
	if res := db.Model(User{}).Where("token = ?", token).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// FindUserById 通过id找到user
func FindUserById(id uint, user *User) error {
	if res := db.Model(User{}).Where("id = ?", id).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find user error")
		return res.Error
	}
	return nil
}

// CreateUser 向数据库写入User
func CreateUser(username string, pwd string) (uid uint, token string, err error) {
	//token 转为md5
	token = MakeToken(username, pwd)

	user := User{
		Name:          username,
		Token:         token,
		FollowCount:   0,
		FollowerCount: 0,
	}

	res := db.Create(&user)

	//UserCount++
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, "", res.Error
	}

	return user.ID, user.Token, nil
}

// CheckUserToken 登陆检查密码，即检查token是否一致，一致返回userid和token，否则返回错误
func CheckUserToken(username string, password string) (uid uint, token string, ok bool) {
	// 根据名字获取token
	var user User
	db.Where("name = ?", username).First(&user)
	// 构造token进行比对

	token = MakeToken(username, password)
	//fmt.Println("token1 : " + token)
	//fmt.Println("token2 : " + user.Token)

	if user.Token != token {
		return 0, "", false
	} else {
		return user.ID, token, true
	}
}

// JudgeAFollowB 判断a是否关注了b
func JudgeAFollowB(uida int64, uidb int64) bool {
	//res := db.Where()
	return true
}

// AFollowB 让a关注/取关b，关注是1
func AFollowB(ctx *gin.Context, op int32) error {
	return nil
}
