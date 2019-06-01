package consumer

import (
    "github.com/go-gin-demo/lib/logger"
    "github.com/go-gin-demo/lib/mq"
)

func Echo(msg *mq.Msg) {
    logger.Info("receive: " + string(msg.Payload))
}