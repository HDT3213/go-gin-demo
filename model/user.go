package model

import (
    "go-close/utils"
    "go-close/entity"
    "strconv"
    "math/rand"
    "encoding/json"
)

func CreateUser(user *entity.User) error {
    user.ID = utils.Hash64(user.Username + strconv.Itoa(int(utils.Now())) + strconv.Itoa(rand.Int()))
    if !db.NewRecord(*user) {
        if err := db.Create(user).Error; err != nil {
            return err
        }
        return nil
    }
    return nil
}

func GetUser(uid uint64) (*entity.User, error) {
    var user entity.User
    err := db.First(&user, uid).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return &user, err
}

func GetUserByName(username string) (*entity.User, error) {
    var user entity.User
    err := db.Where("username = ?", username).First(&user).Error
    if err != nil && err.Error() == "record not found" {
        return nil, nil
    }
    return &user, err
}

func Test() {
    user := &entity.User{ID:1, Username:"fuck"}
    data, err := json.Marshal(user)
    if err != nil {
        panic(err)
    }
    Client.Set("U:1", data, 0)
    val, err := Client.Get("U:1").Result()
    if err != nil {
        panic(err)
    }
    user2 := new(entity.User)
    err = json.Unmarshal([]byte(val), user2)
    if err != nil {
        panic(err)
    }
    println(user2)
}