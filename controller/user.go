package controller

import (
    "github.com/gin-gonic/gin"
    "go-close/model"
    "strconv"
    "fmt"
    "go-close/middleware"
    "go-close/utils/response"
    UserService "go-close/service/user"
)

func GetUser(ctx *gin.Context)  {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
        return
    }
    user, err := UserService.GetUser(uid)
    if err != nil {
        panic(err)
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
    middleware.SetCookie(ctx, uid)
    response.Success(ctx)
}

func Self(ctx *gin.Context) {
    rawUid, ok := ctx.Keys["uid"]
    if !ok {
        response.Forbidden(ctx, "please login")
        return
    }
    uid, ok := rawUid.(uint64)
    if !ok {
        response.Forbidden(ctx, "please login")
        return
    }
    user, err := model.GetUser(uid)
    if err != nil {
        panic(err)
    }
    if user == nil {
        response.NotFound(ctx, fmt.Sprintf("user not found: %d", uid))
        return
    }
    response.Entity(ctx, user)
}