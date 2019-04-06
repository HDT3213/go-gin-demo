package user

import (
    "github.com/go-gin-demo/entity"
    "strconv"
    UserModel "github.com/go-gin-demo/model/user"
    PostModel "github.com/go-gin-demo/model/post"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
    "github.com/go-gin-demo/utils/collections/set"
)

func RenderUsers(users []*entity.User) ([]*entity.UserEntity, error) {
    uidSet := set.MakeUint64Set()
    for _, user := range users {
        uidSet.Add(user.ID)
    }
    uids := uidSet.ToArray()
    postCountMap, err := PostModel.GetUserPostCountMap(uids)
    if err != nil {
        return nil, err
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
        entities[i] = &entity.UserEntity{
            ID:user.ID,
            IDStr:strconv.FormatUint(user.ID, 10),
            Username:user.Username,
            PostCount: postCount,
        }
    }
    return entities, nil
}

func RenderUsersById(uids []uint64) ([]*entity.UserEntity, error) {
    users, err := UserModel.MultiGet(uids)
    if err != nil {
        return nil, err
    }
    return RenderUsers(users)
}

func RenderUser(user *entity.User) (*entity.UserEntity, error) {
    users := []*entity.User{user}
    entities, err := RenderUsers(users)
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
    return RenderUser(user)
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
    userEntity, err := RenderUser(user)
    return userEntity, user.ID, nil
}

func GetUser(uid uint64) (*entity.UserEntity, error) {
    user, err := UserModel.Get(uid)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, BizError.NotFound(fmt.Sprintf("user not found: %d", uid))
    }
    return RenderUser(user)
}