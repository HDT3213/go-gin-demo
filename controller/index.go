package controller

import (
    "github.com/gin-gonic/gin"
    "go-close/utils/response"
)

func Index(context *gin.Context)  {
    response.Success(context)
}