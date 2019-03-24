package controller

import (
    "github.com/gin-gonic/gin"
    "strconv"
    "github.com/go-gin-demo/middleware"
    "github.com/go-gin-demo/utils/response"
    UserService "github.com/go-gin-demo/service/user"
)

func GetUser(ctx *gin.Context)  {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
        return
    }
    user, err := UserService.GetUser(uid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, user)
}

func Register(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    entity, err := UserService.Register(username, password)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, entity)
}

func Login(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    _, uid, err := UserService.Login(username, password)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    middleware.SetCurrentUid(ctx, uid)
    response.Success(ctx)
}

func Self(ctx *gin.Context) {
    uid, err := middleware.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    user, err := UserService.GetUser(uid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, user)
}