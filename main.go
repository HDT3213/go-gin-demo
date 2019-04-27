package main

import (
    "os"
    "github.com/gin-gonic/gin"
    "github.com/fvbock/endless"
    "github.com/go-gin-demo/model"
    "github.com/go-gin-demo/router"
    MQCore "github.com/go-gin-demo/mq/core"
    MQRouter "github.com/go-gin-demo/mq/router"
    "github.com/go-gin-demo/utils/logger"
    "github.com/go-gin-demo/config"
)

func startServer() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    // init db
    model.Setup(&settings.DB, &settings.Redis)
    defer model.Close()

    MQCore.SetupRabbitMQ(&settings.Rabbit)
    defer MQCore.CloseRabbitMQ()

    app := gin.Default()
    router.Setup(app)

    logger.Info("start server")

    // run
    endless.ListenAndServe(":" + port, app)
}

func startMqConsumer() {
    // init db
    model.Setup(&settings.DB, &settings.Redis)
    defer model.Close()

    MQCore.SetupRabbitMQ(&settings.Rabbit)
    defer MQCore.CloseRabbitMQ()

    logger.Info("start consume")
    MQCore.Consume(MQRouter.GetConsumerMap())
}

var settings *config.Settings

func main() {
    configPath := os.Getenv("CONFIG")
    if configPath == "" {
        configPath = "./config.yml"
    }
    settings = config.Setup(configPath)

    logger.Setup(&settings.Log)

    role := os.Getenv("ROLE")
    if role == "server" || role == "" {
        startServer()
    } else if role == "consumer" {
        startMqConsumer()
    } else {
        logger.Fatal("illegal role: %s", role)
    }

}
