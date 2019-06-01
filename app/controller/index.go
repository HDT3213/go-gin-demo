package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/lib/response"
)

func Index(ctx *gin.Context)  {
    response.Success(ctx)
}

func NotFound(ctx *gin.Context) {
    response.NotFound(ctx, "no router found")
}