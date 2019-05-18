package controller

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "time"
    "github.com/go-gin-demo/context/context"
    "github.com/go-gin-demo/lib/mq"
    "github.com/go-gin-demo/lib/response"
    "errors"
    "github.com/go-gin-demo/mq/router"
)

func Panic(ctx *gin.Context) {
    panic(errors.New("test"))
}

func MQEcho(ctx *gin.Context) {
    msg := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC1123), ctx.Query("msg"))
    err := context.MQ.Publish(&mq.Msg{
        Code: router.Echo,
        Payload: []byte(msg),
    })
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
}