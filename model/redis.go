package model

import (
    "github.com/go-redis/redis"
    "fmt"
    "time"
    "github.com/vmihailenco/msgpack"
    "reflect"
)

var Redis *redis.Client

func setupRedis() {
    var err error
    Redis = redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
        PoolSize:     16,
        PoolTimeout:  10 * time.Second,
    })

    pong, err := Redis.Ping().Result()
    fmt.Println(pong, err)
}

func closeCache() {
    if Redis != nil {
        Redis.Close()
    }
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