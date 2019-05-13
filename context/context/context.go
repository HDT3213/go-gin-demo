// runtime status and resources, depends on lib only
package context

import (
    "github.com/go-gin-demo/lib/mq"
    "github.com/jinzhu/gorm"
    "github.com/go-redis/redis"
)

var MQ mq.MQ

var DB *gorm.DB

var Redis *redis.Client