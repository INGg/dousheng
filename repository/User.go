package repository

import (
	"demo1/middleware"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"sync"
)

var UserCount int64

type User struct {
	ID            uint   `gorm:"primaryKey; not null" json:"id"`
	Name          string `gorm:"type:varchar(64); not null;" json:"name"`
	Token         string `gorm:"-" json:"token,omitempty"`
	Password      string `gorm:"char(24) ; not null;" json:"password,omitempty"`
	FollowCount   int64  `gorm:"not null; default:0" json:"follow_count"`   // 关注的人的数量
	FollowerCount int64  `gorm:"not null; default:0" json:"follower_count"` // 粉丝总数
	IsFollow      bool   `gorm:"-" json:"is_follow"`
}

func TableName() string {
	return "users"
}

type UserDAO struct {
}

var (
	userDAO  *UserDAO
	userOnce sync.Once
)

func NewUserDAO() *UserDAO {
	userOnce.Do(func() {
		userDAO = &UserDAO{}
	})
	return userDAO
}

// MakeToken 构造token
func MakeToken(username string, uid uint) (string, error) {
	return middleware.GenToken(username, uid)
}

// FindUserIDByName 通过名字判断用户是否存在
func (u *UserDAO) FindUserIDByName(username string, user *User) error {
	if res := db.Model(User{}).Where("name = ?", username).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// FindUserById 通过id找到user
func (u *UserDAO) FindUserById(id uint, user *User) error {
	if res := db.Model(User{}).Select("name", "id", "follow_count", "follower_count").Where("id = ?", id).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find user error")
		return res.Error
	}
	return nil
}
func (u *UserDAO) FindMUserByIdList(idList []uint, userList *[]User) error {
	return nil
}

// CreateUser 向数据库写入User
func (u *UserDAO) CreateUser(username string, pwd string) (uid uint, token string, err error) {
	user := User{
		Name:          username,
		Password:      base64.StdEncoding.EncodeToString([]byte(pwd)),
		FollowCount:   0,
		FollowerCount: 0,
	}

	res := db.Create(&user)

	//UserCount++
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return 0, "", res.Error
	}

	// 创建token
	token, err = MakeToken(username, user.ID)
	if err != nil {
		return 0, "", err
	}

	return user.ID, token, nil
}

// CheckUserPwd 登陆检查密码，即检查password是否一致，一致返回userid和ok，ok为true为密码正确
func (u *UserDAO) CheckUserPwd(username string, pwd string) (uid uint, ok bool) {
	// 根据名字获取token
	var user User
	db.Where("name = ?", username).First(&user)
	// 构造password进行比对

	if user.Password != base64.StdEncoding.EncodeToString([]byte(pwd)) {
		return 0, false
	} else {
		return user.ID, true
	}
}

// JudgeAFollowB 判断a是否关注了b
func (u *UserDAO) JudgeAFollowB(uida int64, uidb int64) bool {
	//res := db.Where()
	return true
}

// AFollowB 让a关注/取关b，关注是1
func (u *UserDAO) AFollowB(ctx *gin.Context, op int32) error {
	return nil
}

//func (u *UserDAO) GetUserList(uid *[]uint, resUser *[]User) error {
//var res []User
//db.Find(&User{})
//}
