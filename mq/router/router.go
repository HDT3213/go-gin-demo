package router

import (
    "github.com/go-gin-demo/mq/core"
    "github.com/go-gin-demo/mq/consumer"
)

const (
    Echo = iota
)

func GetConsumerMap() map[uint32]func(*core.Msg) {
    consumerMap := make(map[uint32]func(*core.Msg))

    consumerMap[Echo] = consumer.Echo

    return consumerMap
}