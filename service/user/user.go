package user

import (
    "github.com/go-gin-demo/entity"
    "strconv"
    "github.com/go-gin-demo/model"
    BizError "github.com/go-gin-demo/errors"
    "fmt"
)

func RenderUsers(users []*entity.User) (entities []*entity.UserEntity) {
    entities = make([]*entity.UserEntity, len(users))
    for i, user := range users {
        entity := &entity.UserEntity{Id:strconv.FormatUint(user.ID, 10), Username:user.Username}
        entities[i] = entity
    }
    return
}

func RenderUser(user *entity.User) *entity.UserEntity {
    users := make([]*entity.User, 1)
    users[0] = user
    entities := RenderUsers(users)
    return entities[0]
}

func Register(username string, password string) (*entity.UserEntity, error) {
    user := &entity.User{Username:username, Password:password, Valid:true}
    err := model.CreateUser(user)
    if err != nil {
        panic(err)
    }
    return RenderUser(user), nil
}

func Login(username string, password string) (*entity.UserEntity, uint64, error) {
    user, err := model.GetUserByName(username)
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
    user, err := model.GetUser(uid)
    if err != nil {
        panic(err)
    }
    if user == nil {
        return nil, BizError.NotFound(fmt.Sprintf("user not found: %d", uid))
    }
    return RenderUser(user), nil
}