package utils

import "time"

func Now() int32 {
   return int32(time.Now().Unix())
}