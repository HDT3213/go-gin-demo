// runtime status and resources, depends on lib only
package context

import (
    "github.com/go-gin-demo/lib/mq"
    "github.com/jinzhu/gorm"
    "github.com/go-redis/redis"
    "github.com/go-gin-demo/lib/canal"
)

type RuntimeSettings struct {
    EnableCanal bool `yaml:"canal-enable"`
}

var Runtime RuntimeSettings

var MQ mq.MQ

var DB *gorm.DB

var Redis *redis.Client

var Canal *canal.Canal

func EnableCanal() bool {
    return Runtime.EnableCanal
}