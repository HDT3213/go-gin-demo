package recovery

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/lib/logger"
    "runtime/debug"
    "fmt"
    "github.com/go-gin-demo/lib/response"
)

func Recovery() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        defer func(ctx *gin.Context) {
            if err := recover(); err != nil {
                logger.Warn(fmt.Sprintf("error occurs: %v\n%s", err, string(debug.Stack())))
                response.InternalServerError(ctx, "internal server error")
            }
        }(ctx)
        ctx.Next()
    }
}

