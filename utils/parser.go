package utils

import (
    "github.com/gin-gonic/gin"
    "strconv"
    BizError "github.com/go-gin-demo/errors"
)

func GetPage(ctx *gin.Context, defaultStart int32, defaultLength int32) (int32, int32, error) {
    var err error

    rawStart := ctx.Query("start")
    var start int64
    if len(rawStart) > 0 {
        start, err = strconv.ParseInt(rawStart, 10, 32)
        if err != nil {
            return -1, -1, BizError.InvalidForm( "invalid start: " + rawStart)
        }
    } else {
        start = int64(defaultStart)
    }

    rawLength := ctx.Query("length")
    var length int64
    if len(rawLength) > 0 {
        length, err = strconv.ParseInt(rawLength, 10, 32)
        if err != nil {
            return -1, -1, BizError.InvalidForm( "invalid length: " + rawLength)
        }
    } else {
        length = int64(defaultLength)
    }
    return int32(start), int32(length), nil
}
