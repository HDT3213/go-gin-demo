package main

import (
    "github.com/go-gin-demo/context/config"
    "github.com/go-gin-demo/lib/logger"
    "fmt"
    "github.com/go-gin-demo/context/context"
    "github.com/go-gin-demo/model"
    "os"
    "github.com/fvbock/endless"
    MQ "github.com/go-gin-demo/lib/mq"
    DB "github.com/go-gin-demo/lib/db"
    Redis "github.com/go-gin-demo/lib/cache/redis"
    "github.com/go-gin-demo/router"
)

func main() {
    // load config
    configPath := os.Getenv("CONFIG")
    if configPath == "" {
        configPath = "./config.yml"
    }
    settings := config.Setup(configPath)

    // setup logger
    logger.Setup(&settings.Log)

    // set up mysql
    db, err := DB.Setup(&settings.DB)
    if err != nil {
        logger.Fatal(fmt.Sprintf("db start failed, %v", err))
    }
    defer db.Close()
    context.DB = db
    model.Migrate(db)

    // setup redis
    redis, err := Redis.Setup(&settings.Redis)
    if err != nil {
        logger.Fatal(fmt.Sprintf("redis start failed, %v", err))
    }
    defer redis.Close()
    context.Redis = redis

    // setup rabbit mq
    mq, err := MQ.SetupRabbitMQ(&settings.Rabbit)
    if err != nil {
        logger.Fatal(fmt.Sprintf("mq start failed, %v", err))
    }
    defer mq.Close()
    context.MQ = mq


    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    app := router.Setup()

    // run
    logger.Info("start server")
    endless.ListenAndServe(":" + port, app)
}
