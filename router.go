package main

import (
	"demo1/service"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	//r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	apiRouter.GET("/feed/", service.Feed)
	apiRouter.GET("/user/", service.UserInfo)
	apiRouter.POST("/user/register/", service.Register)
	apiRouter.POST("/user/login/", service.Login)
	apiRouter.POST("/publish/action/", service.Publish)
	apiRouter.GET("/publish/list/", service.PublishList)

	// extra apis - I
	//apiRouter.POST("/favorite/action/", controller.FavoriteAction)
	//apiRouter.GET("/favorite/list/", controller.FavoriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)
	//
	//// extra apis - II
	//apiRouter.POST("/relation/action/", controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)
}
