package service

import (
	"demo1/middleware"
	"demo1/model"
	"demo1/model/entity"
	"demo1/repository"
	"errors"
	"fmt"
)

func Register(req *model.UserRegisterRequest) (*model.UserResisterResponse, error) {
	// 判断这个用户名的是否存在

	// 2.向gorm发起请求判断用户是否存在
	var user entity.User

	// 创建单例
	userDAO := repository.NewUserDAO()

	if err := userDAO.FindUserIDByName(req.UserName, &user); err != nil { // 返回错误说明没找到
		// 创建用户
		uid, token, err := userDAO.CreateUser(req.UserName, req.Password)
		if err != nil { // 返回错误说明没有创建成功
			return &model.UserResisterResponse{
				Response: model.Response{
					StatusCode: 1,
					StatusMsg:  "create user error",
				},
				UserID: 0,
				Token:  "",
			}, err
		}

		// 成功创建返回正确的信息
		return &model.UserResisterResponse{
			Response: model.Response{
				StatusCode: 0,
				StatusMsg:  "register user successful",
			},
			UserID: uid,
			Token:  token,
		}, nil
	} else {
		return &model.UserResisterResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user exited",
			},
			UserID: user.ID,
			Token:  "",
		}, nil
	}
}

// ---Login---

func Login(req *model.UserLoginRequest) (*model.UserLoginResponse, error) {

	// 2.向gorm发起请求判断用户是否存在
	var user entity.User
	// 创建单例
	userDAO := repository.NewUserDAO()

	if err := userDAO.FindUserIDByName(req.UserName, &user); err != nil {
		return &model.UserLoginResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "user not existed",
			},
			UserID: 0,
			Token:  "",
		}, err
	}

	// 校验pwd，返回是否登陆成功
	uid, ok := userDAO.CheckUserPwd(req.UserName, req.Password)
	if ok == false {
		return &model.UserLoginResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%v pwd error\n", req.UserName),
			},
			UserID: 0,
			Token:  "",
		}, errors.New("pwd error")
	}

	token, err := middleware.GenToken(user.Name, user.ID)
	if err != nil {
		return &model.UserLoginResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "%v gen token error",
			},
			UserID: 0,
			Token:  "",
		}, err
	}

	return &model.UserLoginResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("%v login succussfully\n", req.UserName),
		},
		UserID: uid,
		Token:  token,
	}, nil
}

// ---UserInfo----

func UserInfo(req *model.UserInfoRequest) (*model.UserInfoResponse, error) {
	// 创建单例
	userDAO := repository.NewUserDAO()
	relationDAO := repository.NewRelationDAO()

	var user entity.User
	if err := userDAO.FindUserById(req.UserID, (&user)); err != nil {
		return &model.UserInfoResponse{
			Response: model.Response{
				StatusCode: 1,
				StatusMsg:  "get user info error (user not exists)",
			},
			User: entity.User{},
		}, err
	}

	user.IsFollow = relationDAO.QueryAFollowB(req.UserID, user.ID)

	return &model.UserInfoResponse{
		Response: model.Response{
			StatusCode: 0,
			StatusMsg:  "get user info successfully",
		},
		User: user,
	}, nil
}
