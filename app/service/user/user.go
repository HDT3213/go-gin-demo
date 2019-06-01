package user

import (
    "github.com/go-gin-demo/app/entity"
    "strconv"
    UserModel "github.com/go-gin-demo/app/model/user"
    PostModel "github.com/go-gin-demo/app/model/post"
    FollowModel "github.com/go-gin-demo/app/model/follow"
    FollowTimelineModel "github.com/go-gin-demo/app/model/timeline/follow"
    BizError "github.com/go-gin-demo/lib/errors"
    "fmt"
    "github.com/go-gin-demo/lib/collections/set"
    "github.com/go-gin-demo/lib/canal"
    "github.com/mitchellh/mapstructure"
    "github.com/go-gin-demo/app/context/context"
)

func RenderUsers(currentUid uint64, users []*entity.User) ([]*entity.UserEntity, error) {
    if len(users) == 0 {
        return []*entity.UserEntity{}, nil
    }
    uidSet := set.MakeUint64Set()
    for _, user := range users {
        uidSet.Add(user.ID)
    }
    uids := uidSet.ToArray()

    postCountMap, err := PostModel.GetUserPostCountMap(uids)
    if err != nil {
        return nil, err
    }

    followingCountMap, err := FollowModel.GetFollowingCountMap(uids)
    if err != nil {
        return nil, err
    }

    followerCountMap, err := FollowModel.GetFollowerCountMap(uids)
    if err != nil {
        return nil, err
    }

    followingSet := set.MakeUint64Set()
    if currentUid != 0 {
        followingUids, err := FollowModel.GetFollowingIn(currentUid, uids)
        if err != nil {
            return nil, err
        }
        for _, followingUid := range followingUids {
            followingSet.Add(followingUid)
        }
    }

    entities := make([]*entity.UserEntity, len(users))
    for i, user := range users {
        uid := user.ID
        if user == nil {
            continue
        }
        postCount, ok := postCountMap[uid]
        if !ok {
            postCount = 0
        }
        followingCount, ok := followingCountMap[uid]
        if !ok {
            followingCount = 0
        }
        followerCount, ok := followerCountMap[uid]
        if !ok {
            followerCount = 0
        }

        entities[i] = &entity.UserEntity{
            ID:uid,
            IDStr:strconv.FormatUint(uid, 10),
            Username:user.Username,
            PostCount: postCount,
            FollowingCount: followingCount,
            FollowerCount: followerCount,
            Following: followingSet.Has(uid),
        }
    }
    return entities, nil
}

func RenderUsersById(currentUid uint64, uids []uint64) ([]*entity.UserEntity, error) {
    users, err := UserModel.MultiGet(uids)
    if err != nil {
        return nil, err
    }
    return RenderUsers(currentUid, users)
}

func RenderUser(currentUid uint64, user *entity.User) (*entity.UserEntity, error) {
    users := []*entity.User{user}
    entities, err := RenderUsers(currentUid, users)
    if err != nil {
        return nil, err
    }
    return entities[0], nil
}

func Register(username string, password string) (*entity.UserEntity, error) {
    exists, err := UserModel.GetByName(username)
    if err != nil { // error during access database
        return nil, err
    }
    if exists != nil {
        return nil, BizError.InvalidForm("username has been used")
    }
    user := &entity.User{Username:username, Password:password, Valid:true}
    err = UserModel.Create(user)
    if err != nil {
        return nil, err
    }
    return RenderUser(user.ID, user)
}

func AfterCreateUser(former canal.Row, row canal.Row) error {
    var user entity.User
    err := mapstructure.Decode(row, &user)
    if err != nil {
        return err
    }
    if !user.Valid {
        return nil
    }
    return UserModel.AfterCreate(&user)
}

func AfterUpdateUser(former canal.Row, row canal.Row) error {
    return nil
}


func Login(username string, password string) (*entity.UserEntity, uint64, error) {
    user, err := UserModel.GetByName(username)
    if err != nil {
        panic(err)
    }
    if user == nil {
        return nil, 0, BizError.NotFound("user not found: " + username)
    }
    if user.Password != password {
        return nil, 0, BizError.InvalidForm("user password mismatch")
    }
    userEntity, err := RenderUser(user.ID, user)
    return userEntity, user.ID, nil
}

func GetUser(currentUid uint64, uid uint64) (*entity.UserEntity, error) {
    user, err := UserModel.Get(uid)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, BizError.NotFound(fmt.Sprintf("user not found: %d", uid))
    }

    return RenderUser(currentUid, user)
}

func Follow(currentUid uint64, followingUid uint64) error {
    if currentUid == followingUid {
        return BizError.InvalidForm("cannot follow yourself")
    }

    // check user exists
    user, err := UserModel.Get(followingUid)
    if err != nil {
        return err
    }
    if user == nil {
        return BizError.NotFound(fmt.Sprintf("user not found: %d", followingUid))
    }

    // check has followed
    isFollowing, err := FollowModel.IsFollowing(currentUid, followingUid)
    if err != nil {
        return err
    }
    if isFollowing {
        return nil
    }

    err = FollowModel.Create(&entity.Follow{
        Uid: currentUid,
        FollowingUid: followingUid,
        Valid: true,
    })
    if err != nil {
        return err
    }
    if !context.EnableCanal() {
        FollowTimelineModel.Del(currentUid)
    }
    return nil
}

func AfterCreateFollow(former canal.Row, row canal.Row) error {
    var follow entity.Follow
    err := mapstructure.Decode(row, &follow)
    if err != nil {
        return err
    }
    if !follow.Valid {
        return nil
    }

    err = FollowModel.AfterCreate(&follow)
    if err != nil {
        return err
    }

    err = FollowTimelineModel.Del(follow.Uid)
    return err
}

func UnFollow(currentUid uint64, followingUid uint64) error {
    if currentUid == followingUid {
        return BizError.InvalidForm("cannot follow yourself")
    }

    // check user exists
    user, err := UserModel.Get(followingUid)
    if err != nil {
        return err
    }
    if user == nil {
        return BizError.NotFound(fmt.Sprintf("user not found: %d", followingUid))
    }

    // check has followed
    isFollowing, err := FollowModel.IsFollowing(currentUid, followingUid)
    if err != nil {
        return err
    }
    if !isFollowing {
        return nil
    }

    err = FollowModel.Delete(currentUid, followingUid)
    if err != nil {
        return err
    }
    if !context.EnableCanal() {
        FollowTimelineModel.Del(currentUid)
    }
    return nil
}

func AfterUpdateFollow(former canal.Row, row canal.Row) error {
    var follow entity.Follow
    err := mapstructure.Decode(row, &follow)
    if err != nil {
        return err
    }
    if !follow.Valid {
        err = FollowModel.AfterDelete(&follow)
        if err != nil {
            return err
        }

        err = FollowTimelineModel.Del(follow.Uid)
        return err
    }
    return nil
}

func GetFollowings(currentUid uint64, uid uint64, start int32, length int32) ([]*entity.UserEntity, int32, error) {
    follows, err := FollowModel.GetFollowings(uid, start, length)
    if err != nil {
        return  nil, 0, err
    }

    followingUids := make([]uint64, len(follows))
    for i, follow := range follows {
        followingUids[i] = follow.FollowingUid
    }

    users, err := UserModel.MultiGet(followingUids)
    if err != nil {
        return nil, 0, err
    }

    followingCount, err := FollowModel.GetFollowingCount(uid)
    if err != nil {
        return nil, 0, err
    }

    entities, err := RenderUsers(currentUid, users)
    if err != nil {
        return nil, 0, err
    }

    return entities, followingCount, nil
}

func GetFollowers(currentUid uint64, uid uint64, start int32, length int32) ([]*entity.UserEntity, int32, error) {
    follows, err := FollowModel.GetFollowers(uid, start, length)
    if err != nil {
        return  nil, 0, err
    }

    followerUids := make([]uint64, len(follows))
    for i, follow := range follows {
        followerUids[i] = follow.Uid
    }

    users, err := UserModel.MultiGet(followerUids)
    if err != nil {
        return nil, 0, err
    }

    followerCount, err := FollowModel.GetFollowerCount(uid)
    if err != nil {
        return nil, 0, err
    }

    entities, err := RenderUsers(currentUid, users)
    if err != nil {
        return nil, 0, err
    }

    return entities, followerCount, nil
}