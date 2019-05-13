package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/lib/response"
    MQRouter "github.com/go-gin-demo/mq/router"
    "fmt"
    "time"
    "github.com/go-gin-demo/context/context"
    "github.com/go-gin-demo/lib/mq"
)

func Index(ctx *gin.Context)  {
    response.Success(ctx)
}

func NotFound(ctx *gin.Context) {
    response.NotFound(ctx, "no router found")
}

func InternalError(ctx *gin.Context) {
    response.InternalServerError(ctx, "internal server error")
}

func MQEcho(ctx *gin.Context) {
    msg := fmt.Sprintf("[%s] %s", time.Now().Format(time.RFC1123), ctx.Query("msg"))
    err := context.MQ.Publish(&mq.Msg{
        Code: MQRouter.Echo,
        Payload: []byte(msg),
    })
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
}