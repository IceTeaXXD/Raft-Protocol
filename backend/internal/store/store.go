package store

import (
	"fmt"
	"strconv"
	"sync"
)

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

func Strln(key string) string {
    value := Get(key)
    length := len(value)
    fmt.Println(length)
    fmt.Println(strconv.Itoa(length))
    return strconv.Itoa(length)
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
