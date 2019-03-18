package model

import (
    "github.com/go-redis/redis"
    "fmt"
    "encoding/json"
    "time"
)

var Redis *redis.Client

func setupRedis() {
    var err error
    Redis = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        PoolSize:     16,
        PoolTimeout:  10 * time.Second,
    })

    pong, err := Redis.Ping().Result()
    fmt.Println(pong, err)
}

func closeCache() {
    if Redis != nil {
        Redis.Close()
    }
}


func Marshal(v interface{}) ([]byte, error) {
    return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
    return json.Unmarshal(data, v)
}