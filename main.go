package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/fvbock/endless"
    "github.com/go-gin-demo/model"
    "github.com/go-gin-demo/router"
)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    // init db
    model.Setup()
    defer model.Close()

    app := gin.Default()
    router.Setup(app)

    // run
    endless.ListenAndServe(":" + port, app)
}
