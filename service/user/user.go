package user

import (
    "github.com/go-gin-demo/entity"
    "strconv"
    UserModel "github.com/go-gin-demo/model/user"
    PostModel "github.com/go-gin-demo/model/post"
    FollowModel "github.com/go-gin-demo/model/follow"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
    "github.com/go-gin-demo/utils/collections/set"
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
    return err
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