package cmap

import "sync"

// ConcurrentSafeMap is a thread safe map type (string to interface type).
// We map this into a bunch of shards to reduce lock bottlenecks.
type ConcurrentSafeMap struct {
	MapParts   []*ConcurrentMapPart
	shardCount int
}

// ConcurrentMapPart is thread safe string to interface type.
type ConcurrentMapPart struct {
	items map[string]interface{}
	sync.RWMutex
}

// Tuple is contains key and associated value.
type Tuple struct {
	key string
	val interface{}
}

// UpsertCallbackFn is a callback function which looks at whether a value already exists
// in the map, value already in the map, and the new value to be inserted.
type UpsertCallbackFn func(valueExists bool, valueInMap interface{}, newValue interface{}) interface{}

// RemoveCallbackFn is a callback function which if it is true, element will be removed.
type RemoveCallbackFn func(key string, v interface{}, valueExists bool) bool

