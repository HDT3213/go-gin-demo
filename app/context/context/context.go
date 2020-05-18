// runtime status and resources, depends on lib only
package context

import (
    "github.com/HDT3213/go-gin-demo/lib/canal"
    "github.com/HDT3213/go-gin-demo/lib/mq"
    "github.com/go-redis/redis"
    "github.com/jinzhu/gorm"
)

type RuntimeSettings struct {
    EnableCanal bool `yaml:"canal-enable"`
    AllowMQFallback bool `yaml:"mq-fallback"`
}

var Runtime RuntimeSettings

var MQ mq.MQ

var DB *gorm.DB

var Redis *redis.Client

var Canal *canal.Canal

func EnableCanal() bool {
    return Runtime.EnableCanal
}