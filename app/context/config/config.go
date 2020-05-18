// static system config, depends on lib only
package config

import (
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/lib/cache/redis"
    "github.com/HDT3213/go-gin-demo/lib/canal"
    "github.com/HDT3213/go-gin-demo/lib/db"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    "github.com/HDT3213/go-gin-demo/lib/mq"
    "gopkg.in/yaml.v2"
    "io/ioutil"
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