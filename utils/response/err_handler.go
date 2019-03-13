package response

import (
    "github.com/gin-gonic/gin"
    BizError "go-close/errors"
)

func Error(ctx *gin.Context, err error) {
    if BizError.IsInvalidForm(err) {
        BadRequest(ctx, err.Error())
    } else if BizError.IsNotFound(err) {
        NotFound(ctx, err.Error())
    }
}
