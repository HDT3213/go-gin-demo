package entity

import (
    "time"
)

type User struct {
    ID        uint64     `gorm:"primary_key;auto_increment:false"`
    CreatedAt time.Time
    UpdatedAt time.Time
    Valid     bool
    Username  string
    Password  string
}

type UserEntity struct {
    ID       uint64 `json:"-"`
    IDStr    string `json:"id"`
    Username string `json:"username"`
    PostCount int32 `json:"postCount"`
    FollowingCount int32 `json:"followingCount"`
    FollowerCount int32 `json:"followerCount"`
    Following bool `json:"following"`
}
