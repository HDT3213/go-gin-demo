package set

// avoid reflect for better performance
type Uint64Set struct {
    hash map[uint64]bool
}

func MakeUint64Set() *Uint64Set {
    return &Uint64Set{hash: make(map[uint64]bool)}
}

func (set *Uint64Set)Add(e uint64) {
    set.hash[e] = true
}

func (set *Uint64Set)Remove(e uint64) {
    delete(set.hash, e)
}

func (set *Uint64Set)Len() int {
    return len(set.hash)
}

func (set *Uint64Set)Has(e uint64) bool {
    _, has := set.hash[e]
    return has
}

func (set *Uint64Set)ToArray() []uint64 {
    arr := make([]uint64, set.Len())
    i := 0
    for k := range set.hash {
        arr[i] = k
        i++
    }
    return arr
}

func (set *Uint64Set)ToInterfaceArray() []interface{} {
    arr := make([]interface{}, set.Len())
    i := 0
    for k := range set.hash {
        arr[i] = k
        i++
    }
    return arr
}