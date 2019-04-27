package consumer

import (
    "github.com/go-gin-demo/mq/core"
    "github.com/go-gin-demo/utils/logger"
)

func Echo(msg *core.Msg) {
    logger.Info("receive: " + string(msg.Payload))
}