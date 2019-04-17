package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/utils/response"
    MQ "github.com/go-gin-demo/mq/core"
    MQRouter "github.com/go-gin-demo/mq/router"
    "fmt"
    "time"
)

func Index(ctx *gin.Context)  {
    response.Success(ctx)
}

func MQEcho(ctx *gin.Context) {
    msg := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC1123), ctx.Query("msg"))
    err := MQ.Publish(&MQ.Msg{
        Code: MQRouter.Echo,
        Payload: []byte(msg),
    })
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
}