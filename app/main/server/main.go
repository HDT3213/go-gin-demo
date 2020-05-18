package main

import (
    "fmt"
    "github.com/HDT3213/go-gin-demo/app/context/config"
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/app/model"
    MQRouter "github.com/HDT3213/go-gin-demo/app/mq/router"
    "github.com/HDT3213/go-gin-demo/app/router"
    Redis "github.com/HDT3213/go-gin-demo/lib/cache/redis"
    DB "github.com/HDT3213/go-gin-demo/lib/db"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    MQ "github.com/HDT3213/go-gin-demo/lib/mq"
    "github.com/fvbock/endless"
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


    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    app := router.Setup()

    // run
    logger.Info("start server")
    endless.ListenAndServe(":" + port, app)
}
