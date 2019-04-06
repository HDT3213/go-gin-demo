package controller

import (
    "github.com/gin-gonic/gin"
    "strconv"
    "github.com/go-gin-demo/utils/response"
    PostService "github.com/go-gin-demo/service/post"
    "github.com/go-gin-demo/utils"
    "github.com/go-gin-demo/middleware"
)

func GetUserTimeline(ctx *gin.Context) {
    uid, err := strconv.ParseUint(ctx.Param("uid"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("uid"))
        return
    }
    start, length, err := utils.GetPage(ctx,0 , 10)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    posts, total, err := PostService.GetUserTimeline(uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, posts, total)
}

func GetSelfTimeline(ctx *gin.Context) {
    uid, err := middleware.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    start, length, err := utils.GetPage(ctx, 0, 10)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    posts, total, err := PostService.GetUserTimeline(uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, posts, total)
}