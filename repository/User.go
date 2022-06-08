package repository

import (
	"demo1/middleware"
	"demo1/model/entity"
	"demo1/util"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
	"sync"
)

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
func (u *UserDAO) FindUserIDByName(username string, user *entity.User) error {
	if res := db.Model(entity.User{}).Where("name = ?", username).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return res.Error
	} else {
		return nil
	}
}

// FindUserById 通过id找到user
func (u *UserDAO) FindUserById(id uint, user *entity.User) error {
	if res := db.Model(entity.User{}).Select("name", "id", "follow_count", "follower_count").Where("id = ?", id).First(user); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find user error")
		return res.Error
	}
	return nil
}

// FindUsersByIdList 通过id数组找到user数组
func (u *UserDAO) FindUsersByIdList(idList []uint, userList *[]entity.User) error {
	if res := db.Model(entity.User{}).Select("name", "id", "follow_count", "follower_count").Where("id IN ?", idList).Find(userList); errors.Is(res.Error, gorm.ErrRecordNotFound) {
		fmt.Println("find user error")
		return res.Error
	}
	return nil
}

// CreateUser 向数据库写入User
func (u *UserDAO) CreateUser(username string, pwd string) (uid uint, token string, err error) {
	user := entity.User{
		Name:          username,
		Password:      util.MakeMD5(pwd),
		FollowCount:   0,
		FollowerCount: 0,
	}

	res := db.Create(&user)
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
	var user entity.User
	db.Where("name = ?", username).First(&user)
	// 构造password进行比对

	if user.Password != util.MakeMD5(pwd) {
		return 0, false
	} else {
		return user.ID, true
	}
}

// UpdateUserFollowerCount 增加某人的粉丝数量
func (u *UserDAO) UpdateUserFollowerCount(uid uint) error {
	res := db.Model(&entity.User{ID: uid}).UpdateColumn("follower_count", gorm.Expr("follower_count+?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Update FollowerCount error")
		return res.Error
	}
	return nil
}

// ReduceFollowerCount 减少某人的粉丝数量
func (u *UserDAO) ReduceFollowerCount(uid uint) error {
	res := db.Model(&entity.User{ID: uid}).UpdateColumn("follower_count", gorm.Expr("follower_count-?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Update FollowerCount error")
		return res.Error
	}
	return nil
}

// UpdateUserFollowCount 增加某人的关注数量
func (u *UserDAO) UpdateUserFollowCount(uid uint) error {
	res := db.Model(&entity.User{ID: uid}).UpdateColumn("follow_count", gorm.Expr("follow_count+?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Update FollowCount error")
		return res.Error
	}
	return nil
}

// ReduceFollowCount 减少某人的关注数量
func (u *UserDAO) ReduceFollowCount(uid uint) error {
	res := db.Model(&entity.User{ID: uid}).UpdateColumn("follow_count", gorm.Expr("follow_count-?", 1))
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		log.Print("Update FollowerCount error")
		return res.Error
	}
	return nil
}
