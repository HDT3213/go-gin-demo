package post

import (
    "fmt"
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/app/entity"
    "github.com/HDT3213/go-gin-demo/lib/cache/redis"
    "github.com/HDT3213/go-gin-demo/lib/cache/redis/counter"
    "github.com/HDT3213/go-gin-demo/lib/collections"
    "github.com/HDT3213/go-gin-demo/lib/errors"
    "github.com/HDT3213/go-gin-demo/lib/hash"
    "github.com/HDT3213/go-gin-demo/lib/time"
    "math/rand"
    "strconv"
)

const userPostCounterKeyPrefix = "Count:User:P"

func genCacheKey(pid uint64) string {
    return fmt.Sprintf("P:%d", pid)
}

func setCache(post *entity.Post) error {
    key := genCacheKey(post.ID)
    val, err := redis.Marshal(post)
    if err != nil {
        return err
    }
    _, err = context.Redis.Set(key, val, 0).Result()
    return err
}

func Create(post *entity.Post) error {
    post.ID = uint64(hash.Hash32(strconv.FormatUint(post.Uid, 10) + strconv.Itoa(int(time.Now())) + strconv.Itoa(rand.Int())))
    if context.DB.NewRecord(post) {
        return errors.InvalidForm("post exists")
    }
    if err := context.DB.Create(post).Error; err != nil {
        return err
    }
    if !context.EnableCanal() {
       return AfterCreate(post)
    }
    return nil
}

func AfterCreate(post *entity.Post) error {
    err := setCache(post)
    if err != nil {
        return err
    }
    return counter.Increase(context.Redis, userPostCounterKeyPrefix, post.Uid, 1)
}

func getFromCache(pid uint64) (*entity.Post, error) { // may get a post which valid=false
    key := genCacheKey(pid)
    val, err := context.Redis.Get(key).Result()
    if err != nil {
        if err.Error() == "redis: nil" {
            return nil, nil
        } else {
            return nil, err
        }
    }
    post := new(entity.Post)
    err = redis.Unmarshal([]byte(val), post)
    if err != nil {
        return nil, err
    }
    return post, nil
}

func getFromDB(pid uint64) (*entity.Post, error) {
    post := new(entity.Post)
    err := context.DB.Where("valid = 1").First(&post, pid).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return post, err
}

func Get(pid uint64) (*entity.Post, error) {
    post, err := getFromCache(pid)
    if err != nil { // must be caused by redis or system error, abort
        return nil, err
    }
    if post != nil {
        if !post.Valid {
            return nil, nil
        }
        return post, nil
    }
    post, err  = getFromDB(pid)
    if err != nil {
        return nil, err
    }
    if post != nil {
        setCache(post)  // set invalid post as placeholder
        if !post.Valid {
            return nil, nil
        }
    }
    return post, nil
}

func MultiGet(pids []uint64) ([]*entity.Post, error) {
    if len(pids) == 0 {
        return make([]*entity.Post, 0), nil
    }

    keys := make([]string, len(pids))
    for i, uid := range pids {
        keys[i] = genCacheKey(uid)
    }

    vals, err := context.Redis.MGet(keys...).Result()
    if err != nil {
        return nil, err
    }
    posts := make([]*entity.Post, len(vals))
    redis.MultiUnmarshal(vals, &posts)

    for i, pid := range pids {
        if posts[i] == nil {
            post, err := getFromDB(pid)
            if err != nil {
                continue
            }
            if post != nil {
                setCache(post) // set invalid post as placeholder
                if post.Valid {
                    posts[i] = post
                }
            }
        }
    }
    return posts, nil
}

func GetMap(pids []uint64) (map[uint64]*entity.Post, error) {
    posts, err := MultiGet(pids)
    if err != nil {
        return nil, err
    }
    postMap := make(map[uint64]*entity.Post, len(pids))
    for _, post := range posts {
        postMap[post.ID] = post
    }
    return postMap, nil
}

func Delete(post *entity.Post) error {
    pid := post.ID
    if err := context.DB.Model(&entity.Post{}).Where("id = ? AND valid = 1", pid).Update("valid", 0).Error; err != nil {
        return err
    }
    if !context.EnableCanal() {
        return AfterDelete(post)
    }
    return nil
}

func AfterDelete(post *entity.Post) error {
    err := setCache(&entity.Post{ID:post.ID, Valid:false})
    if err != nil {
        return err
    }
    return counter.Increase(context.Redis, userPostCounterKeyPrefix, post.Uid, -1)
}

func getUserPostCountFromDB(uid uint64) (int32, error) {
    var count int32
    err := context.DB.Model(&entity.Post{}).Where("uid = ? AND valid = 1", uid).Count(&count).Error
    if err != nil {
        return -1, err
    }
    return count, nil
}

func GetUserPostCount(uid uint64) (int32, error){
    return counter.Get(context.Redis, userPostCounterKeyPrefix, uid, getUserPostCountFromDB)
}

func multiGetPostCountFromDB(uids []uint64) (map[uint64]int32, error) {
    pairs := make([]*collections.IdCountPair, len(uids))
    err := context.DB.Table("posts").Select("uid AS id, count(*) AS num").Where("uid IN (?) AND valid = 1", uids).Group("uid").Scan(&pairs).Error
    if err != nil {
        return nil, err
    }
    return collections.ToCountMap(pairs), nil
}

func GetUserPostCountMap(uids []uint64) (map[uint64]int32, error) {
    return counter.GetMap(context.Redis, userPostCounterKeyPrefix, uids, multiGetPostCountFromDB)
}