package user

import (
    "github.com/go-gin-demo/entity"
    "strconv"
    "math/rand"
    "fmt"
    BizError "github.com/go-gin-demo/lib/errors"
    "github.com/go-gin-demo/lib/hash"
    "github.com/go-gin-demo/lib/time"
    "github.com/go-gin-demo/lib/cache/redis"
    "github.com/go-gin-demo/context/context"
)

func genCacheKey(uid uint64) string {
    return fmt.Sprintf("U:%d", uid)
}

func setCache(user *entity.User) error {
    key := genCacheKey(user.ID)
    val, err := redis.Marshal(user)
    if err != nil {
        return err
    }
    _, err = context.Redis.Set(key, val, 0).Result()
    return err
}

func Create(user *entity.User) error {
    user.ID = uint64(hash.Hash32(user.Username + strconv.Itoa(int(time.Now())) + strconv.Itoa(rand.Int())))
    if context.DB.NewRecord(*user) {
        return BizError.InvalidForm("user exists")
    }
    if err := context.DB.Create(user).Error; err != nil {
        return err
    }
    setCache(user)
    return nil
}

func getFromCache(uid uint64) (*entity.User, error) {
    key := genCacheKey(uid)
    val, err := context.Redis.Get(key).Result()
    if err != nil {
        if err.Error() == "redis: nil" {
            return nil, nil
        } else {
            return nil, err
        }
    }
    user := new(entity.User)
    err = redis.Unmarshal([]byte(val), user)
    if err != nil {
        return nil, err
    }
    return user, nil
}

func getFromDB(uid uint64) (*entity.User, error) {
    user := new(entity.User)
    err := context.DB.Where("valid = 1").First(&user, uid).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return user, err
}

func Get(uid uint64) (*entity.User, error) {
    user, err := getFromCache(uid)
    if err != nil {
        return nil, err
    }
    if user != nil {
        return user, nil
    }
    user, err = getFromDB(uid)
    if err != nil {
        return nil, err
    }
    if user != nil {
        setCache(user)
    }
    return user, err
}

/**
    return slice may contains nil
 */
func MultiGet(uids []uint64) ([]*entity.User, error) {
    if len(uids) == 0 {
        return make([]*entity.User, 0), nil
    }

    keys := make([]string, len(uids))
    for i, uid := range uids {
        keys[i] = genCacheKey(uid)
    }

    vals, err := context.Redis.MGet(keys...).Result()
    if err != nil {
        return nil, err
    }

    users := make([]*entity.User, len(vals))
    redis.MultiUnmarshal(vals, &users)

    for i, uid := range uids {
        if users[i] == nil {
            user, err := getFromDB(uid)
            if err != nil {
                continue
            }
            if user != nil {
                users[i] = user
                setCache(user)
            }
        }
    }
    return users, nil
}

func GetMap(uids []uint64) (map[uint64]*entity.User, error) {
    users, err := MultiGet(uids)
    if err != nil {
        return nil, err
    }
    userMap := make(map[uint64]*entity.User, len(uids))
    for _, user := range users {
        if user != nil {
            userMap[user.ID] = user
        }
    }
    return userMap, nil
}

func GetByName(username string) (*entity.User, error) {
    var user entity.User
    err := context.DB.Where("username = ? AND valid = 1", username).First(&user).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return &user, err
}
