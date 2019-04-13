package follow

import (
    "fmt"
    "github.com/go-gin-demo/entity"
    "github.com/go-gin-demo/model"
    "github.com/go-gin-demo/model/counter"
    "github.com/go-gin-demo/utils/collections/set"
    RLock "github.com/bsm/redis-lock"
    "errors"
    "github.com/go-gin-demo/utils/collections"
)

const (
    userFollowingCounterKeyPrefix = "Count:User:Following"
    userFollowerCounterKeyPrefix = "Count:User:Follower"
)

// uid -> set(FollowingUid)
func genKey(uid uint64) string {
    return fmt.Sprintf("U:Following:%d", uid)
}

func cached(uid uint64) (bool, error) {
    followingCount, err := GetFollowingCount(uid)
    if err != nil {
        return false, err
    }
    if followingCount == 0 { // if following no one, consider empty key as cached too
        return true, nil
    }

    key := genKey(uid)
    exists, err := model.Redis.Exists(key).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, nil
}

func setCache(follow *entity.Follow) error {
    uid := follow.Uid
    exists, err := cached(uid)
    if err != nil {
        return err
    }
    if !exists {
        return nil  // lazy load
    }
    key := genKey(uid)
    _, err = model.Redis.SAdd(key, follow.FollowingUid).Result()
    return err
}

func Create(follow *entity.Follow) error {
    // check existed record
    existed := new(entity.Follow)
    err := model.DB.Where("uid = ? AND following_uid = ?", follow.Uid, follow.FollowingUid).First(&existed).Error
    if err != nil {
        if err.Error() == "record not found" {
            existed = nil
        } else {
            return err
        }
    }
    if existed != nil {
        if existed.Valid {
            return nil // exists do nothing
        } else {
            // existed but deleted, recover
            if err := model.DB.Model(&entity.Follow{}).
                Where("uid = ? AND following_uid = ?", follow.Uid, follow.FollowingUid).
                Update("valid", true).Error; err != nil {
                return err
            }
        }
    } else {
        // new record
        if err := model.DB.Create(follow).Error; err != nil {
            return err
        }
    }

    setCache(follow)
    counter.Increase(userFollowingCounterKeyPrefix, follow.Uid, 1)
    counter.Increase(userFollowerCounterKeyPrefix, follow.FollowingUid, 1)
    return nil
}

func removeCache(uid uint64, followingUid uint64) error {
    key := genKey(uid)
    _, err := model.Redis.SRem(key, followingUid).Result()
    return err
}

func Delete(uid uint64, followingUid uint64) error {
    err := model.DB.Model(&entity.Follow{}).Where("uid = ? AND following_uid = ?", uid, followingUid).Update("valid", 0).Error
    if err != nil {
        return err
    }
    removeCache(uid, followingUid)
    counter.Increase(userFollowingCounterKeyPrefix, uid, -1)
    counter.Increase(userFollowerCounterKeyPrefix, followingUid, -1)
    return nil
}

func rebuildInternal(uid uint64) (*set.Uint64Set, error) {
    var follows []*entity.Follow
    if err := model.DB.Where("uid = ? AND valid = 1", uid).Find(&follows).Error; err != nil {
        return nil, err
    }
    followingSet := set.MakeUint64Set()
    for _, follow := range follows {
        followingSet.Add(follow.FollowingUid)
    }
    followingUids := followingSet.ToInterfaceArray()
    _, err := model.Redis.SAdd(genKey(uid), followingUids...).Result()
    if err != nil {
        return nil, err
    }
    return followingSet, nil
}

func Rebuild(uid uint64) error {
    key := genKey(uid)

    // lock
    lock, err := RLock.Obtain(model.Redis, "lock:" + key, &RLock.Options{
        RetryCount: 3,
    })
    if err != nil {
        return  err
    }
    if lock == nil {
        return errors.New("cannot obtain lock")
    }
    defer lock.Unlock()

    // check again
    existed, err := cached(uid)
    if err != nil {
        return err
    }
    if existed {
        return nil
    }

    _, err = rebuildInternal(uid)
    return err
}

func IsFollowing(uid uint64, followingUid uint64) (bool, error) {
    // check and rebuild
    isCached, err := cached(uid)
    if err != nil {
        return false, err
    }
    if !isCached {
        err = Rebuild(uid)
        if err != nil {
            return false, err
        }
    }
    key := genKey(uid)
    result, err := model.Redis.SIsMember(key, followingUid).Result()
    if err != nil {
        return false, nil
    }
    return result, nil
}


func GetFollowingIn(currentUid uint64, uids []uint64) ([]uint64, error) {
    // check an rebuild
    isCached, err := cached(currentUid)
    if err != nil {
        return nil, err
    }
    if !isCached {
        err = Rebuild(currentUid)
        if err != nil {
            return nil, err
        }
    }

    key := genKey(currentUid)
    return model.Intersect(key, uids)
}

func getFollowingCountFromDB(uid uint64) (int32, error) {
    var count int32
    err := model.DB.Model(&entity.Follow{}).Where("uid = ? AND valid = 1", uid).Count(&count).Error
    if err != nil {
        return -1, err
    }
    return count, nil
}

func GetFollowingCount(uid uint64) (int32, error) {
    return counter.Get(userFollowingCounterKeyPrefix, uid, getFollowingCountFromDB)
}

func multiGetFollowingCountFromDB(uids []uint64) (map[uint64]int32, error) {
    pairs := make([]*collections.IdCountPair, len(uids))
    err := model.DB.Model(&entity.Follow{}).Select("uid AS id, count(*) AS num").Where("uid IN (?) AND valid = 1", uids).Group("uid").Scan(&pairs).Error
    if err != nil {
        return nil, err
    }
    return collections.ToCountMap(pairs), nil
}

func GetFollowingCountMap(uids []uint64) (map[uint64]int32, error) {
    return counter.GetMap(userFollowingCounterKeyPrefix, uids, multiGetFollowingCountFromDB)
}

func getFollowerCountFromDB(uid uint64) (int32, error) {
    var count int32
    err := model.DB.Model(&entity.Follow{}).Where("following_uid = ? AND valid = 1", uid).Count(&count).Error
    if err != nil {
        return -1, err
    }
    return count, nil
}

func GetFollowerCount(uid uint64) (int32, error) {
    return counter.Get(userFollowerCounterKeyPrefix, uid, getFollowerCountFromDB)
}


func multiGetFollowerCountFromDB(uids []uint64) (map[uint64]int32, error) {
    pairs := make([]*collections.IdCountPair, len(uids))
    err := model.DB.Model(&entity.Follow{}).Select("following_uid AS id, count(*) AS num").Where("following_uid in (?) AND valid = 1", uids).Group("following_uid").Scan(&pairs).Error
    if err != nil {
        return nil, err
    }
    return collections.ToCountMap(pairs), nil
}

func GetFollowerCountMap(uids []uint64) (map[uint64]int32, error) {
    return counter.GetMap(userFollowerCounterKeyPrefix, uids, multiGetFollowerCountFromDB)
}

func GetFollowings(uid uint64, start int32, length int32) ([]*entity.Follow, error) {
    var follows []*entity.Follow
    err := model.DB.
        Where("uid = ? AND valid = 1", uid).
        Order("created_at DESC").
        Limit(length).Offset(start).
        Find(&follows).Error
    if err != nil {
        return nil, err
    }
    return follows, nil
}

func GetFollowers(uid uint64, start int32, length int32) ([]*entity.Follow, error) {
    var follows []*entity.Follow
    err := model.DB.
        Where("following_uid = ? AND valid = 1", uid).
        Order("created_at DESC").
        Limit(length).Offset(start).
        Find(&follows).Error
    if err != nil {
        return nil, err
    }
    return follows, nil
}