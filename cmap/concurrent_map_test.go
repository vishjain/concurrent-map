package cmap

import (
	"github.com/stretchr/testify/require"
	"sort"
	"strconv"
	"sync"
	"testing"
)

// ValTest will be used as a basic
//
type ValTest struct {
	value string
}

// TestBasicSetAndGet of the function.
func TestBasicSetAndGet(t *testing.T) {
	csmap := New(10)
	csmap.Set("a", ValTest{value: "b"})
	newVal, inMap := csmap.Get("a")
	require.Equal(t, inMap, true)
	require.Equal(t, newVal.(ValTest).value, "b")
}

// TestIterator tests the correctness of the iterator.
func TestIterator(t *testing.T) {
	csmap := New(5)
	csmap.Set("a", ValTest{value: "b"})
	csmap.Set("c", ValTest{value: "c"})
	csmap.Set("d", ValTest{value: "d"})
	csmap.Set("e", ValTest{value: "e"})
	csmap.Set("f", ValTest{value: "f"})
	newVal, inMap := csmap.Get("a")
	require.Equal(t, inMap, true)
	require.Equal(t, newVal.(ValTest).value, "b")

	tupChan := csmap.Iterator()
	keyLst := make([]string, 0)
	for tup := range tupChan {
		keyLst = append(keyLst, tup.key)
	}
	require.Equal(t, len(keyLst), 5)
}

// TestKeysWithConcurrency tests whether getting all the keys is correct
// adds a few concurrent goroutines to test this logic.
func TestKeysWithConcurrency(t *testing.T) {
	csmap := New(10)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()

		for i := 0; i < 50; i++ {
			valTest := ValTest{
				value: strconv.Itoa(i),
			}
			csmap.Set(strconv.Itoa(i), valTest)
		}
	}()


	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 50; i < 60; i++ {
			valTest := ValTest{
				value: strconv.Itoa(i),
			}
			csmap.Set(strconv.Itoa(i), valTest)
		}
	}()
	wg.Wait()

	for i := 0; i < 60; i++ {
		require.Equal(t, csmap.Has(strconv.Itoa(i)), true)
	}

	keyLst := csmap.Keys()
	var keyIntLst = []int{}

	for _, i := range keyLst {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		require.Equal(t, err, nil)
		keyIntLst = append(keyIntLst, j)
	}
	sort.Ints(keyIntLst)
	i := 0
	for _, keyInt := range keyIntLst {
		require.Equal(t, keyInt, i)
		i += 1
	}
}

// Test that we can get all values from different shards.
func TestValues(t *testing.T) {
	csmap := New(5)
	csmap.Set("1", ValTest{value: "b"})
	csmap.Set("2", ValTest{value: "c"})
	csmap.Set("3", ValTest{value: "d"})
	csmap.Set("4", ValTest{value: "e"})
	csmap.Set("5", ValTest{value: "f"})
	newVal, inMap := csmap.Get("1")
	require.Equal(t, inMap, true)
	require.Equal(t, newVal.(ValTest).value, "b")

	valueLst := csmap.Values()
	require.Equal(t, 5, len(valueLst))

	valueStrLst := make([]string, csmap.Count())
	for idx, val := range valueLst {
		valueStrLst[idx] = val.(ValTest).value
	}
	sort.Strings(valueStrLst)
	expectedStrLst := []string{"b", "c", "d", "e", "f"}
	for idx, val := range expectedStrLst {
		require.Equal(t, val, valueStrLst[idx])
	}
}