package post

import (
    "fmt"
    "github.com/go-gin-demo/entity"
    "strconv"
    "github.com/go-gin-demo/utils"
    "math/rand"
    "github.com/go-gin-demo/errors"
    "github.com/go-gin-demo/model"
    "github.com/go-gin-demo/model/counter"
    "github.com/go-gin-demo/utils/collections"
)

const counterKeyPrefix = "Count:P:"

func genCacheKey(pid uint64) string {
    return fmt.Sprintf("P:%d", pid)
}

func setCache(post *entity.Post) error {
    key := genCacheKey(post.ID)
    val, err := model.Marshal(post)
    if err != nil {
        return err
    }
    _, err = model.Redis.Set(key, val, 0).Result()
    return err
}

func Create(post *entity.Post) error {
    post.ID = utils.Hash64(strconv.FormatUint(post.Uid, 10) + strconv.Itoa(int(utils.Now())) + strconv.Itoa(rand.Int()))
    if model.DB.NewRecord(post) {
        return errors.InvalidForm("post exists")
    }
    if err := model.DB.Create(post).Error; err != nil {
        return err
    }
    setCache(post)
    counter.Increase(counterKeyPrefix, post.Uid, 1)
    return nil
}

func getFromCache(pid uint64) (*entity.Post, error) { // may get a post which valid=false
    key := genCacheKey(pid)
    val, err := model.Redis.Get(key).Result()
    if err != nil {
        if err.Error() == "redis: nil" {
            return nil, nil
        } else {
            return nil, err
        }
    }
    post := new(entity.Post)
    err = model.Unmarshal([]byte(val), post)
    if err != nil {
        return nil, err
    }
    return post, nil
}

func getFromDB(pid uint64) (*entity.Post, error) {
    post := new(entity.Post)
    err := model.DB.Where("valid = 1").First(&post, pid).Error
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

    vals, err := model.Redis.MGet(keys...).Result()
    if err != nil {
        return nil, err
    }

    posts := make([]*entity.Post, len(vals))
    for i, val := range vals {
        post := new(entity.Post)
        str, ok := val.(string)
        if !ok {
            continue
        }
        err = model.Unmarshal([]byte(str), post)
        if err != nil {
            continue
        }
        posts[i] = post
    }

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
    if err := model.DB.Model(&entity.Post{}).Where("id = ? AND valid = 1", pid).Update("valid", 0).Error; err != nil {
        return err
    }
    err := setCache(&entity.Post{ID:pid, Valid:false})
    counter.Increase(counterKeyPrefix, post.Uid, -1)
    return err
}

func getUserPostCountFromDB(uid uint64) (int32, error) {
    var count int32
    err := model.DB.Model(&entity.Post{}).Where("uid = ? AND valid = 1", uid).Count(&count).Error
    if err != nil {
        return -1, err
    }
    return count, nil
}

func GetUserPostCount(uid uint64) (int32, error){
    return counter.Get(counterKeyPrefix, uid, getUserPostCountFromDB)
}

func multiGetPostCountFromDB(uids []uint64) (map[uint64]int32, error) {
    pairs := make([]*collections.IdCountPair, len(uids))
    err := model.DB.Table("posts").Select("uid AS id, count(*) AS num").Where("uid IN (?) AND valid = 1", uids).Group("uid").Scan(&pairs).Error
    if err != nil {
        return nil, err
    }
    return collections.ToCountMap(pairs), nil
}

func GetUserPostCountMap(uids []uint64) (map[uint64]int32, error) {
    return counter.GetMap(counterKeyPrefix, uids, multiGetPostCountFromDB)
}