package controller

import (
    "github.com/gin-gonic/gin"
    PostService "github.com/go-gin-demo/service/post"
    "github.com/go-gin-demo/middleware/auth"
    "github.com/go-gin-demo/utils/response"
    "strconv"
    BizError "github.com/go-gin-demo/errors"
)

func CreatePost(ctx *gin.Context) {
    uid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    text := ctx.PostForm("text")
    if text == "" {
        response.Error(ctx, BizError.InvalidForm("text is required"))
        return
    }
    post, err := PostService.CreatePost(uid, text)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, post)
}

func GetPost(ctx *gin.Context) {
    pid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid id: " + ctx.Param("id"))
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

    post, err := PostService.GetPost(currentUid, pid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, post)
}


func DeletePost(ctx *gin.Context) {
    pid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid id: " + ctx.Param("id"))
        return
    }
    uid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    err = PostService.DeletePost(uid, pid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
}