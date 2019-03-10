package service

import (
    "go-close/model"
    "go-close/entity"
    "strconv"
)

func RenderUsers(users []*model.User) (entities []*entity.UserEntity) {
    entities = make([]*entity.UserEntity, len(users))
    for i, user := range users {
        entity := &entity.UserEntity{Id:strconv.FormatUint(user.ID, 10), Username:user.Username}
        entities[i] = entity
    }
    return
}

func RenderUser(user *model.User) *entity.UserEntity {
    users := make([]*model.User, 1)
    users[0] = user
    entities := RenderUsers(users)
    return entities[0]
}