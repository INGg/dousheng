package service

import (
	"demo1/middleware"
	"demo1/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Register(c *gin.Context) {
	// 判断这个用户名的是否存在

	// 1.拿去请求信息
	var req UserRegisterRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserResisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "login should bind error",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	//fmt.Println("Register")
	//fmt.Println(req.UserName, req.Password)

	// 2.向gorm发起请求判断用户是否存在
	var user repository.User

	// 创建单例
	userDAO := repository.NewUserDAO()

	if err := userDAO.FindUserByName(req.UserName, &user); err != nil { // 返回错误说明没找到
		// 创建用户
		uid, token, err := userDAO.CreateUser(req.UserName, req.Password)
		if err != nil { // 返回错误说明没有创建成功
			c.JSON(http.StatusOK, UserResisterResponse{
				Response: Response{
					StatusCode: 1,
					StatusMsg:  "create user error",
				},
				UserId: 0,
				Token:  "",
			})
			return
		}

		// 成功创建返回正确的信息
		c.JSON(http.StatusOK, UserResisterResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "register user successful",
			},
			UserId: uid,
			Token:  token,
		})
	} else {
		c.JSON(http.StatusOK, UserResisterResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "user exited",
			},
			UserId: user.ID,
			Token:  "",
		})
	}
}

// ---Login---

func Login(c *gin.Context) {
	// 1.获取请求参数
	var req UserLoginRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "login should bind error",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	// 2.向gorm发起请求判断用户是否存在
	var user repository.User
	// 创建单例
	userDAO := repository.NewUserDAO()

	if err := userDAO.FindUserByName(req.UserName, &user); err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "user not existed",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	// 校验token，返回是否登陆成功
	uid, ok := userDAO.CheckUserPwd(req.UserName, req.Password)
	if ok == false {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  fmt.Sprintf("%v pwd error\n", req.UserName),
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	token, err := middleware.GenToken(user.Name)
	if err != nil {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "%v gen token error",
			},
			UserId: 0,
			Token:  "",
		})
		return
	}

	c.JSON(http.StatusOK, UserLoginResponse{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  fmt.Sprintf("%v login succussfully\n", req.UserName),
		},
		UserId: uid,
		Token:  token,
	})
}

// ---UserInfo----

func UserInfo(c *gin.Context) {
	var req UserInfoRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "UserInfo should bind error",
			},
			User: User{},
		})
		return
	}
	req.UserName = c.GetString("username")

	// 创建单例
	userDAO := repository.NewUserDAO()

	var user User
	if err := userDAO.FindUserById(req.UserId, &user.User); err == nil {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{
				StatusCode: 0,
				StatusMsg:  "get user info successfully",
			},
			User: user,
		})
	} else {
		c.JSON(http.StatusOK, UserInfoResponse{
			Response: Response{
				StatusCode: 1,
				StatusMsg:  "get user info error (user not exists)",
			},
			User: User{},
		})
	}
}
