package set

import (
    "reflect"
)

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

func (set *Set)ToArray(out interface{}) { // out should be a pointer of slice
    slice := reflect.ValueOf(out).Elem()
    if slice.Kind() != reflect.Slice {
        panic("out is not slice")
    }
    if slice.Cap() < set.Len() {
        panic("out is too small")
    }
    slice.SetLen(set.Len())
    i := 0
    for k := range set.hash {
        slice.Index(i).Set(reflect.ValueOf(k))
        i++
    }
}

