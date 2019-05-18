package router

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/middleware/auth"
    "github.com/go-gin-demo/controller"
    "github.com/go-gin-demo/middleware/recovery"
)

func Setup() (*gin.Engine) {
    app := gin.New()
    app.Use(gin.Logger())
    app.Use(recovery.Recovery())

    app.GET("/", controller.Index)
    app.GET("/panic", controller.Panic)
    app.NoRoute(controller.NotFound)
    app.NoMethod(controller.NotFound)

    app.GET("/mq", controller.MQEcho)

    root := app.Group("")
    root.Use(auth.JWT())

    root.POST("/register", controller.Register)
    root.POST("/login", controller.Login)
    root.GET("/user/:id", controller.GetUser)
    root.GET("/self", controller.Self)

    root.GET("/post/:id", controller.GetPost)
    root.POST("/post", controller.CreatePost)
    root.DELETE("/post/:id", controller.DeletePost)
    root.GET("/timeline/self", controller.GetSelfTimeline)
    root.GET("/timeline/following", controller.GetFollowingTimeline)
    root.GET("/timeline/user/:uid", controller.GetUserTimeline)

    root.POST("/user/:id/follow", controller.Follow)
    root.POST("/user/:id/unfollow", controller.UnFollow)
    root.GET("/following/user/:id", controller.GetUserFollowings)
    root.GET("/follower/user/:id", controller.GetUserFollowers)
    root.GET("/following/self", controller.GetSelfFollowings)
    root.GET("/follower/self", controller.GetSelfFollowers)
    return app
}
