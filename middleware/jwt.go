package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/go-gin-demo/utils"
    "github.com/go-gin-demo/utils/response"
    "github.com/dgrijalva/jwt-go"
    "time"
    BizError "github.com/go-gin-demo/errors"
)

const (
    authCookie = "auth"
    expireTime = 24
    jwtSecret = "embarrassing_caanan"
)


func JWT() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token, err := ctx.Cookie(authCookie)
        if err != nil && err.Error() != "http: named cookie not present" {
            response.Forbidden(ctx,"fail to get cookie")
            ctx.Abort()
            return
        }
        if token == "" {
            //response.Forbidden(ctx,"no auth cookie")
            //ctx.Abort()
            ctx.Next()
            return
        }
        uid, err := utils.ParseAuthToken(token, jwtSecret)
        if err != nil {
            if err, ok := err.(*jwt.ValidationError); ok && err.Errors == jwt.ValidationErrorExpired {
                response.Forbidden(ctx,"token expired")
            } else {
                response.Forbidden(ctx,"invalid token")
            }
            ctx.Abort()
        }
        if ctx.Keys == nil {
            ctx.Keys = make(map[string]interface{})
        }
        ctx.Keys["uid"] = uid
        ctx.Next()
    }
}

func SetCurrentUid(ctx *gin.Context, uid uint64) {
    token, err := utils.GenAuthToken(uid, expireTime * time.Hour, jwtSecret)
    if err != nil {
        panic(err)
    }
    ctx.SetCookie(authCookie, token, expireTime * 60 * 60, "/", "", false, false)
}

func GetCurrentUid(ctx *gin.Context) (uint64, error) {
    rawUid, ok := ctx.Keys["uid"]
    if !ok {
        return 0, BizError.Forbidden("login required")
    }
    uid, ok := rawUid.(uint64)
    if !ok {
        return 0, BizError.Forbidden("login required")
    }
    return uid, nil
}