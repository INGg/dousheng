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

type User repository.User

type FeedRequest struct {
	LatestTime int64  `json:"latest_time"`
	Token      string `json:"token"`
}

type FeedResponse struct {
	Response
	VideoList *[]Video `json:"video_list"`
	NextTime  int64    `json:"next_time"`
}

type PublishActionRequest struct {
	Token    string                `json:"token" form:"token"`
	Data     *multipart.FileHeader `json:"data" form:"data"`
	Title    string                `json:"title" form:"title"`
	UserName string
}

type PublishActionResponse struct {
	Response
}

type PublishListRequest struct {
	Token    string `json:"token" form:"token"`
	UserId   uint   `json:"user_id" form:"user_id"`
	UserName string
}

type PublishListResponse struct {
	Response
	VideoList *[]Video `json:"video_list"`
}

type UserInfoRequest struct {
	Token    string `json:"token" form:"token"`
	UserId   uint   `json:"user_id" form:"user_id" `
	UserName string
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

// Comment 这里的comment和repository里的comment本质是没有太大的区别，这里只是为了请求好返回
type Comment struct {
	ID uint
	User
	Content    string `json:"content"`
	CreateDate string `json:"create_date"` // 评论发布日期，格式 mm-dd
}

type CommentActionRequest struct {
	UserId      uint   `json:"user_id"`
	Token       string `json:"token"`
	VideoId     uint   `json:"video_id"`
	ActionType  uint8  `json:"action_type"`
	CommentText string `json:"comment_text"` // 用户填写的评论内容，在action_type=1的时候使用
	CommentId   uint   `json:"comment_id"`   // 要删除的评论id，在action_type=2的时候使用
	UserName    string
}

type CommentActionResponse struct {
	Response
	Comment
}

type CommentListRequest struct {
	Token    string `json:"token"`
	VideoId  uint   `json:"video_id"`
	UserName string
}

type CommentListResponse struct {
	Response
	CommentList *[]Comment
}

type UserFavoriteRequest struct {
	UserId     uint   `json:"user_id"`
	Token      string `json:"token"`
	VideoId    uint   `json:"video_id"`
	ActionType int32  `json:"action_type"`
}

type UserFavoriteResponse struct {
	Response
}

type UserFavoriteListRequest struct {
	UserId uint   `json:"user_id"`
	Token  string `json:"token"`
}

type UserFavoriteListResponse struct {
	Response
	VideoList []repository.Video `json:"video_list"`
}
