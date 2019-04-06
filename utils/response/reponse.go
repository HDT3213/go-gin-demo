package response

import "github.com/gin-gonic/gin"

func Success(context *gin.Context) {
    context.JSON(200, gin.H{
        "success": true,
    })
}

func Entity(ctx *gin.Context, entity interface{}) {
    ctx.JSON(200, gin.H{
        "success": true,
        "entity": entity,
    })
}

func Entities(ctx *gin.Context, entities interface{}, total int32) {
    ctx.JSON(200, gin.H{
        "success": true,
        "entities": entities,
        "total": total,
    })
}

func BadRequest(ctx *gin.Context, msg string) {
    ctx.JSON(400, gin.H{
        "success": false,
        "msg": msg,
    })
}

func NotFound(ctx *gin.Context, msg string) {
    ctx.JSON(404, gin.H{
        "success": false,
        "msg": msg,
    })
}

func Forbidden(ctx *gin.Context, msg string) {
    ctx.JSON(400, gin.H{
        "success": false,
        "msg": msg,
    })
}

func InternalServerError(ctx *gin.Context, msg string) {
    ctx.JSON(500, gin.H{
        "success": false,
        "msg": msg,
    })
}