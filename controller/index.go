package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/utils/response"
)

func Index(context *gin.Context)  {
    response.Success(context)
}