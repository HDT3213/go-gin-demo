package controller

import (
    "github.com/HDT3213/go-gin-demo/lib/response"
    "github.com/gin-gonic/gin"
)

func Index(ctx *gin.Context)  {
    response.Success(ctx)
}

func NotFound(ctx *gin.Context) {
    response.NotFound(ctx, "no router found")
}