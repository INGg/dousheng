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
	// --feed--
	feedRouter := apiRouter.Group("/feed")
	feedRouter.GET("/", service.Feed) // /feed/

	// --user--
	userRouter := apiRouter.Group("/user")
	userRouter.GET("/", service.UserInfo)           // /user/
	userRouter.POST("/register/", service.Register) // /user/register/
	userRouter.POST("/login/", service.Login)       // /user/login/

	// --publish--
	publishRouter := apiRouter.Group("/publish")
	publishRouter.Use(middleware.JWTAuth())          // 声明中间件
	publishRouter.POST("/action/", service.Publish)  // /publish/action/
	publishRouter.GET("/list/", service.PublishList) // /publish/list/

	// extra apis - I

	// --favorite--
	favoriteRouter := apiRouter.Group("/favorite")
	favoriteRouter.Use(middleware.JWTAuth())                // 声明中间件
	favoriteRouter.POST("/action/", service.FavoriteAction) // /favorite/action/
	favoriteRouter.GET("/list/", service.FavoriteList)      // /favorite/list/

	// --comment--
	commentRouter := apiRouter.Group("/comment")
	commentRouter.Use(middleware.JWTAuth())               // 声明中间件
	commentRouter.POST("/action/", service.CommentAction) // /comment/action/
	commentRouter.GET("/list/", service.CommentList)      // /comment/list/

	// extra apis - II
	relationRouter := apiRouter.Group("/relation")
	relationRouter.Use(middleware.JWTAuth())                    // 声明中间件
	relationRouter.POST("/action/", service.RelationAction)     // /relation/action/
	relationRouter.GET("/follow/list/", service.FollowList)     // /relation/follow/list/
	relationRouter.GET("/follower/list/", service.FollowerList) // /relation/follower/list/
}
