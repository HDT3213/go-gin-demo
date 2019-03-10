package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/fvbock/endless"
    "go-close/model"
    "go-close/router"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    // init db
    model.Setup()
    defer model.Close()

    model.Test()

    app := gin.Default()
    router.Setup(app)

    // run
    endless.ListenAndServe(":" + port, app)
}