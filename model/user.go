package model

import (
    "time"
    "../utils"
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
    if !db.NewRecord(user) {
        db.Create(user)
        return true
    }
    return false
}

func GetUser(uid uint64) *User {
    var user User
    db.Where("id = ?", uid).First(&user)
    return &user
}

func AllUsers() []User {
    var users []User
    db.Find(&users)
    return users
}
