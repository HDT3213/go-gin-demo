package model

import (
    "time"
    "go-close/utils"
    "strconv"
    "math/rand"
)

type User struct {
    ID        uint64     `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index"`
    Username  string     `json:"username"`
    Password  string
}

func CreateUser(user *User) bool {
    user.ID = utils.Hash64(user.Username + strconv.Itoa(int(utils.Now())) + strconv.Itoa(rand.Int()))
    if !db.NewRecord(*user) {
        if err := db.Create(user).Error; err != nil {
            panic(err)
        }
        return true
    }
    return false
}

func GetUser(uid uint64) *User {
    var user User
    db.First(&user, uid)
    return &user
}

func GetUserByName(username string) (*User, error) {
    var user User
    err := db.Where("username = ?", username).First(&user).Error
    return &user, err
}

func AllUsers() []User {
    var users []User
    err := db.Find(&users).Error
    if err != nil {
        panic(err)
    }
    return users
}