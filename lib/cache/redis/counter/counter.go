package counter

import (
    "fmt"
    "strconv"
    "time"
    "github.com/go-gin-demo/lib/collections/set"
    "github.com/go-redis/redis"
)

const TTL = 1 * time.Hour

func genKey(keyPrefix string, id uint64) string {
    return fmt.Sprintf("%s:%d", keyPrefix, id)
}

func Get(red *redis.Client, keyPrefix string, id uint64, getter func (uint64)(int32, error)) (int32, error) {
    key := genKey(keyPrefix, id)
    result, err := red.Get(key).Result()
    if err != nil {
        if err.Error() == "redis: nil" {
            count, err := getter(id) // trust getter
            if err != nil {
                return -1, err
            }
            Set(red, keyPrefix, id, count)
            return count, nil
        } else {
            return -1, err
        }
    }
    count, err := strconv.ParseInt(result, 10, 32)
    if err != nil {
        return -1, err
    }
    if count < 0 { // cache inconsistent
        count, err := getter(id) // trust getter
        if err != nil {
            return -1, err
        }
        Del(red, keyPrefix, id) // must not set cache, other thread may about to increase count
        return count, nil
    }
    return int32(count), nil
}

func Set(red *redis.Client, keyPrefix string, id uint64, count int32) error {
    key := genKey(keyPrefix, id)
    _, err := red.Set(key, count, TTL).Result()
    return err

}

func Del(red *redis.Client, keyPrefix string, id uint64) error {
    key := genKey(keyPrefix, id)
    _, err := red.Del(key).Result()
    return err
}

func Increase(red *redis.Client, keyPrefix string, id uint64, delta int32) error {
    key := genKey(keyPrefix, id)
    exists, err := red.Exists(key).Result()
    if err != nil {
        return err
    }
    if exists > 0 {
        _, err := red.IncrBy(key, int64(delta)).Result()
        if err != nil {
            return err
        }
    }
    return nil
}

func GetMap(red *redis.Client, keyPrefix string, ids []uint64, getter func ([]uint64)(map[uint64]int32, error)) (map[uint64]int32, error){
    size := len(ids)
    keys := make([]string, size)
    for i, id := range ids {
        keys[i] = genKey(keyPrefix, id)
    }
    vals, err := red.MGet(keys...).Result()
    if err != nil {
        return nil, err
    }

    countMap := make(map[uint64]int32)
    failedIdSet := set.MakeUint64Set()
    for i, id := range ids {
        val, ok := vals[i].(string)
        if !ok {
            failedIdSet.Add(id)
            continue
        }
        count, err := strconv.ParseInt(val, 10, 32)
        if err != nil {
            failedIdSet.Add(id)
            continue
        }
        countMap[id] = int32(count)
    }

    if failedIdSet.Len() > 0 {
        fillMap, err := getter(failedIdSet.ToArray())
        if err != nil {
            return countMap, err
        }
        SetMap(red, keyPrefix, fillMap)
        for id, count := range fillMap {
            countMap[id] = count
        }
    }
    return countMap, nil
}

func SetMap(red *redis.Client, keyPrefix string, countMap map[uint64]int32) error {
    size := len(countMap)
    pairs := make([]interface{}, size * 2)
    i := 0
    for id, count := range countMap {
        pairs[i] = genKey(keyPrefix, id)
        pairs[i+1] = count
        i += 2
    }
    _, err := red.MSet(pairs...).Result()
    return err
}