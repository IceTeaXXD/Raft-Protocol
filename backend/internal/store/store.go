package store

import "sync"

var store sync.Map

func Get(key string) string {
    value, ok := store.Load(key)
    if !ok {
        return ""
    }
    return value.(string)
}

func Set(key, value string) {
    store.Store(key, value)
}

func Strln(key string) int {
    value := Get(key)
    return len(value)
}

func Del(key string) string {
    value, ok := store.LoadAndDelete(key)
    if !ok {
        return ""
    }
    return value.(string)
}

func Append(key, value string) {
    existingValue := Get(key)
    newValue := existingValue + value
    Set(key, newValue)
}
