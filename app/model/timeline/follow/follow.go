package follow

import (
    "fmt"
    "github.com/go-gin-demo/app/entity"
    FollowModel "github.com/go-gin-demo/app/model/follow"
    UserTimelineModel "github.com/go-gin-demo/app/model/timeline/user"
    RLock "github.com/bsm/redis-lock"
    "errors"
    "time"
    "container/heap"
    "github.com/go-gin-demo/app/context/context"
    "github.com/go-gin-demo/lib/cache/redis"
)

const (
    fetchLimit = 128
    cacheLimit = 1024
    cacheTTL = 15 * time.Minute
)

// uid -> list(TimelineItem)
func genKey(uid uint64) string {
    return fmt.Sprintf("TL:Fo:%d", uid)
}

func cached(uid uint64) (bool, error) {
    key := genKey(uid)
    exists, err := context.Redis.Exists(key).Result()
    if err != nil {
        return false, err
    }
    return exists > 0, err
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
    exists, err := cached(uid)
    if err != nil {
        return nil, err
    }
    if exists {
        return getFromCache(uid, start, length)
    }

    // rebuild
    timeline, err := rebuildInternal(uid)
    if int(start + length) > len(timeline) {
        return timeline[start:], nil
    } else {
        return timeline[start : start + length], nil
    }
}

func rebuildInternal(uid uint64) ([]*entity.TimelineItem, error) {
    // get following
    followings, err := FollowModel.GetAllFollowings(uid)
    if err != nil {
        return nil, err
    }
    // get following posts
    timelineMap, err := UserTimelineModel.MultiGet(followings, fetchLimit)

    // merge and sort
    timelineHeap := &entity.TimelineHeap{}
    heap.Init(timelineHeap)
    for _, timeline := range timelineMap {
        for _, item := range timeline {
            heap.Push(timelineHeap, item)
        }
    }

    // trim
    size := cacheLimit
    if size > timelineHeap.Len() {
        size = timelineHeap.Len()
    }

    followingTimeline := make([]*entity.TimelineItem, size)
    for i := 0; i < size; i++ {
        item, _ := heap.Pop(timelineHeap).(*entity.TimelineItem)
        followingTimeline[i] = item
    }

    // set cache
    if size > 0 {
        key := genKey(uid)
        vals := make([]interface{}, len(followingTimeline))
        for i, item := range followingTimeline {
            val, err := redis.Marshal(item)
            if err != nil {
                return nil, err
            }
            vals[i] = val
        }

        _, err = context.Redis.Del(key).Result()
        if err != nil {
            return nil, err
        }

        _, err = context.Redis.RPush(key, vals...).Result()
        if err != nil {
            return nil, err
        }

        _, err = context.Redis.Expire(key, cacheTTL).Result()
        if err != nil {
            return nil, err
        }
    }
    return followingTimeline, nil
}


// invoker should ensure timeline cached
func getTotalCount(uid uint64) (int32, error) {
    key := genKey(uid)
    count, err := context.Redis.LLen(key).Result()
    if err != nil {
        return 0, err
    }
    return int32(count), nil
}

func Get(uid uint64, start int32, length int32) ([]*entity.TimelineItem, int32, error) {
    exists, err := cached(uid)
    if err != nil {
        return nil, 0, err
    }

    var timeline []*entity.TimelineItem
    if exists {
        timeline, err = getFromCache(uid, start, length)
    } else {
        timeline, err = Rebuild(uid, start, length)
    }
    if err != nil {
        return nil, 0, err
    }

    totalCount, err := getTotalCount(uid)
    if err != nil {
        return nil, 0, err
    }

    return timeline, totalCount, nil
}

func Del(uid uint64) error {
    key := genKey(uid)
    _, err := context.Redis.Del(key).Result()
    return err
}