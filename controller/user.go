package controller

import (
    "github.com/gin-gonic/gin"
    "go-close/model"
    "strconv"
    "fmt"
    "go-close/middleware"
    "go-close/utils/response"
)

func GetUser(ctx *gin.Context)  {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
        return
    }
    user := model.GetUser(uid)
    if user == nil {
        response.NotFound(ctx, fmt.Sprintf("user not found: %d", uid))
        return
    }
    response.Entity(ctx, user)
}

func AllUser(ctx *gin.Context) {
    users := model.AllUsers()
    response.Entity(ctx, users)
}

func Register(ctx *gin.Context) {
    user := new(model.User)
    user.Username = ctx.PostForm("username")
    user.Password = ctx.PostForm("password")
    model.CreateUser(user)
    response.Entity(ctx, user)
}

func Login(ctx *gin.Context) {
    username := ctx.PostForm("username")
    password := ctx.PostForm("password")
    user, err := model.GetUserByName(username)
    if err != nil {
        panic(err)
    }
    if user != nil && user.Password == password {
        middleware.SetCookie(ctx, user.ID)
        response.Success(ctx)
    } else {
        response.BadRequest(ctx, "invalid username or password")
    }

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
    user := model.GetUser(uid)
    if user == nil {
        response.NotFound(ctx, fmt.Sprintf("user not found: %d", uid))
        return
    }
    response.Entity(ctx, user)
}