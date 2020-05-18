package main

import (
    "fmt"
    Router "github.com/HDT3213/go-gin-demo/app/canal"
    "github.com/HDT3213/go-gin-demo/app/context/config"
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/app/model"
    MQRouter "github.com/HDT3213/go-gin-demo/app/mq/router"
    Redis "github.com/HDT3213/go-gin-demo/lib/cache/redis"
    Canal "github.com/HDT3213/go-gin-demo/lib/canal"
    DB "github.com/HDT3213/go-gin-demo/lib/db"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    MQ "github.com/HDT3213/go-gin-demo/lib/mq"
    "os"
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
        logger.Error(fmt.Sprintf("mq start failed, %v", err))
        if settings.Runtime.AllowMQFallback {
            mq, _ := MQ.SetupGoRoutineMQ(MQRouter.GetConsumerMap())
            defer mq.Close()
            context.MQ = mq
        } else {
            logger.Fatal("abort start")
        }
    } else {
        defer mq.Close()
        context.MQ = mq
    }

    canal, err := Canal.Setup(&settings.Canal)
    if err != nil {
        logger.Fatal(fmt.Sprintf("mq start failed, %v", err))
    }

    logger.Info("start consume")
    canal.Listen(Router.GetConsumerMap())
}
