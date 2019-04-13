package router

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/middleware"
    "github.com/go-gin-demo/controller"
)

func Setup(app *gin.Engine) {
    root := app.Group("")
    root.Use(middleware.JWT())

    app.GET("/", controller.Index)

    root.POST("/register", controller.Register)
    root.POST("/login", controller.Login)
    root.GET("/user/:id", controller.GetUser)
    root.GET("/self", controller.Self)

    root.GET("/post/:id", controller.GetPost)
    root.POST("/post", controller.CreatePost)
    root.DELETE("/post/:id", controller.DeletePost)
    root.GET("/timeline/self", controller.GetSelfTimeline)
    root.GET("/timeline/user/:uid", controller.GetUserTimeline)

    root.POST("/user/:id/follow", controller.Follow)
    root.POST("/user/:id/unfollow", controller.UnFollow)
    root.GET("/following/user/:id", controller.GetUserFollowings)
    root.GET("/follower/user/:id", controller.GetUserFollowers)
    root.GET("/following/self", controller.GetSelfFollowings)
    root.GET("/follower/self", controller.GetSelfFollowers)
}
