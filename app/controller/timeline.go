package controller

import (
    "github.com/HDT3213/go-gin-demo/app/middleware/auth"
    PostService "github.com/HDT3213/go-gin-demo/app/service/post"
    BizError "github.com/HDT3213/go-gin-demo/lib/errors"
    "github.com/HDT3213/go-gin-demo/lib/request"
    "github.com/HDT3213/go-gin-demo/lib/response"
    "github.com/gin-gonic/gin"
    "strconv"
)

func GetUserTimeline(ctx *gin.Context) {
    uid, err := strconv.ParseUint(ctx.Param("uid"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("uid"))
        return
    }
    start, length, err := request.GetPage(ctx,0 , 10)
    if err != nil {
        response.Error(ctx, err)
        return
    }

    currentUid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        if BizError.IsForbidden(err) {
            currentUid = 0
        } else {
            response.Error(ctx, err)
            return
        }
    }

    posts, total, err := PostService.GetUserTimeline(currentUid, uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, posts, total)
}

func GetSelfTimeline(ctx *gin.Context) {
    uid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    start, length, err := request.GetPage(ctx, 0, 10)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    posts, total, err := PostService.GetUserTimeline(uid, uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, posts, total)
}

func GetFollowingTimeline(ctx *gin.Context) {
    uid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    start, length, err := request.GetPage(ctx, 0, 10)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    posts, total, err := PostService.GetFollowingTimeline(uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, posts, total)
}