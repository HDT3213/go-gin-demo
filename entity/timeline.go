package entity

import "time"

type TimelineItem struct {
    PostId uint64
    Uid uint64
    Time time.Time
}

func MakeTimelineItem(post *Post) *TimelineItem {
    return &TimelineItem{
        PostId: post.ID,
        Uid: post.Uid,
        Time: post.CreatedAt,
    }
}
