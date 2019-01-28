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