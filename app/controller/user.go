package controller

import (
    "github.com/HDT3213/go-gin-demo/app/middleware/auth"
    UserService "github.com/HDT3213/go-gin-demo/app/service/user"
    BizError "github.com/HDT3213/go-gin-demo/lib/errors"
    "github.com/HDT3213/go-gin-demo/lib/request"
    "github.com/HDT3213/go-gin-demo/lib/response"
    "github.com/gin-gonic/gin"
    "strconv"
)

func GetUser(ctx *gin.Context)  {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
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

    user, err := UserService.GetUser(currentUid, uid)
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
    auth.SetCurrentUid(ctx, entity.ID)
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
    auth.SetCurrentUid(ctx, uid)
    response.Success(ctx)
}

func Self(ctx *gin.Context) {
    uid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    user, err := UserService.GetUser(uid, uid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entity(ctx, user)
}

func Follow(ctx *gin.Context) {
    currentUid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
        return
    }
    err = UserService.Follow(currentUid, uid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
    return
}

func UnFollow(ctx *gin.Context) {
    currentUid, err := auth.GetCurrentUid(ctx)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
        return
    }
    err = UserService.UnFollow(currentUid, uid)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Success(ctx)
    return
}

func GetUserFollowings(ctx *gin.Context) {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
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

    users, total, err := UserService.GetFollowings(currentUid, uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, users, total)
}

func GetSelfFollowings(ctx *gin.Context) {
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

    users, total, err := UserService.GetFollowings(currentUid, currentUid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, users, total)
}

func GetUserFollowers(ctx *gin.Context) {
    uid, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
    if err != nil {
        response.BadRequest(ctx, "invalid uid: " + ctx.Param("id"))
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

    users, total, err := UserService.GetFollowers(currentUid, uid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, users, total)
}

func GetSelfFollowers(ctx *gin.Context) {
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

    users, total, err := UserService.GetFollowers(currentUid, currentUid, start, length)
    if err != nil {
        response.Error(ctx, err)
        return
    }
    response.Entities(ctx, users, total)
}