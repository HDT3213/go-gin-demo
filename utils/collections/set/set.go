package set

type Set struct {
    hash map[interface{}]bool
}

func Make() *Set {
    return &Set{hash: make(map[interface{}]bool)}
}

func (set *Set)Add(e interface{}) {
    set.hash[e] = true
}

func (set *Set)Remove(e interface{}) {
    delete(set.hash, e)
}

func (set *Set)Len() int {
    return len(set.hash)
}

func (set *Set)Has(e interface{}) bool {
    _, has := set.hash[e]
    return has
}

func (set *Set)ToArray() []interface{} {
    arr := make([]interface{}, set.Len())
    i := 0
    for k := range set.hash {
       arr[i] = k
       i++
    }
    return arr
}