package post

import (
    "github.com/go-gin-demo/entity"
    UserTimelineModel "github.com/go-gin-demo/model/timeline/user"
    UserModel "github.com/go-gin-demo/model/user"
    PostModel "github.com/go-gin-demo/model/post"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
)

func GetUserTimeline(uid uint64, start int32, length int32) ([]*entity.PostEntity, int32, error) {
    user, err := UserModel.Get(uid)
    if err != nil {
        return nil, -1, err
    }
    if user == nil {
        return nil, -1, BizError.NotFound(fmt.Sprintf("user not found: %d", uid))
    }
    timeline, err := UserTimelineModel.Get(uid, start, length)
    if err != nil {
        return nil, -1, err
    }
    pids := make([]uint64, len(timeline))
    for i, item := range timeline {
        pids[i] = item.PostId
    }

    totalCount, err := PostModel.GetUserPostCount(uid)
    if err != nil {
        return nil, -1, err
    }

    posts, err := PostModel.MultiGet(pids)
    if err != nil {
        return nil, -1, err
    }
    entities, err := RenderPosts(posts)
    if err != nil {
        return nil, -1, err
    }
    return entities, totalCount, nil
}