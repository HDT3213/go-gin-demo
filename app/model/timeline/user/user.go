package user

import (
    "errors"
    "fmt"
    "github.com/HDT3213/go-gin-demo/app/context/context"
    "github.com/HDT3213/go-gin-demo/app/entity"
    PostModel "github.com/HDT3213/go-gin-demo/app/model/post"
    "github.com/HDT3213/go-gin-demo/lib/cache/redis"
    "github.com/HDT3213/go-gin-demo/lib/collections/set"
    "github.com/HDT3213/go-gin-demo/lib/logger"
    RLock "github.com/bsm/redis-lock"
    GoRedis "github.com/go-redis/redis"
    "time"
)

const (
    TTL = 1 * time.Hour
    defaultFetchLimit = 1024
)

// uid -> list(TimelineItem)
func genKey(uid uint64) string {
    return fmt.Sprintf("TL:U:%d", uid)
}

func cached(uid uint64) (bool, error) {
    key := genKey(uid)
    exists, err := context.Redis.Exists(key).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, err
}

func Push(post *entity.Post) error {
    uid := post.Uid
    exist, err := cached(uid)
    if err != nil {
        return err
    }
    if !exist { // lazy push
        return nil
    }
    timelineItem := entity.MakeTimelineItem(post)
    key := genKey(uid)
    bytes, err := redis.Marshal(timelineItem)
    if err != nil {
        return err
    }
    _, err = context.Redis.LPush(key, bytes).Result()
    return err
}

func Remove(post *entity.Post) error {
    uid := post.Uid
    exists, err := cached(uid)
    if err != nil {
        return err
    }
    if !exists { // lazy push
        return nil
    }
    timelineItem := entity.MakeTimelineItem(post)
    key := genKey(uid)
    bytes, err := redis.Marshal(timelineItem)
    if err != nil {
        return err
    }
    _, err = context.Redis.LRem(key, -1, bytes).Result()
    return err
}

func getFromDB(uid uint64, start int32, length int32) ([]*entity.TimelineItem, error) {
    posts := make([]*entity.Post, length)
    err := context.DB.Table("posts").
        Where("uid = ? AND valid = 1", uid).
        Order("created_at DESC").
        Limit(length).Offset(start).
        Find(&posts).Error
    if err != nil {
        return nil, err
    }
    timeline := make([]*entity.TimelineItem, len(posts))
    for i, post := range posts {
        timeline[i] = entity.MakeTimelineItem(post)
    }
    return timeline, nil
}

func getFromCache(uid uint64, start int32, length int32) ([]*entity.TimelineItem, error) {
    key := genKey(uid)
    vals, err := context.Redis.LRange(key, int64(start), int64(start + length -1)).Result()
    if err != nil {
        return nil, err
    }
    timeline := make([]*entity.TimelineItem, len(vals))
    redis.MultiUnmarshalStr(vals, &timeline)
    return timeline, nil
}

/**
 * rebuild cache and return [start:start+length]
 * lazy rebuild: rebuild only if timeline is not cached
 * using check-lock-check pattern
 */
func Rebuild(uid uint64, start int32, length int32) ([]*entity.TimelineItem, error) {
    key := genKey(uid)

    // lock
    lock, err := RLock.Obtain(context.Redis, "lock:" + key, &RLock.Options{
        RetryCount: 3,
    })
    if err != nil {
        return nil, err
    }
    if lock == nil {
        return nil, errors.New("cannot obtain lock")
    }
    defer lock.Unlock()

    // check again
    cachedSize, err := context.Redis.LLen(key).Result()
    if err != nil {
        return nil, err
    }
    postCount, err := PostModel.GetUserPostCount(uid)
    if err != nil {
        return nil, err
    }
    // if cached needs or cached all, get from cache, do not rebuild
    if int64(start + length) <= cachedSize || cachedSize == int64(postCount) {
        return getFromCache(uid, start, length)
    }

    // rebuild
    timeline, err := rebuildInternal(uid, start + length)
    if int(start + length) > len(timeline) {
        return timeline[start:], nil
    } else {
        return timeline[start : start + length], nil
    }
}

/**
 * rebuild timeline cache, ignore whether is cached
 * invoker should obtain lock to ensure thread safe
 */
func rebuildInternal(uid uint64, limit int32) ([]*entity.TimelineItem, error) {
    if limit < defaultFetchLimit {
        limit = defaultFetchLimit
    }

    // get from db
    timeline, err := getFromDB(uid, 0, limit)
    if err != nil {
        return nil, err
    }

    // marshal
    vals := make([]interface{}, len(timeline))
    for i, item := range timeline {
        val, err := redis.Marshal(item)
        if err != nil {
            return nil, err
        }
        vals[i] = val
    }

    key := genKey(uid)

    _, err = context.Redis.Del(key).Result()
    if err != nil {
        return nil, err
    }

    _, err = context.Redis.RPush(key, vals...).Result()
    if err != nil {
        return nil, err
    }
    return timeline, nil
}

func Get(uid uint64, start int32, length int32) ([]*entity.TimelineItem, error) {
    key := genKey(uid)
    cachedSize, err := context.Redis.LLen(key).Result()
    if err != nil {
        return nil, err
    }

    // if size exceed cached range, rebuild cache
    if int64(start + length) > cachedSize {
        postCount, err := PostModel.GetUserPostCount(uid)
        if err != nil {
            return nil, err
        }
        // if postCount == 0 return empty array, don't rebuild to avoid cache breakdown
        if start > postCount || postCount == 0 { // no more posts
            return make([]*entity.TimelineItem, 0), nil
        }

        timeline, err := Rebuild(uid, start, length)
        if err != nil {
            return nil, err
        }
        return timeline, nil
    }

    // get from cache
    return getFromCache(uid, start, length)
}

// limit should not be greater than defaultFetchLimit
func MultiGet(uids []uint64, limit int32) (map[uint64][]*entity.TimelineItem, error) {
    if limit == 0 {
        return make(map[uint64][]*entity.TimelineItem), nil
    }

    cmdMap := make(map[uint64]*GoRedis.StringSliceCmd)
    pipe := context.Redis.Pipeline()
    for _, uid := range uids {
        key := genKey(uid)
        cmdMap[uid] = pipe.LRange(key, 0, int64(limit - 1))
    }

    _, err := pipe.Exec()
    if err != nil {
        return nil, err
    }

    timelineMap := make(map[uint64][]*entity.TimelineItem)
    missedUidSet := set.MakeUint64Set()
    for uid, cmd := range cmdMap {
        vals, err := cmd.Result()
        if err != nil {
            return nil, err
        }
        timeline := make([]*entity.TimelineItem, len(vals))
        redis.MultiUnmarshalStr(vals, &timeline)
        // if len(timeline) > 0 then len(timeline) >= defaultFetchLimit
        // ignore len(timeline) < limit
        if len(timeline) > 0 {
            timelineMap[uid] = timeline
        } else {
            missedUidSet.Add(uid)
        }
    }

    missedUids := missedUidSet.ToArray()
    for _, uid := range missedUids {
        timeline, err := Rebuild(uid, 0, limit)
        if err != nil {
            logger.Warn(fmt.Sprintf("err during rebuild, uid: %d, err: %s", uid, err.Error()))
            continue
        }
        timelineMap[uid] = timeline
    }
    return timelineMap, nil
}
