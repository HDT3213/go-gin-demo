package controller

import (
    "errors"
    "fmt"
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/app/mq/router"
    "github.com/HDT3213/go-gin-demo/lib/mq"
    "github.com/HDT3213/go-gin-demo/lib/response"
    "github.com/gin-gonic/gin"
    "time"
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