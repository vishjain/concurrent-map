package cmap

import (
	"sync"
)


// getAllChannels is a helper function to put all data from the concurrent map
// into a channel.
func (csmap *ConcurrentSafeMap) getAllChannels() chan Tuple {
	shardChans := make(chan Tuple, csmap.Count())
	wg := sync.WaitGroup{}
	wg.Add(csmap.shardCount)
	for idx, shard := range csmap.MapParts {

		go func(index int, shardPart *ConcurrentMapPart) {

			defer func() {
				wg.Done()
			}()

			shardPart.RLock()
			for key, val := range shardPart.items {
				shardChans <- 	Tuple {
					key: key,
					val: val,
				}
			}
			shardPart.RUnlock()
		}(idx, shard)
	}
	wg.Wait()
	close(shardChans)
	return shardChans
}

// Returns a channel, you can read from that channel.
func (csmap *ConcurrentSafeMap) Iterator() <-chan Tuple {
	return csmap.getAllChannels()
}
