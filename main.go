package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/fvbock/endless"
    "github.com/go-gin-demo/model"
    "github.com/go-gin-demo/router"
    MQCore "github.com/go-gin-demo/mq/core"
    MQRouter "github.com/go-gin-demo/mq/router"
    "log"
)

func startServer() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    // init db
    model.Setup()
    defer model.Close()

    MQCore.SetupRabbitMQ()
    defer MQCore.CloseRabbitMQ()

    app := gin.Default()
    router.Setup(app)

    // run
    endless.ListenAndServe(":" + port, app)
}

func startMqConsumer() {
    // init db
    model.Setup()
    defer model.Close()

    MQCore.SetupRabbitMQ()
    defer MQCore.CloseRabbitMQ()

    log.Println("start consume")
    MQCore.Consume(MQRouter.GetConsumerMap())
}

func main() {
    role := os.Getenv("ROLE")
    if role == "server" || role == "" {
        startServer()
    } else if role == "consumer" {
        startMqConsumer()
    } else {
        log.Fatalf("illegal role: %s", role)
    }

}
