package router

import (
    "github.com/go-gin-demo/mq/consumer"
    "github.com/go-gin-demo/lib/mq"
)

const (
    Echo = iota
)

func GetConsumerMap() map[uint32]func(*mq.Msg) {
    consumerMap := make(map[uint32]func(*mq.Msg))

    consumerMap[Echo] = consumer.Echo

    return consumerMap
}