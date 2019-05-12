package config

import (
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "github.com/go-gin-demo/utils/logger"
    "github.com/go-gin-demo/model"
    MQCore "github.com/go-gin-demo/mq/core"
)

type Settings struct {
   Log logger.Settings
   DB model.DBSettings
   Redis model.RedisSettings
   Rabbit MQCore.Settings
   Canal model.CanalSettings
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