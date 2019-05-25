package main

import (
    "github.com/go-gin-demo/lib/logger"
    "os"
    "github.com/go-gin-demo/context/config"
    MQ "github.com/go-gin-demo/lib/mq"
    "fmt"
    "github.com/go-gin-demo/context/context"
    DB "github.com/go-gin-demo/lib/db"
    "github.com/go-gin-demo/model"
    Redis "github.com/go-gin-demo/lib/cache/redis"
    Canal "github.com/go-gin-demo/lib/canal"
    Router "github.com/go-gin-demo/canal"
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

    canal, err := Canal.Setup(&settings.Canal)
    if err != nil {
        logger.Fatal(fmt.Sprintf("mq start failed, %v", err))
    }

    logger.Info("start consume")
    canal.Listen(Router.GetConsumerMap())
}
