package middleware

import (
    "github.com/gin-gonic/gin"
    "go-close/utils"
    "go-close/utils/response"
    "github.com/dgrijalva/jwt-go"
    "time"
)

const (
    authCookie = "auth"
    expireTime = 24
    jwtSecret = "embarrassing_caanan"
)


func JWT() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token, err := ctx.Cookie(authCookie)
        if err != nil {
            response.Forbidden(ctx,"no auth cookie")
            ctx.Abort()
            return
        }
        if token == "" {
            response.Forbidden(ctx,"no auth cookie")
            ctx.Abort()
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

func SetCookie(ctx *gin.Context, uid uint64) {
    token, err := utils.GenAuthToken(uid, expireTime * time.Hour, jwtSecret)
    if err != nil {
        panic(err)
    }
    ctx.SetCookie(authCookie, token, expireTime * 60 * 60, "/", "", false, false)
}