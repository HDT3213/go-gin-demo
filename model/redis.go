package model

import (
    "github.com/go-redis/redis"
    "fmt"
)

var Client *redis.Client

func setupRedis() {
    var err error
    Client = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    pong, err := Client.Ping().Result()
    fmt.Println(pong, err)
}

func closeCache() {
    if Client != nil {
        Client.Close()
    }
}



