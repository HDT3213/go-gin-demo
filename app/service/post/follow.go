package post

import (
    "github.com/go-gin-demo/app/entity"
    FollowTimelineModel "github.com/go-gin-demo/app/model/timeline/follow"
    PostModel "github.com/go-gin-demo/app/model/post"
)

func GetFollowingTimeline(currentUid uint64, start int32, length int32) ([]*entity.PostEntity, int32, error) {
    timeline, totalCount, err := FollowTimelineModel.Get(currentUid, start, length)
    if err != nil {
        return nil, -1, err
    }
    pids := make([]uint64, len(timeline))
    for i, item := range timeline {
        pids[i] = item.PostId
    }

    posts, err := PostModel.MultiGet(pids)
    if err != nil {
        return nil, -1, err
    }
    entities, err := RenderPosts(currentUid, posts)
    if err != nil {
        return nil, -1, err
    }
    return entities, totalCount, nil
}
