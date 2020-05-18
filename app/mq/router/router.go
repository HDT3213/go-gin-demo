package router

import (
    "github.com/HDT3213/go-gin-demo/app/mq/consumer"
    "github.com/HDT3213/go-gin-demo/lib/mq"
)

const (
    Echo = iota
)

func GetConsumerMap() map[uint32]mq.ConsumerFunc {
    consumerMap := make(map[uint32]mq.ConsumerFunc)

    consumerMap[Echo] = consumer.Echo

    return consumerMap
}