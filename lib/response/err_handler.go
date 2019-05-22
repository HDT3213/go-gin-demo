package response

import (
    "github.com/gin-gonic/gin"
    BizError "github.com/go-gin-demo/lib/errors"
)

func Error(ctx *gin.Context, err error) {
    if BizError.IsInvalidForm(err) {
        BadRequest(ctx, err.Error())
    } else if BizError.IsNotFound(err) {
        NotFound(ctx, err.Error())
    } else if BizError.IsForbidden(err) {
        Forbidden(ctx, err.Error())
    } else {
        InternalServerError(ctx, err.Error())
    }
}