package main

import (
	"demo1/controller"
	"demo1/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	//r.Static("/static", "./public")

	apiRouter := r.Group("/douyin")

	// basic apis
	// --feed--
	feedRouter := apiRouter.Group("/feed")
	feedRouter.GET("/", controller.Feed) // /feed/

	// --user--
	userRouter := apiRouter.Group("/user")
	userRouter.GET("/", controller.UserInfo)           // /user/
	userRouter.POST("/register/", controller.Register) // /user/register/
	userRouter.POST("/login/", controller.Login)       // /user/login/

	// --publish--
	publishRouter := apiRouter.Group("/publish")
	publishRouter.Use(middleware.JWTAuth())             // 声明中间件
	publishRouter.POST("/action/", controller.Publish)  // /publish/action/
	publishRouter.GET("/list/", controller.PublishList) // /publish/list/

	// extra apis - I

	// --favorite--
	favoriteRouter := apiRouter.Group("/favorite")
	favoriteRouter.Use(middleware.JWTAuth())                   // 声明中间件
	favoriteRouter.POST("/action/", controller.FavoriteAction) // /favorite/action/
	favoriteRouter.GET("/list/", controller.FavoriteList)      // /favorite/list/

	// --comment--
	commentRouter := apiRouter.Group("/comment")
	commentRouter.Use(middleware.JWTAuth())                  // 声明中间件
	commentRouter.POST("/action/", controller.CommentAction) // /comment/action/
	commentRouter.GET("/list/", controller.CommentList)      // /comment/list/

	// extra apis - II
	relationRouter := apiRouter.Group("/relation")
	relationRouter.Use(middleware.JWTAuth())                       // 声明中间件
	relationRouter.POST("/action/", controller.RelationAction)     // /relation/action/
	relationRouter.GET("/follow/list/", controller.FollowList)     // /relation/follow/list/
	relationRouter.GET("/follower/list/", controller.FollowerList) // /relation/follower/list/
}
