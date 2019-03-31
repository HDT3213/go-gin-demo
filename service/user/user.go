package user

import (
    "github.com/go-gin-demo/entity"
    "strconv"
    UserModel "github.com/go-gin-demo/model/user"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
)

func RenderUsers(users []*entity.User) (entities []*entity.UserEntity) {
    entities = make([]*entity.UserEntity, len(users))
    for i, user := range users {
        if user == nil {
            continue
        }
        entities[i] = &entity.UserEntity{ID:user.ID, IDStr:strconv.FormatUint(user.ID, 10), Username:user.Username}
    }
    return
}

func RenderUsersById(uids []uint64) ([]*entity.UserEntity, error) {
    users, err := UserModel.MultiGet(uids)
    if err != nil {
        return nil, err
    }
    return RenderUsers(users), nil
}

func RenderUser(user *entity.User) *entity.UserEntity {
    users := []*entity.User{user}
    entities := RenderUsers(users)
    return entities[0]
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
    return RenderUser(user), nil
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
    return RenderUser(user), user.ID, nil
}

func GetUser(uid uint64) (*entity.UserEntity, error) {
    user, err := UserModel.Get(uid)
    if err != nil {
        panic(err)
    }
    if user == nil {
        return nil, BizError.NotFound(fmt.Sprintf("user not found: %d", uid))
    }
    return RenderUser(user), nil
}