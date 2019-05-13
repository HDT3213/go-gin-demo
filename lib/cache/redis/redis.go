package redis

import (
    "github.com/go-redis/redis"
    "time"
    "github.com/vmihailenco/msgpack"
    "reflect"
    "github.com/go-gin-demo/lib/logger"
)

type Settings struct {
    Host string `yaml:"host"`
    Password string `yaml:"password"`
    DB int `yaml:"db"`
    PoolSize int `yaml:"pool-size"`
    PoolTimeout time.Duration `yaml:"pool-timeout"`
}

func Setup(settings *Settings) (*redis.Client, error) {
    redis := redis.NewClient(&redis.Options{
        Addr:     settings.Host,
        Password: settings.Password,
        DB:       settings.DB,
        PoolSize:     settings.PoolSize,
        PoolTimeout:  settings.PoolTimeout,
    })

    _, err := redis.Ping().Result()
    if err != nil {
        return nil, err
    }
    return redis, nil
}


func Marshal(v interface{}) ([]byte, error) {
    return msgpack.Marshal(v)
    //return json.Marshal(v)
}

func Unmarshal(data []byte, v interface{}) error {
    return msgpack.Unmarshal(data, v)
    //return json.Unmarshal(data, v)
}

func MultiUnmarshal(vals []interface{}, out interface{}) { // out should be a *[]*T, such as &[]*entity.Post
    slice := reflect.ValueOf(out).Elem()
    if slice.Kind() != reflect.Slice {
        panic("out is not slice")
    }
    if slice.Cap() < len(vals) {
        panic("out is too small")
    }
    slice.SetLen(len(vals))

    elemType := reflect.TypeOf(slice.Interface()).Elem().Elem() // slice.Elem() is *T, elemType is T

    for i, val := range vals {
        elem := reflect.New(elemType).Interface() // elem is *T
        str, ok := val.(string)
        if !ok {
            continue
        }
        err := Unmarshal([]byte(str), elem)
        if err != nil {
            logger.Warn("unmarshal failed, raw: " + str + ", err: " + err.Error())
            continue
        }
        slice.Index(i).Set(reflect.ValueOf(elem))
    }
}

func MultiUnmarshalStr(vals []string, out interface{}) { // out should be a *[]*T, such as &[]*entity.Post
    slice := reflect.ValueOf(out).Elem()
    if slice.Kind() != reflect.Slice {
        panic("out is not slice")
    }
    if slice.Cap() < len(vals) {
        panic("out is too small")
    }
    slice.SetLen(len(vals))

    elemType := reflect.TypeOf(slice.Interface()).Elem().Elem() // slice.Elem() is *T, elemType is T

    for i, val := range vals {
        elem := reflect.New(elemType).Interface() // elem is *T
        err := Unmarshal([]byte(val), elem)
        if err != nil {
            continue
        }
        slice.Index(i).Set(reflect.ValueOf(elem))
    }
}