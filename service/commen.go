package service

import (
	"demo1/repository"
	"mime/multipart"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	repository.Video
	IsFavorite bool
}

type User struct {
	repository.User
	IsFollow bool `json:"is_follow"`
}

type FeedRequest struct {
	LatestTime int64  `json:"latest_time,omitempty"`
	Token      string `json:"token,omitempty"`
}

type FeedResponse struct {
	StatusCode int32   `json:"status_code,omitempty"`
	StatusMsg  string  `json:"status_msg,omitempty"`
	VideoList  []Video `json:"video_list,omitempty"`
	NextTime   int64   `json:"next_time,omitempty"`
}

type PublishActionRequest struct {
	Token string                `json:"token" form:"token"`
	Data  *multipart.FileHeader `json:"data" form:"data"`
	Title string                `json:"title" form:"title"`
}

type PublishActionResponse struct {
	Response
}

type PublishListRequest struct {
	Token  string `json:"token" form:"token"`
	UserId uint   `json:"user_id" form:"user_id"`
}

type PublishListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

type UserInfoRequest struct {
	UserId uint   `json:"user_id" form:"user_id" `
	Token  string `json:"token" form:"token"`
}

type UserInfoResponse struct {
	Response
	User User `json:"user"`
}

type UserLoginRequest struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserLoginResponse struct {
	Response
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserRegisterRequest struct {
	UserName string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserResisterResponse struct {
	Response
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}
