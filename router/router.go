package router

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/middleware"
    "github.com/go-gin-demo/controller"
)

func Setup(app *gin.Engine) {
    app.GET("/", controller.Index)

    app.POST("/register", controller.Register)
    app.POST("/login", controller.Login)
    app.GET("/user/:id", controller.GetUser)

    loginRequired := app.Group("")
    loginRequired.Use(middleware.JWT())
    {
        //loginRequired.GET("/users", controller.AllUser)
        loginRequired.GET("/self", controller.Self)
    }
}
