package cmap

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

// BenchmarkInsert benchmarks a normal set of insertion operations.
func BenchmarkInsert(t *testing.B) {
	csmap := New(10)
	insertElements := 20000

	start := time.Now()
	for i := 0; i < insertElements; i++ {
		csmap.Set(strconv.Itoa(i), ValTest{value: strconv.Itoa(i)})
	}
	duration := time.Since(start)

	fmt.Printf("Concurrent Map Duration to insert %v %v \n", insertElements,
		duration)
}

// BenchmarkInsertConcurrentInsert benchmarks a concurrent set of
// insertions.
func BenchmarkInsertConcurrentInsert(t *testing.B) {
	csmap := New(10)
	insertElements := 2000
	totalGoroutinesInserting := 10

	wg := &sync.WaitGroup{}
	startTime := time.Now()

	for i := 0; i < totalGoroutinesInserting; i++ {
		wg.Add(1)
		go func(multiple, totalElements int) {
			defer wg.Done()
			startInsert := i*insertElements
			endInsert := (i + 1)*insertElements

			for j := startInsert; j < endInsert; j++ {
				csmap.Set(strconv.Itoa(startInsert), ValTest{value: strconv.Itoa(i)})
			}
		}(i, insertElements)
	}
	wg.Wait()
	endTime := time.Since(startTime)
	fmt.Printf("10 Goroutines Map Duration to insert %v %v \n",
		insertElements, endTime)
}


// BenchmarkGetKeys benchmarks getting the all the keys.
func BenchmarkGetKeys(t *testing.B) {
	csmap := New(10)
	insertElements := 200000

	for i := 0; i < insertElements; i++ {
		csmap.Set(strconv.Itoa(i), ValTest{value: strconv.Itoa(i)})
	}

	start := time.Now()
	csmap.Keys()
	duration := time.Since(start)
	fmt.Printf("Get Keys Time %v \n", duration)
}