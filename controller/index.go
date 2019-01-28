package controller

import "github.com/gin-gonic/gin"

func Index(context *gin.Context)  {
    Success(context)
}

func Success(context *gin.Context) {
    context.JSON(200, gin.H{
        "success": true,
    })
}

func Entity(context *gin.Context, entity interface{}) {
    context.JSON(200, gin.H{
        "success": true,
        "entity": entity,
    })
}

func BadRequest(context *gin.Context, msg string) {
    context.JSON(400, gin.H{
        "success": false,
        "msg": msg,
    })
}

func NotFound(context *gin.Context, msg string) {
    context.JSON(404, gin.H{
        "success": false,
        "msg": msg,
    })
}