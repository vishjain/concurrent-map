## Golang Thread Safe Map

Golang's in-built map go has poor support for concurrent operations.
Moreover, many advanced operations are not thread-safe and/or
performant. This version allows for concurrent writes and reads, while
minimizing lock contention. Ideal if you want a quick in-memory datastore (cache, 
key-value store, etc). Also includes custom & useful operations. Basic unit-tests
and benchmarks are given. 
Some operations supported include: 

// MSet takes all data in Golang built-in map and places into the thread-safe map.
 
// Set takes a key and value and inserts into the thread-safe map. 

// Upsert - Insert or Update according to the callback function passed in. 
Returns value in the shared map.

// SetIfAbsent set's a given value if key is not in the map. Returns whether new value was set.

// Get returns value in the map for a particular key. 

// Count returns total number of items in thread-safe map. 

// Has returns whether the map has a key or not. 

// Remove deletes the given key. 

// IsEmpty returns whether thread-safe map is empty or not.

// RemoveCb takes a key and removes the value associated with that key
according to a callback.  

// Keys returns all the keys in the thread-safe map. 

// Values returns all the values in the thread-safe map. 

// Iterator returns an iterator (in the form of a buffered Golang channel) so the user can retrieve all
keys and values. 

// Merge merges two thread-safe maps in a parallelized fashion. 
