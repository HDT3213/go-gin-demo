package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/jinzhu/gorm"
    "github.com/fvbock/endless"
    "./controller"
    "./model"
)

var DB *gorm.DB

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    // init db
    model.Setup()
    defer model.Close()

    app := gin.Default()

    // router
    app.GET("/", controller.Index)
    app.GET("/user/:id", controller.GetUser)
    app.GET("/users", controller.AllUser)
    app.POST("/register", controller.Register)

    // run
    endless.ListenAndServe(":" + port, app)
}
