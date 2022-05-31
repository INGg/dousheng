package main

import (
	"demo1/middleware"
	"demo1/service"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	//r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	feedRouter := apiRouter.Group("/feed")
	feedRouter.GET("/", service.Feed) // /feed/

	userRouter := apiRouter.Group("/user")
	userRouter.GET("/", service.UserInfo)           // /user/
	userRouter.POST("/register/", service.Register) // /user/register/
	userRouter.POST("/login/", service.Login)       // /user/login/

	publishRouter := apiRouter.Group("/publish")

	publishRouter.Use(middleware.JWTAuth())

	publishRouter.POST("/action/", service.Publish)  // /publish/action/
	publishRouter.GET("/list/", service.PublishList) // /publish/list/

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
