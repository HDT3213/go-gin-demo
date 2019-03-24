package entity

import (
    "time"
)

type Post struct {
    ID        uint64     `gorm:"primary_key"`
    Uid       uint64
    CreatedAt time.Time
    UpdatedAt time.Time
    Valid     bool
    Text      string    `gorm:"type:text"`
}

type PostEntity struct {
    Id string `json:"id"`
    Text string `json:"text"`
    User *UserEntity `json:"user"`
    CreatedAt time.Time `json:"created_at"`
}
