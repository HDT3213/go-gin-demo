// static system config, depends on lib only
package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "github.com/go-gin-demo/lib/logger"
    "github.com/go-gin-demo/lib/mq"
    "github.com/go-gin-demo/lib/db"
    "github.com/go-gin-demo/lib/cache/redis"
    "github.com/go-gin-demo/lib/canal"
    "github.com/go-gin-demo/context/context"
)

type Settings struct {
   Runtime context.RuntimeSettings
   Log logger.Settings
   DB db.Settings
   Redis redis.Settings
   Rabbit mq.Settings
   Canal canal.Settings
}

func Setup(path string) (*Settings) {
    bytes, err := ioutil.ReadFile(path)
    if err != nil {
        panic(err)
    }
    settings := Settings{}
    err = yaml.Unmarshal(bytes, &settings)
    if err != nil {
        panic(err)
    }
    return &settings
}