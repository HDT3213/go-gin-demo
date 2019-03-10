package entity

import "time"

type User struct {
    ID        uint64     `gorm:"primary_key"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time `sql:"index"`
    Username  string     `json:"username"`
    Password  string
}

type UserEntity struct {
    Id string `json:"id"`
    Username string `json:"username"`
}
