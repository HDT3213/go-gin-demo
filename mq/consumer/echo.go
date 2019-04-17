package consumer

import (
    "github.com/go-gin-demo/mq/core"
    "log"
)

func Echo(msg *core.Msg) {
    log.Println("receive: " + string(msg.Payload))
}