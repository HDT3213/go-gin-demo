package collections

type Pair interface {
    Key() interface{}
    Value() interface{}
}

type IdCountPair struct {
    Id uint64
    Num int32
}

func (pair *IdCountPair)Key() uint64 {
    return pair.Id
}

func (pair *IdCountPair)Value() int32 {
    return pair.Num
}

func ToCountMap(pairs []*IdCountPair) map[uint64]int32 {
    countMap := make(map[uint64]int32)
    for _, pair := range pairs {
        countMap[pair.Id] = pair.Num
    }
    return countMap
}