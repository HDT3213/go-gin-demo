package redis

import (
    "fmt"
    "github.com/go-redis/redis"
    "math/rand"
    "strconv"
)

const lockSuffix  = ":lock"

/**
 * gen a temp key which in same slot as originKey
 * see: https://www.cnblogs.com/Finley/p/10674101.html#%E4%B8%B4%E6%97%B6%E9%94%AE%E7%9A%84%E7%94%9F%E6%88%90
 */
func GenTempKey(red *redis.Client, originKey string, suffix string) (string, error) {
    for {
        // use tmp suffix to avoid collision
        tmpKey := fmt.Sprintf("{%s}%s#%d", originKey, suffix, rand.Int() % 1000000)
        tmpKeyLock := tmpKey + lockSuffix
        // try to setnx tmpKeyLock, if success means lock the temp key
        locked, err := red.SetNX(tmpKeyLock, "1", 0).Result()
        if err != nil {
            return "", err
        }
        if locked {
            return tmpKey, nil
        }
    }
}

func GenConsumeKey(red *redis.Client, originKey string, suffix string) (string, error) {
    for i := 0; ; i++ {
        // use tmp suffix to avoid collision
        tmpKey := fmt.Sprintf("{%s}%s#%d", originKey, suffix, i)
        tmpKeyLock := tmpKey + lockSuffix
        // try to setnx tmpKeyLock, if success means lock the temp key
        locked, err := red.SetNX(tmpKeyLock, "1", 0).Result()
        if err != nil {
            return "", err
        }
        if locked {
            return tmpKey, nil
        }
    }
}

func ReleaseTempKey(red *redis.Client, key string) error {
    lockKey := key + lockSuffix
    _, err := red.Del(key, lockKey).Result()
    return err
}

/**
 * 与 Redis 中的 Set 取交集
 * 将 arr 放入临时键中再使用 SINTER 命令
 */
func Intersect(red *redis.Client, key string, arr []uint64) ([]uint64, error) {
    // gen temp key
    tmpKey, err := GenTempKey(red, key,"inter")
    if err != nil {
        return nil, err
    }
    defer ReleaseTempKey(red, tmpKey)

    // put to temp key
    members := make([]interface{}, len(arr))
    for i, m := range arr {
        members[i] = m
    }
    _, err = red.SAdd(tmpKey, members...).Result()
    if err != nil {
        return nil, err
    }

    // SINTER
    vals, err := red.SInter(key, tmpKey).Result()
    if err != nil {
        return nil, err
    }

    // unmarshal results
    results := make([]uint64, len(vals))
    for i, val := range vals {
        uintVal, _ := strconv.ParseUint(val, 10, 64)
        results[i] = uintVal
    }
    return results, nil


}
