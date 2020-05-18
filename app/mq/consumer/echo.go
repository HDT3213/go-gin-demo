package consumer

import (
    "github.com/HDT3213/go-gin-demo/lib/logger"
    "github.com/HDT3213/go-gin-demo/lib/mq"
)

func Echo(msg *mq.Msg) {
    logger.Info("receive: " + string(msg.Payload))
}