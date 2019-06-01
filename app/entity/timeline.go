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

type TimelineHeap []*TimelineItem

func (h TimelineHeap) Len() int  {
    return len(h)
}

func (h TimelineHeap) Less(i, j int) bool {
    return h[i].Time.After(h[j].Time)
}

func (h TimelineHeap) Swap(i, j int) {
    h[i], h[j] = h[j], h[i]
}

func (h *TimelineHeap) Push(x interface{}) {
    *h = append(*h, x.(*TimelineItem))
}

func (h *TimelineHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[0 : n-1]
    return x
}
