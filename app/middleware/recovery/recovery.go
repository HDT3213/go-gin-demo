package recovery

import (
    "fmt"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    "github.com/HDT3213/go-gin-demo/lib/response"
    "github.com/gin-gonic/gin"
    "runtime/debug"
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

