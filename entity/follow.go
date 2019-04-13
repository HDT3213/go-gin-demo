package entity

import "time"

type Follow struct {
    Uid uint64 `gorm:"primary_key;auto_increment:false"`
    FollowingUid uint64 `gorm:"primary_key;auto_increment:false"`
    CreatedAt time.Time
    Valid bool
}
