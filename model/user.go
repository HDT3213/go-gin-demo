package model

import (
    "github.com/go-gin-demo/utils"
    "github.com/go-gin-demo/entity"
    "strconv"
    "math/rand"
    "fmt"
    BizError "github.com/go-gin-demo/errors"
)

func genUserCacheKey(uid uint64) string {
    return fmt.Sprintf("U:%d", uid)
}

func setUserCache(user *entity.User) error {
    key := genUserCacheKey(user.ID)
    val, err := Marshal(user)
    if err != nil {
        return err
    }
    _, err = Redis.Set(key, val, 0).Result()
    return err
}

func CreateUser(user *entity.User) error {
    user.ID = utils.Hash64(user.Username + strconv.Itoa(int(utils.Now())) + strconv.Itoa(rand.Int()))
    if db.NewRecord(*user) {
        return BizError.InvalidForm("user exists")
    }
    if err := db.Create(user).Error; err != nil {
        return err
    }
    setUserCache(user)
    return nil
}

func getUserFromCache(uid uint64) (*entity.User, error) {
    key := genUserCacheKey(uid)
    val, err := Redis.Get(key).Result()
    if err != nil {
        if err.Error() == "redis: nil" {
            return nil, nil
        } else {
            panic(err)
        }
    }
    user := new(entity.User)
    err = Unmarshal([]byte(val), user)
    if err != nil {
        panic(err)
    }
    return user, nil
}

func getUserFromDB(uid uint64) (*entity.User, error) {
    user := new(entity.User)
    err := db.Where("valid = 1").First(&user, uid).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return user, err
}

func GetUser(uid uint64) (*entity.User, error) {
    user, err := getUserFromCache(uid)
    if err != nil {
        return nil, err
    }
    if user != nil {
        return user, nil
    }
    user, err = getUserFromDB(uid)
    if err != nil {
        return nil, err
    }
    if user != nil {
        setUserCache(user)
    }
    return user, err
}

func MultiGetUser(uids []uint64) ([]*entity.User, error) {
    if len(uids) == 0 {
        return make([]*entity.User, 0), nil
    }

    keys := make([]string, len(uids))
    for i, uid := range uids {
        keys[i] = genUserCacheKey(uid)
    }

    vals, err := Redis.MGet(keys...).Result()
    if err != nil {
        return nil, err
    }

    users := make([]*entity.User, len(vals))
    for i, val := range vals {
        user := new(entity.User)
        str, ok := val.(string)
        if !ok {
            continue
        }
        err = Unmarshal([]byte(str), user)
        if err != nil {
            continue
        }
        users[i] = user
    }

    for i, uid := range uids {
        if users[i] == nil {
            user, err := getUserFromDB(uid)
            if err != nil {
                continue
            }
            if user != nil {
                users[i] = user
                setUserCache(user)
            }
        }
    }
    return users, nil
}

func GetUserMap(uids []uint64) (map[uint64]*entity.User, error) {
    users, err := MultiGetUser(uids)
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

func GetUserByName(username string) (*entity.User, error) {
    var user entity.User
    err := db.Where("username = ? AND valid = 1", username).First(&user).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return &user, err
}
