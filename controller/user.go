package controller

import (
    "github.com/gin-gonic/gin"
    "../model"
    "strconv"
    "fmt"
)

func GetUser(context *gin.Context)  {
    uid, err := strconv.ParseUint(context.Param("id"), 10, 32)
    if err != nil {
        BadRequest(context, "invalid uid: " + context.Param("id"))
        return
    }
    user := model.GetUser(uid)
    if user == nil {
        NotFound(context, fmt.Sprintf("user not found: %d", uid))
        return
    }
    Entity(context, user)
}

func AllUser(context *gin.Context) {
    users := model.AllUsers()
    Entity(context, users)
}

func Register(context *gin.Context) {
    user := new(model.User)
    user.Username = context.PostForm("username")
    user.Password = context.PostForm("password")
    model.CreateUser(user)
    Entity(context, user)
}