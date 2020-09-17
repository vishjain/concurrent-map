package cmap

import (
	"sync"
)

// Create a new concurrent map. Specify the shard count.
func New(shardCount int) *ConcurrentSafeMap {
	mapParts := make([]*ConcurrentMapPart, shardCount)
	for i := 0; i < shardCount; i++ {
		mapParts[i] = &ConcurrentMapPart{
			items: make(map[string]interface{}),
		}
	}
	return &ConcurrentSafeMap{
		MapParts:   mapParts,
		shardCount: shardCount,
	}
}

// GetShard gives the shard for a certain key.
func (csmap *ConcurrentSafeMap) GetShardPartition(key string) *ConcurrentMapPart {
	return csmap.MapParts[uint(fnv32(key))%uint(csmap.shardCount)]
}

// MSet takes all data in the map and places into the map (after getting the shard).
func (csmap *ConcurrentSafeMap) MSet(data map[string]interface{}) {
	for key, value := range data {
		shard := csmap.GetShardPartition(key)
		shard.Lock()
		shard.items[key] = value
		shard.Unlock()
	}
}

// Set takes a key and a value and inserts inside map.
func (csmap *ConcurrentSafeMap) Set(key string, value interface{}) {
	shard := csmap.GetShardPartition(key)
	shard.Lock()
	shard.items[key] = value
	shard.Unlock()
}

// Insert or Update according to the callback function passed in. Returns value in the shared map.
func (csmap *ConcurrentSafeMap) Upsert(key string, value interface{}, cb UpsertCallbackFn) (res interface{}) {
	// Get shard
	shard := csmap.GetShardPartition(key)
	shard.Lock()
	// check if value is in the shard
	oldVal, ok := shard.items[key]
	valToInsert := cb(ok, oldVal, value)
	// figure out if needs to replace and replace
	shard.items[key] = valToInsert
	shard.Unlock()
	return valToInsert
}

// SetIfAbsent set's the given value if there's not specified key (if no value associated with it).
// Returns a new value was set or not.
func (csmap *ConcurrentSafeMap) SetIfAbsent(key string, value interface{}) bool {
	shard := csmap.GetShardPartition(key)
	shard.Lock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	shard.Unlock()
	return !ok
}

// Get function will return the value in the map.
func (csmap *ConcurrentSafeMap) Get(key string) (interface{}, bool) {
	shard := csmap.GetShardPartition(key)
	shard.RLock()
	val, ok := shard.items[key]
	if !ok {
		return nil, ok
	}
	shard.RUnlock()
	return val, ok
}

// Get the total number of items (protected with a read lock).
func (csmap *ConcurrentSafeMap) Count() int {
	count := 0
	for i := 0; i < csmap.shardCount; i++ {
		shard := csmap.MapParts[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Returns whether the map has a key or not.
func (csmap *ConcurrentSafeMap) Has(key string) bool {
	shard := csmap.GetShardPartition(key)
	shard.RLock()
	_, ok := shard.items[key]
	shard.RUnlock()
	return ok
}

// Remove just deletes the given key.
func (csmap *ConcurrentSafeMap) Remove(key string) {
	shard := csmap.GetShardPartition(key)
	shard.Lock()
	delete(shard.items, key)
	shard.Unlock()
}

// IsEmpty returns whether the map
func (csmap *ConcurrentSafeMap) IsEmpty() bool {
	return csmap.Count() == 0
}

// RemoveCb will remove a key if its in the map and the cb callback function
// says to do so.
func (csmap *ConcurrentSafeMap) RemoveCb(key string, cb RemoveCallbackFn) bool {
	shard := csmap.GetShardPartition(key)
	shard.Lock()
	v, ok := shard.items[key]
	removeBool := cb(key, v, ok)
	if removeBool && ok {
		delete(shard.items, key)
	}
	shard.Unlock()
	return removeBool
}

// Keys function returns all the keys from the concurrent map.
func (csmap *ConcurrentSafeMap) Keys() []string {
	keyChan := make(chan string)

	// First goroutine to send the keys.
	WriteAndReadWaitGroup := sync.WaitGroup{}
	WriteAndReadWaitGroup.Add(1)
	go func() {
		defer func() {
			close(keyChan)
			WriteAndReadWaitGroup.Done()
		}()

		for _, shard := range csmap.MapParts {
			for key := range shard.items {
				keyChan <- key
			}
		}
	}()

	// Second goroutine listening to take the keys.
	WriteAndReadWaitGroup.Add(1)
	keys := make([]string, 0, csmap.Count())
	go func() {
		defer func() {
			WriteAndReadWaitGroup.Done()
		}()
		for k := range keyChan {
			keys = append(keys, k)
		}
	}()
	WriteAndReadWaitGroup.Wait()
	return keys
}

func (csmap *ConcurrentSafeMap) collectShardValues(
	shardWorkersWaitGroup *sync.WaitGroup,
	mapShard *ConcurrentMapPart, valChan chan interface{}) {

	defer shardWorkersWaitGroup.Done()
	for _, value := range mapShard.items {
		valChan <- value
	}
}

// runShardSearchValueWorkers is a helper function that starts new goroutine
// for each shard and searches for values in the map.
func (csmap *ConcurrentSafeMap) runShardSearchValueWorkers(valChan chan interface{}) {

	shardWorkersWaitGroup := &sync.WaitGroup{}
	shardWorkersWaitGroup.Add(csmap.shardCount)
	for _, shard := range csmap.MapParts {
		go csmap.collectShardValues(shardWorkersWaitGroup, shard, valChan)
	}
	shardWorkersWaitGroup.Wait()
}

// Values implements parallelized mechanism to retrieve all the values in
// the whole map.
func (csmap *ConcurrentSafeMap) Values() []interface{} {
	values := make([]interface{}, csmap.Count())
	valChan := make(chan interface{})
	wgWaitGroup := &sync.WaitGroup{}

	wgWaitGroup.Add(1)
	go func() {
		defer func() {
			close(valChan)
			wgWaitGroup.Done()
		}()

		csmap.runShardSearchValueWorkers(valChan)
	}()

	wgWaitGroup.Add(1)
	go func() {
		defer wgWaitGroup.Done()
		idx := 0
		for val := range valChan {
			values[idx] = val
			idx += 1
		}
	}()
	wgWaitGroup.Wait()
	return values
}
